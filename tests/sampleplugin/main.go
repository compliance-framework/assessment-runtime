package main

import (
	"context"
	"fmt"
	"github.com/compliance-framework/assessment-runtime/internal/plugin"
)

func main() {
}

type SamplePlugin struct{}

func (p *SamplePlugin) Init() error {
	fmt.Println("Plugin initialized")
	return nil
}

func (p *SamplePlugin) Execute(ctx context.Context, in *plugin.ActionInput) (*plugin.ActionOutput, error) {
	fmt.Println("Plugin executed")
	return &plugin.ActionOutput{}, nil
}

func (p *SamplePlugin) Shutdown(ctx context.Context) error {
	fmt.Println("Plugin shutdown")
	return nil
}
