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
}

func NewAssessment(cfg config.AssessmentConfig) *Assessment {
	return &Assessment{
		cfg:     cfg,
		clients: make(map[string]*goplugin.Client),
	}
}

func (pm *Assessment) Init() error {
	pluginMap := make(map[string][]config.PluginConfig)
	for _, plugin := range pm.cfg.Plugins {
		pluginMap[plugin.Package] = append(pluginMap[plugin.Package], plugin)
	}

	ex, err := os.Executable()
	if err != nil {
		panic(err)
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
		pm.clients[pkg] = client
	}

	return nil
}

func (pm *Assessment) Execute(name string, input ActionInput) error {
	client, ok := pm.clients[name]
	if !ok {
		err := fmt.Errorf("plugin %s not found", name)
		log.WithField("plugin", name).Error(err)
		return err
	}

	grpcClient, err := client.Client()
	if err != nil {
		log.WithFields(log.Fields{
			"plugin": name,
			"error":  err,
		}).Error("Failed to get GRPC client for plugin")
		return err
	}

	raw, err := grpcClient.Dispense(name)
	if err != nil {
		log.WithFields(log.Fields{
			"plugin": name,
			"error":  err,
		}).Error("Failed to dispense plugin")
		return err
	}

	plugin := raw.(Plugin)
	output, err := plugin.Execute(&input)
	if err != nil {
		log.WithFields(log.Fields{
			"plugin": name,
			"error":  err,
		}).Error("Failed to execute plugin")
		return err
	}
	log.WithFields(log.Fields{
		"plugin": name,
		"output": output,
	}).Info("Plugin executed successfully")

	return nil
}

func (pm *Assessment) Stop() {
	var wg sync.WaitGroup

	for _, client := range pm.clients {
		wg.Add(1)
		go func(c *goplugin.Client) {
			defer wg.Done()
			c.Kill()
		}(client)
	}

	wg.Wait()
}
