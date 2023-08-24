package plugins

import goplugin "github.com/hashicorp/go-plugin"

type Plugin interface {
	Execute(*ActionInput) (*ActionOutput, error)
	Shutdown() error
}

type GRPCPlugin struct {
	goplugin.Plugin
	Impl Plugin
}

func (p *GRPCPlugin) GRPCServer(broker *goplugin.GRPCBroker, s *goplugin.GRPCServer) error {
	RegisterActionServiceServer(s, &ActionService{Plugin: p.Impl})
	return nil
}

func (p *GRPCPlugin) GRPCClient(ctx *goplugin.GRPCBroker, c *goplugin.GRPCClient) (interface{}, error) {
	return &GRPCClient{client: NewPluginClient(c)}, nil
}
