package main

import (
	"context"
	"fmt"
	"github.com/compliance-framework/assessment-runtime/plugins"
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
	plugins.Register(map[string]plugins.Plugin{
		"do-nothing": &SamplePlugin{},
	})
}
