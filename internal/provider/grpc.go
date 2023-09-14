package provider

import "context"

type grpcClient struct {
	client JobServiceClient
}

func (c *grpcClient) EvaluateSelector(selector *SubjectSelector) (*SubjectList, error) {
	return c.client.EvaluateSelector(context.Background(), selector)
}

func (c *grpcClient) Execute(input *ActionInput) (*ActionOutput, error) {
	return c.client.Execute(context.Background(), input)
}

type grpcServer struct {
	Impl Plugin
}

func (c *grpcServer) EvaluateSelector(ctx context.Context, selector *SubjectSelector) (*SubjectList, error) {
	return c.Impl.EvaluateSelector(selector)
}

func (c *grpcServer) Execute(ctx context.Context, input *ActionInput) (*ActionOutput, error) {
	return c.Impl.Execute(input)
}
