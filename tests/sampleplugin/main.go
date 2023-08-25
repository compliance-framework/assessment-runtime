package main

import (
	"context"
	"fmt"
	"github.com/compliance-framework/assessment-runtime/plugins"
	"github.com/hashicorp/go-plugin"
)

type SamplePlugin struct{}

func (p *SamplePlugin) Init() error {
	fmt.Println("Plugin initialized")
	return nil
}

func (p *SamplePlugin) Execute(in *plugins.ActionInput) (*plugins.ActionOutput, error) {
	fmt.Println("Plugin executed")
	return &plugins.ActionOutput{}, nil
}

func (p *SamplePlugin) Shutdown(ctx context.Context) error {
	fmt.Println("Plugin shutdown")
	return nil
}

func main() {
	pluginSet := plugin.PluginSet{
		"sample": &plugins.AssessmentActionGRPCPlugin{Impl: &SamplePlugin{}},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugins.HandshakeConfig,
		Plugins:         pluginSet,
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}
