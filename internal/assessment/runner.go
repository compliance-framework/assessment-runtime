package assessment

import (
	"context"
	"fmt"
	"github.com/compliance-framework/assessment-runtime/internal/config"
	"github.com/compliance-framework/assessment-runtime/internal/plugin"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Runner struct {
	cfg      config.JobConfig
	pack     *plugin.Pack
	executor *plugin.Executor
}

func NewRunner(cfg config.JobConfig) (*Runner, error) {
	a := &Runner{
		cfg: cfg,
	}

	pluginManager, err := plugin.NewPluginPack(cfg)
	if err != nil {
		return nil, err
	}
	a.pack = pluginManager
	a.executor = plugin.NewExecutor(pluginManager)

	return a, nil
}

func (r *Runner) Run(ctx context.Context) map[string]*plugin.ActionOutput {
	outputs := make(map[string]*plugin.ActionOutput)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, pluginConfig := range r.cfg.Plugins {
		wg.Add(1)
		go func(pluginConfig config.PluginConfig) {
			defer wg.Done()

			pluginName := pluginConfig.Name

			select {
			case <-ctx.Done():
				log.WithField("plugin", pluginName).Info("execution cancelled")
				mu.Lock()
				outputs[pluginName] = &plugin.ActionOutput{
					Error: fmt.Errorf("execution cancelled").Error(),
				}
				mu.Unlock()
				return
			default:
				input := plugin.ActionInput{
					AssessmentId: r.cfg.AssessmentId,
					SSPId:        r.cfg.SspId,
					ControlId:    r.cfg.ControlId,
					ComponentId:  r.cfg.ControlId,
				}

				output, err := r.executor.ExecutePlugin(pluginName, &input)
				mu.Lock()
				if err != nil {
					outputs[pluginName] = &plugin.ActionOutput{
						Error: err.Error(),
					}
					log.WithField("plugin", pluginName).Error(err)
				} else {
					outputs[pluginName] = output
				}
				mu.Unlock()
			}
		}(pluginConfig)
	}

	wg.Wait()

	return outputs
}

func (r *Runner) Stop() {
	r.pack.UnloadPlugins()
}
