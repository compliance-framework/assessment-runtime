package plugin

import (
	"context"
	"fmt"
	"github.com/compliance-framework/assessment-runtime/internal/config"
	goplugin "github.com/hashicorp/go-plugin"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

type AssessmentRunner struct {
	cfg     config.AssessmentConfig
	clients map[string]*goplugin.Client
}

func NewAssessmentRunner(cfg config.AssessmentConfig) (*AssessmentRunner, error) {
	a := &AssessmentRunner{
		cfg:     cfg,
		clients: make(map[string]*goplugin.Client),
	}

	err := a.loadPlugins()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *AssessmentRunner) Run(ctx context.Context) map[string]*ActionOutput {
	outputs := make(map[string]*ActionOutput)
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
				outputs[pluginName] = &ActionOutput{
					Error: fmt.Errorf("execution cancelled").Error(),
				}
				mu.Unlock()
				return
			default:
				input := ActionInput{
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
					outputs[pluginName] = &ActionOutput{
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

func (a *AssessmentRunner) Stop() {
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

func (a *AssessmentRunner) loadPlugins() error {
	pluginMap := make(map[string][]config.PluginConfig)
	for _, plugin := range a.cfg.Plugins {
		pluginMap[plugin.Package] = append(pluginMap[plugin.Package], plugin)
	}

	ex, err := os.Executable()
	if err != nil {
		return err
	}

	for pkg, plugins := range pluginMap {
		log.WithField("package", pkg).Info("Loading package")

		pluginMap := make(map[string]goplugin.Plugin)
		for _, plugin := range plugins {
			log.WithField("plugin", plugin.Name).Info("Loading plugin")
			pluginMap[plugin.Name] = &AssessmentActionGRPCPlugin{}
		}
		pluginsPath := filepath.Join(filepath.Dir(ex), "./plugins")
		packagePath := fmt.Sprintf("%s/%s/%s/%s", pluginsPath, pkg, plugins[0].Version, pkg)

		log.WithFields(log.Fields{
			"package":     pkg,
			"pluginsPath": pluginsPath,
			"packagePath": packagePath,
		}).Info("Loading plugin package")

		client := goplugin.NewClient(&goplugin.ClientConfig{
			HandshakeConfig:  HandshakeConfig,
			Plugins:          pluginMap,
			Cmd:              exec.Command(packagePath),
			AllowedProtocols: []goplugin.Protocol{goplugin.ProtocolGRPC},
		})

		for _, plugin := range plugins {
			a.clients[plugin.Name] = client
		}
	}

	return nil
}

func (a *AssessmentRunner) executePlugin(name string, input *ActionInput) (*ActionOutput, error) {
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

	plugin := raw.(Plugin)
	output, err := plugin.Execute(input)
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
