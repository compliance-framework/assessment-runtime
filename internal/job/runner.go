package job

import (
	"context"
	"fmt"
	"github.com/compliance-framework/assessment-runtime/internal/model"
	"github.com/compliance-framework/assessment-runtime/internal/provider"
	goplugin "github.com/hashicorp/go-plugin"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

type Runner struct {
	spec    model.JobSpec
	clients map[string]*goplugin.Client
}

func NewRunner(spec model.JobSpec) (*Runner, error) {
	a := &Runner{
		spec:    spec,
		clients: make(map[string]*goplugin.Client),
	}

	err := a.loadProviders()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (r *Runner) loadProviders() error {
	pluginMap := make(map[string][]model.Plugin)
	for _, activity := range r.spec.Activities {
		pluginMap[activity.Plugin.Package] = append(pluginMap[activity.Plugin.Package], *activity.Plugin)
	}

	ex, err := os.Executable()
	if err != nil {
		return err
	}

	for pkg, plugins := range pluginMap {
		log.WithField("package", pkg).Info("Loading package")

		pluginMap := make(map[string]goplugin.Plugin)
		for _, pluginConfig := range plugins {
			log.WithField("plugin", pluginConfig.Name).Info("Loading plugin")
			pluginMap[pluginConfig.Name] = &provider.GrpcPlugin{}
		}
		pluginsPath := filepath.Join(filepath.Dir(ex), "./plugins")
		packagePath := fmt.Sprintf("%s/%s/%s/%s", pluginsPath, pkg, plugins[0].Version, pkg)

		log.WithFields(log.Fields{
			"package":     pkg,
			"pluginsPath": pluginsPath,
			"packagePath": packagePath,
		}).Info("Loading plugin package")

		client := goplugin.NewClient(&goplugin.ClientConfig{
			HandshakeConfig:  provider.HandshakeConfig,
			Plugins:          pluginMap,
			Cmd:              exec.Command(packagePath),
			AllowedProtocols: []goplugin.Protocol{goplugin.ProtocolGRPC},
		})

		for _, pluginConfig := range plugins {
			r.clients[pluginConfig.Name] = client
		}
	}

	return nil
}

func (r *Runner) execute(name string, input *provider.ActionInput) (*provider.ActionOutput, error) {
	client, ok := r.clients[name]
	if !ok {
		err := fmt.Errorf("plugin %s not found", name)
		log.WithField("plugin", name).Error(err)
		return nil, err
	}

	grpcClient, err := client.Client()
	if err != nil {
		log.WithFields(log.Fields{
			"plugin": name,
			"error":  err,
		}).Error("Failed to get GRPC client for plugin")
		return nil, err
	}

	raw, err := grpcClient.Dispense(name)
	if err != nil {
		log.WithFields(log.Fields{
			"plugin": name,
			"error":  err,
		}).Error("Failed to dispense plugin")
		return nil, err
	}

	plg := raw.(provider.Provider)
	output, err := plg.Execute(input)
	if err != nil {
		log.WithFields(log.Fields{
			"plugin": name,
			"error":  err,
		}).Error("Failed to execute plugin")
		return nil, err
	}
	log.WithFields(log.Fields{
		"plugin": name,
		"output": output,
	}).Info("Provider executed successfully")

	return output, nil
}

func (r *Runner) Run(ctx context.Context) map[string]*provider.ActionOutput {
	outputs := make(map[string]*provider.ActionOutput)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, activity := range r.spec.Activities {
		wg.Add(1)
		go func(pluginConfig *model.Plugin) {
			defer wg.Done()

			pluginName := pluginConfig.Name

			select {
			case <-ctx.Done():
				// TODO: Propagate cancellation to GRPC plugins
				log.WithField("plugin", pluginName).Info("execution cancelled")
				mu.Lock()
				outputs[pluginName] = &provider.ActionOutput{
					Error: fmt.Errorf("execution cancelled").Error(),
				}
				mu.Unlock()
				return
			default:
				input := provider.ActionInput{
					AssessmentId: r.spec.AssessmentId,
					SSPId:        r.spec.SspId,
				}

				output, err := r.execute(pluginName, &input)
				mu.Lock()
				if err != nil {
					outputs[pluginName] = &provider.ActionOutput{
						Error: err.Error(),
					}
					log.WithField("plugin", pluginName).Error(err)
				} else {
					outputs[pluginName] = output
				}
				mu.Unlock()
			}
		}(activity.Plugin)
	}

	wg.Wait()

	return outputs
}

func (r *Runner) Stop() {
	log.Info("unloading providers")

	var wg sync.WaitGroup

	for _, client := range r.clients {
		wg.Add(1)
		go func(c *goplugin.Client) {
			defer wg.Done()
			c.Kill()
		}(client)
	}

	wg.Wait()
}
