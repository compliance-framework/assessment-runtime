package plugin

import (
	"context"
	"fmt"
	"github.com/compliance-framework/assessment-runtime/plugin"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ActionService struct {
	plugin plugin.Plugin
}

func (s *ActionService) Execute(ctx context.Context, in *ActionInput) (*ActionOutput, error) {
	if s.plugin == nil {
		return nil, fmt.Errorf("plugin is not initialized")
	}
	return s.plugin.Execute(ctx, in)
}

func (s *ActionService) Shutdown(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	if s.plugin == nil {
		return nil, fmt.Errorf("plugin is nil")
	}

	if err := s.plugin.Shutdown(ctx); err != nil {
		return nil, fmt.Errorf("failed to shutdown plugin: %w", err)
	}

	return &emptypb.Empty{}, nil
}
