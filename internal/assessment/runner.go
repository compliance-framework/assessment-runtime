package assessment

import (
	"context"
	"fmt"
	"github.com/compliance-framework/assessment-runtime/internal/config"
	"github.com/compliance-framework/assessment-runtime/internal/plugin"
	goplugin "github.com/hashicorp/go-plugin"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

type Runner struct {
	cfg     config.AssessmentConfig
	clients map[string]*goplugin.Client
}

func NewRunner(cfg config.AssessmentConfig) (*Runner, error) {
	a := &Runner{
		cfg:     cfg,
		clients: make(map[string]*goplugin.Client),
	}

	err := a.loadPlugins()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *Runner) Run(ctx context.Context) map[string]*plugin.ActionOutput {
	outputs := make(map[string]*plugin.ActionOutput)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, pluginConfig := range a.cfg.Plugins {
		wg.Add(1)
		go func(pluginConfig config.PluginConfig) {
			defer wg.Done()

			pluginName := pluginConfig.Name

			select {
			case <-ctx.Done():
				log.WithField("plugin", pluginName).Info("execution cancelled")
				mu.Lock()
				outputs[pluginName] = &plugin.ActionOutput{
					Error: fmt.Errorf("execution cancelled").Error(),
				}
				mu.Unlock()
				return
			default:
				input := plugin.ActionInput{
					AssessmentId: a.cfg.AssessmentId,
					SSPId:        a.cfg.SSPId,
					ControlId:    a.cfg.ControlId,
					ComponentId:  a.cfg.ComponentId,
					Config:       pluginConfig.Configuration,
					Parameters:   pluginConfig.Parameters,
				}

				output, err := a.executePlugin(pluginName, &input)
				mu.Lock()
				if err != nil {
					outputs[pluginName] = &plugin.ActionOutput{
						Error: err.Error(),
					}
					log.WithField("plugin", pluginName).Error(err)
				} else {
					outputs[pluginName] = output
				}
				mu.Unlock()
			}
		}(pluginConfig)
	}

	wg.Wait()

	return outputs
}

func (a *Runner) Stop() {
	log.Infof("stopping assessment %s", a.cfg.AssessmentId)

	var wg sync.WaitGroup

	for _, client := range a.clients {
		wg.Add(1)
		go func(c *goplugin.Client) {
			defer wg.Done()
			c.Kill()
		}(client)
	}

	wg.Wait()

	log.Infof("stopped assessment %s", a.cfg.AssessmentId)
}

func (a *Runner) loadPlugins() error {
	pluginMap := make(map[string][]config.PluginConfig)
	for _, pluginConfig := range a.cfg.Plugins {
		pluginMap[pluginConfig.Package] = append(pluginMap[pluginConfig.Package], pluginConfig)
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
			pluginMap[pluginConfig.Name] = &plugin.AssessmentActionGRPCPlugin{}
		}
		pluginsPath := filepath.Join(filepath.Dir(ex), "./plugins")
		packagePath := fmt.Sprintf("%s/%s/%s/%s", pluginsPath, pkg, plugins[0].Version, pkg)

		log.WithFields(log.Fields{
			"package":     pkg,
			"pluginsPath": pluginsPath,
			"packagePath": packagePath,
		}).Info("Loading plugin package")

		client := goplugin.NewClient(&goplugin.ClientConfig{
			HandshakeConfig:  plugin.HandshakeConfig,
			Plugins:          pluginMap,
			Cmd:              exec.Command(packagePath),
			AllowedProtocols: []goplugin.Protocol{goplugin.ProtocolGRPC},
		})

		for _, pluginConfig := range plugins {
			a.clients[pluginConfig.Name] = client
		}
	}

	return nil
}

func (a *Runner) executePlugin(name string, input *plugin.ActionInput) (*plugin.ActionOutput, error) {
	client, ok := a.clients[name]
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

	plg := raw.(plugin.Plugin)
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
	}).Info("Plugin executed successfully")

	return output, nil
}
