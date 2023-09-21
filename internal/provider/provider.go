package provider

import (
	"context"
	goplugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type Provider interface {
	EvaluateSelector(*SubjectSelector) (*SubjectList, error)
	Execute(input *JobInput) (*JobResult, error)
}

type GrpcPlugin struct {
	goplugin.Plugin
	Impl Provider
}

func (p *GrpcPlugin) GRPCServer(broker *goplugin.GRPCBroker, s *grpc.Server) error {
	RegisterJobServiceServer(s, &grpcServer{Impl: p.Impl})
	return nil
}

func (p *GrpcPlugin) GRPCClient(_ context.Context, _ *goplugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &grpcClient{client: NewJobServiceClient(c)}, nil
}
