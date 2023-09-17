package job

import (
	"context"
	"fmt"
	"github.com/compliance-framework/assessment-runtime/internal/model"
	"github.com/compliance-framework/assessment-runtime/internal/provider"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Runner struct {
	spec     model.JobSpec
	pack     *provider.Pack
	executor *provider.Executor
}

func NewRunner(spec model.JobSpec) (*Runner, error) {
	a := &Runner{
		spec: spec,
	}

	pack, err := provider.NewPluginPack(spec)
	if err != nil {
		return nil, err
	}
	a.pack = pack
	a.executor = provider.NewExecutor(pack)
	return a, nil
}

func (r *Runner) Run(ctx context.Context) map[string]*provider.ActionOutput {
	outputs := make(map[string]*provider.ActionOutput)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, activity := range r.spec.Activities {
		wg.Add(1)
		go func(pluginConfig *model.Plugin) {
			defer wg.Done()

			pluginName := pluginConfig.Name

			select {
			case <-ctx.Done():
				// TODO: Propagate cancellation to GRPC plugins
				log.WithField("plugin", pluginName).Info("execution cancelled")
				mu.Lock()
				outputs[pluginName] = &provider.ActionOutput{
					Error: fmt.Errorf("execution cancelled").Error(),
				}
				mu.Unlock()
				return
			default:
				input := provider.ActionInput{
					AssessmentId: r.spec.AssessmentId,
					SSPId:        r.spec.SspId,
				}

				output, err := r.executor.Execute(pluginName, &input)
				mu.Lock()
				if err != nil {
					outputs[pluginName] = &provider.ActionOutput{
						Error: err.Error(),
					}
					log.WithField("plugin", pluginName).Error(err)
				} else {
					outputs[pluginName] = output
				}
				mu.Unlock()
			}
		}(activity.Plugin)
	}

	wg.Wait()

	return outputs
}

func (r *Runner) Stop() {
	r.pack.Unload()
}
