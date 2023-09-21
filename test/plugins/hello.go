package main

import (
	. "github.com/compliance-framework/assessment-runtime/internal/provider"
)

type Hello struct {
}

func (p *Hello) EvaluateSelector(_ *SubjectSelector) (*SubjectList, error) {
	return nil, nil
}

func (p *Hello) Execute(_ *JobInput) (*JobResult, error) {
	return &JobResult{}, nil
}

func main() {
	Register(&Hello{})
}
