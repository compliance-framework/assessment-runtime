package provider

import (
	goplugin "github.com/hashicorp/go-plugin"
	"os"
	"path/filepath"
	"strings"
)

func Register(provider Provider) {
	executablePath := os.Args[0]
	executableNameWithExtension := filepath.Base(executablePath)
	name := strings.TrimSuffix(executableNameWithExtension, ".exe")

	pluginSet := goplugin.PluginSet{}
	pluginSet[name] = &GrpcPlugin{Impl: provider}

	goplugin.Serve(&goplugin.ServeConfig{
		HandshakeConfig: HandshakeConfig,
		Plugins:         pluginSet,
		GRPCServer:      goplugin.DefaultGRPCServer,
	})
}
