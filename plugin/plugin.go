package plugin

import (
	"context"
	"github.com/compliance-framework/assessment-runtime/internal/plugin"
)

type Plugin interface {
	Init() error
	Execute(ctx context.Context, in *plugin.ActionInput) (*plugin.ActionOutput, error)
	Shutdown(ctx context.Context) error
}
