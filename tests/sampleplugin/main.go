package main

import (
	"context"
	"github.com/compliance-framework/assessment-runtime/plugins"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

type SamplePlugin struct{}

func (p *SamplePlugin) Init() error {
	return nil
}

func (p *SamplePlugin) Execute(in *plugins.ActionInput) (*plugins.ActionOutput, error) {
	data := map[string]interface{}{
		"foo": "bar",
	}
	s, err := structpb.NewStruct(data)
	if err != nil {
		return nil, err
	}
	return &plugins.ActionOutput{
		ResultData: s,
	}, nil
}

func (p *SamplePlugin) Shutdown(ctx context.Context) error {
	return nil
}

func main() {
	plugins.Register(map[string]plugins.Plugin{
		"do-nothing": &SamplePlugin{},
	})
}
