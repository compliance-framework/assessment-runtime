package plugin

import (
	"context"

	"github.com/compliance-framework/assessment-runtime/plugin/proto"
)

type PluginServer struct {
	proto.UnimplementedActionServiceServer
}

func (s *PluginServer) Execute(ctx context.Context, in *proto.ActionInput) (*proto.ActionOutput, error) {
	return &proto.ActionOutput{}, nil
}
