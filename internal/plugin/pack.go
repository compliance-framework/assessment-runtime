package plugin

import (
	"fmt"
	"github.com/compliance-framework/assessment-runtime/internal/config"
	goplugin "github.com/hashicorp/go-plugin"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

type Pack struct {
	cfg     config.JobConfig
	Clients map[string]*goplugin.Client
}

func NewPluginPack(cfg config.JobConfig) (*Pack, error) {
	p := &Pack{
		cfg:     cfg,
		Clients: make(map[string]*goplugin.Client),
	}

	err := p.LoadPlugins()
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Pack) LoadPlugins() error {
	pluginMap := make(map[string][]config.PluginConfig)
	for _, pluginConfig := range p.cfg.Plugins {
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
			pluginMap[pluginConfig.Name] = &AssessmentActionGRPCPlugin{}
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

		for _, pluginConfig := range plugins {
			p.Clients[pluginConfig.Name] = client
		}
	}

	return nil
}

func (p *Pack) UnloadPlugins() {
	log.Info("unloading plugins")

	var wg sync.WaitGroup

	for _, client := range p.Clients {
		wg.Add(1)
		go func(c *goplugin.Client) {
			defer wg.Done()
			c.Kill()
		}(client)
	}

	wg.Wait()
}
