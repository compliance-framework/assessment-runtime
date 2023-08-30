package main

import (
	"context"
	"github.com/compliance-framework/assessment-runtime/plugins"
	"google.golang.org/protobuf/types/known/structpb"
)

type Hello struct {
}

func (p *Hello) Init() error {
	return nil
}

func (p *Hello) Execute(_ *plugins.ActionInput) (*plugins.ActionOutput, error) {
	data := map[string]interface{}{
		"message": "Hello World",
	}
	s, err := structpb.NewStruct(data)
	if err != nil {
		return nil, err
	}
	return &plugins.ActionOutput{
		ResultData: s,
	}, nil
}

func (p *Hello) Shutdown(ctx context.Context) error {
	return nil
}

func main() {
	plugins.Register(map[string]plugins.Plugin{
		"hello-plugin": &Hello{},
	})
}
