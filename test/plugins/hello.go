package main

import (
	"context"
	. "github.com/compliance-framework/assessment-runtime/internal/plugin"
	"google.golang.org/protobuf/types/known/structpb"
)

type Hello struct {
}

func (p *Hello) Init() error {
	return nil
}

func (p *Hello) Execute(_ *ActionInput) (*ActionOutput, error) {
	data := map[string]interface{}{
		"message": "Hello World",
	}
	s, err := structpb.NewStruct(data)
	if err != nil {
		return nil, err
	}
	return &ActionOutput{
		ResultData: s,
	}, nil
}

func (p *Hello) Shutdown(context.Context) error {
	return nil
}

func main() {
	Register(map[string]Plugin{
		"hello-plugin": &Hello{},
	})
}
