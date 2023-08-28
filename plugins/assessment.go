package plugins

import (
	"fmt"
	"github.com/compliance-framework/assessment-runtime/config"
	goplugin "github.com/hashicorp/go-plugin"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

type Assessment struct {
	cfg     config.AssessmentConfig
	clients map[string]*goplugin.Client
	outputs map[string]*ActionOutput
}

func NewAssessment(cfg config.AssessmentConfig) (*Assessment, error) {
	a := &Assessment{
		cfg:     cfg,
		clients: make(map[string]*goplugin.Client),
		outputs: make(map[string]*ActionOutput),
	}

	pluginMap := make(map[string][]config.PluginConfig)
	for _, plugin := range a.cfg.Plugins {
		pluginMap[plugin.Package] = append(pluginMap[plugin.Package], plugin)
	}

	ex, err := os.Executable()
	if err != nil {
		return nil, err
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

	return a, nil
}

func (a *Assessment) Run() error {
	var wg sync.WaitGroup

	for _, plugin := range a.cfg.Plugins {
		wg.Add(1)
		go func(pluginName string) {
			defer wg.Done()

			for _, pluginConfig := range a.cfg.Plugins {
				if pluginConfig.Name != pluginName {
					continue
				}

				input := ActionInput{
					AssessmentId: a.cfg.AssessmentId,
					SSPId:        a.cfg.SSPId,
					ControlId:    a.cfg.ControlId,
					ComponentId:  a.cfg.ComponentId,
					Config:       pluginConfig.Configuration,
					Parameters:   pluginConfig.Parameters,
				}

				output, err := a.executePlugin(pluginName, input)
				if err != nil {
					log.WithField("plugin", pluginName).Error(err)
				}
				a.outputs[pluginName] = output
			}

		}(plugin.Name)
	}

	wg.Wait()

	return nil
}

func (a *Assessment) executePlugin(name string, input ActionInput) (*ActionOutput, error) {
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
	output, err := plugin.Execute(&input)
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

func (a *Assessment) Stop() {
	var wg sync.WaitGroup

	for _, client := range a.clients {
		wg.Add(1)
		go func(c *goplugin.Client) {
			defer wg.Done()
			c.Kill()
		}(client)
	}

	wg.Wait()
}
