package plugins

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ActionService struct {
	Plugin Plugin
}

func (s *ActionService) Execute(ctx context.Context, in *ActionInput) (*ActionOutput, error) {
	if s.Plugin == nil {
		return nil, fmt.Errorf("plugins is not initialized")
	}
	return s.Plugin.Execute(ctx, in)
}

func (s *ActionService) Shutdown(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	if s.Plugin == nil {
		return nil, fmt.Errorf("plugins is nil")
	}

	if err := s.Plugin.Shutdown(ctx); err != nil {
		return nil, fmt.Errorf("failed to shutdown plugins: %w", err)
	}

	return &emptypb.Empty{}, nil
}
