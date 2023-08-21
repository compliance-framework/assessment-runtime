package plugins

import (
	"github.com/compliance-framework/assessment-runtime/config"
	log "github.com/sirupsen/logrus"
)

type PluginManager struct {
}

func NewPluginManager(cfg config.Config) *PluginManager {
	return &PluginManager{}
}

func (p *PluginManager) InitPlugins() error {
	log.Info("Initializing plugins")
	return nil
}

func (p *PluginManager) StartPlugin(name string) error {
	log.Infof("Starting plugin: %s", name)
	return nil
}

func (p *PluginManager) StopPlugin(name string) error {
	log.Infof("Stopping plugin: %s", name)
	return nil
}
