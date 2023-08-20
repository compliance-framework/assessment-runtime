package plugins

import (
	"context"
)

type Plugin interface {
	Init() error
	Execute(ctx context.Context, in *ActionInput) (*ActionOutput, error)
	Shutdown(ctx context.Context) error
}
