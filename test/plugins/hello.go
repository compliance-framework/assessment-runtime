package main

import (
	"context"
	plugins2 "github.com/compliance-framework/assessment-runtime/internal/plugins"
	"google.golang.org/protobuf/types/known/structpb"
)

type Hello struct {
}

func (p *Hello) Init() error {
	return nil
}

func (p *Hello) Execute(_ *plugins2.ActionInput) (*plugins2.ActionOutput, error) {
	data := map[string]interface{}{
		"message": "Hello World",
	}
	s, err := structpb.NewStruct(data)
	if err != nil {
		return nil, err
	}
	return &plugins2.ActionOutput{
		ResultData: s,
	}, nil
}

func (p *Hello) Shutdown(context.Context) error {
	return nil
}

func main() {
	plugins2.Register(map[string]plugins2.Plugin{
		"hello-plugin": &Hello{},
	})
}
