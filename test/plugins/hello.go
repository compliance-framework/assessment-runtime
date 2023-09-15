package main

import (
	. "github.com/compliance-framework/assessment-runtime/internal/provider"
	"google.golang.org/protobuf/types/known/structpb"
)

type Hello struct {
}

func (p *Hello) EvaluateSelector(_ *SubjectSelector) (*SubjectList, error) {
	return nil, nil
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

func main() {
	Register(&Hello{})
}
