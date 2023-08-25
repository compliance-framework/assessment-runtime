package plugins

import (
	goplugin "github.com/hashicorp/go-plugin"
	"os"

	log "github.com/sirupsen/logrus"
)

func Register(plugins map[string]Plugin) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)

	pluginSet := goplugin.PluginSet{}
	for name, plugin := range plugins {
		pluginSet[name] = &AssessmentActionGRPCPlugin{Impl: plugin}
	}

	goplugin.Serve(&goplugin.ServeConfig{
		HandshakeConfig: HandshakeConfig,
		Plugins:         pluginSet,
		GRPCServer:      goplugin.DefaultGRPCServer,
	})
}
