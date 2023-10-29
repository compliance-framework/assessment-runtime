package provider

import (
	"context"
	goplugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type Provider interface {
	// Evaluate evaluates the subject selector for the activity and returns a list of subjects that match the selector
	// This method runs for each Activity in the Task. The results are used to determine which subjects to run the activity against.
	// For example, if the activity is to run a compliance check on all Linux hosts, the selector would be something like:
	// "os == 'linux'"
	// The provider would then return a list of all subjects that match the selector.
	// The runtime would then run the activity against each subject in the list.
	Evaluate(*EvaluateInput) (*EvaluateResult, error)

	// Execute runs the activity against the subject and returns the results
	Execute(input *ExecuteInput) (*ExecuteResult, error)
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
