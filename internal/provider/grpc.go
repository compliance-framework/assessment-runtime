package provider

import "context"

type grpcClient struct {
	client JobServiceClient
}

func (c *grpcClient) Evaluate(input *EvaluateInput) (*EvaluateResult, error) {
	return c.client.Evaluate(context.Background(), input)
}

func (c *grpcClient) Execute(input *ExecuteInput) (*ExecuteResult, error) {
	return c.client.Execute(context.Background(), input)
}

type grpcServer struct {
	Impl Provider
}

func (c *grpcServer) Evaluate(ctx context.Context, input *EvaluateInput) (*EvaluateResult, error) {
	return c.Impl.Evaluate(input)
}

func (c *grpcServer) Execute(ctx context.Context, input *ExecuteInput) (*ExecuteResult, error) {
	return c.Impl.Execute(input)
}
