package main

import (
	"fmt"

	"github.com/compliance-framework/assessment-runtime/plugin"
	"github.com/compliance-framework/assessment-runtime/plugin/proto"
)

func main() {
	plugin.Activate(func(in *proto.ActionInput) (*proto.ActionOutput, error) {
		fmt.Println(in)
		return &proto.ActionOutput{}, nil
	})
}
