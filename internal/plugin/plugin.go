package plugin

import (
	"context"
	goplugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type Plugin interface {
	Execute(*ActionInput) (*ActionOutput, error)
}

type AssessmentActionGRPCPlugin struct {
	goplugin.Plugin
	Impl Plugin
}

func (p *AssessmentActionGRPCPlugin) GRPCServer(broker *goplugin.GRPCBroker, s *grpc.Server) error {
	RegisterActionServiceServer(s, &grpcServer{Impl: p.Impl})
	return nil
}

func (p *AssessmentActionGRPCPlugin) GRPCClient(_ context.Context, _ *goplugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &grpcClient{client: NewActionServiceClient(c)}, nil
}
