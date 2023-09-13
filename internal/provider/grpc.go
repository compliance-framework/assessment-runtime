package provider

import "context"

type grpcClient struct {
	client ActionServiceClient
}

func (c *grpcClient) Execute(input *ActionInput) (*ActionOutput, error) {
	return c.client.Execute(context.Background(), input)
}

type grpcServer struct {
	Impl Plugin
}

func (c *grpcServer) Execute(ctx context.Context, input *ActionInput) (*ActionOutput, error) {
	return c.Impl.Execute(input)
}
