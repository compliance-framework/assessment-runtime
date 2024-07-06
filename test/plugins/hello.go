package main

import (
	. "github.com/compliance-framework/assessment-runtime/provider"
)

type Hello struct {
}

func (p *Hello) Evaluate(*EvaluateInput) (*EvaluateResult, error) {
	return nil, nil
}

func (p *Hello) Execute(input *ExecuteInput) (*ExecuteResult, error) {
	return &ExecuteResult{}, nil
}

func main() {
	Register(&Hello{})
}
