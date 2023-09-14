package provider

import (
	goplugin "github.com/hashicorp/go-plugin"
)

func Register(name string, provider Provider) {
	pluginSet := goplugin.PluginSet{}
	pluginSet[name] = &GrpcPlugin{Impl: provider}

	goplugin.Serve(&goplugin.ServeConfig{
		HandshakeConfig: HandshakeConfig,
		Plugins:         pluginSet,
		GRPCServer:      goplugin.DefaultGRPCServer,
	})
}
