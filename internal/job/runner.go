package job

import (
	"context"
	"fmt"
	"github.com/compliance-framework/assessment-runtime/internal/model"
	"github.com/compliance-framework/assessment-runtime/internal/provider"
	goplugin "github.com/hashicorp/go-plugin"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

type Runner struct {
	spec    model.JobSpec
	clients map[string]*goplugin.Client
}

func NewRunner(spec model.JobSpec) (*Runner, error) {
	a := &Runner{
		spec:    spec,
		clients: make(map[string]*goplugin.Client),
	}

	err := a.loadProviders()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (r *Runner) loadProviders() error {
	pluginMap := make(map[string][]model.Provider)
	for _, task := range r.spec.Tasks {
		for _, activity := range task.Activities {
			pluginMap[activity.Provider.Package] = append(pluginMap[activity.Provider.Package], activity.Provider)
		}
	}

	ex, err := os.Executable()
	if err != nil {
		return err
	}

	for pkg, plugins := range pluginMap {
		log.WithField("package", pkg).Info("Loading package")

		pluginMap := make(map[string]goplugin.Plugin)
		for _, pluginConfig := range plugins {
			log.WithField("plugin", pluginConfig.Name).Info("Loading plugin")
			pluginMap[pluginConfig.Name] = &provider.GrpcPlugin{}
		}
		pluginsPath := filepath.Join(filepath.Dir(ex), "./plugins")
		packagePath := fmt.Sprintf("%s/%s/%s/%s", pluginsPath, pkg, plugins[0].Version, pkg)

		log.WithFields(log.Fields{
			"package":     pkg,
			"pluginsPath": pluginsPath,
			"packagePath": packagePath,
		}).Info("Loading plugin package")

		client := goplugin.NewClient(&goplugin.ClientConfig{
			HandshakeConfig:  provider.HandshakeConfig,
			Plugins:          pluginMap,
			Cmd:              exec.Command(packagePath),
			AllowedProtocols: []goplugin.Protocol{goplugin.ProtocolGRPC},
		})

		for _, pluginConfig := range plugins {
			r.clients[pluginConfig.Name] = client
		}
	}

	return nil
}

func (r *Runner) provider(name string) (provider.Provider, error) {
	client, ok := r.clients[name]
	if !ok {
		err := fmt.Errorf("plugin %s not found", name)
		log.WithField("plugin", name).Error(err)
		return nil, err
	}

	grpcClient, err := client.Client()
	if err != nil {
		log.WithFields(log.Fields{
			"plugin": name,
			"error":  err,
		}).Error("Failed to get GRPC client for plugin")
		return nil, err
	}

	raw, err := grpcClient.Dispense(name)
	if err != nil {
		log.WithFields(log.Fields{
			"plugin": name,
			"error":  err,
		}).Error("Failed to dispense plugin")
		return nil, err
	}

	return raw.(provider.Provider), nil
}

func (r *Runner) subjects(activityId string) ([]*provider.Subject, error) {
	for _, task := range r.spec.Tasks {
		for _, activity := range task.Activities {
			if activity.Id == activityId {
				p, err := r.provider(activity.Provider.Name)
				if err != nil {
					log.WithFields(log.Fields{
						"plugin": activity.Provider.Name,
						"error":  err,
					}).Error("failed to get provider")
					return nil, err
				}

				expressions := make([]*provider.Expression, 0)
				for _, expression := range activity.Selector.Expressions {
					expressions = append(expressions, &provider.Expression{
						Key:      expression.Key,
						Operator: expression.Operator,
						Values:   expression.Values,
					})
				}

				// TODO: Add missing information to the input: ComponentId, ControlId, etc.
				input := &provider.EvaluateInput{
					Plan: &provider.Plan{
						Id:         r.spec.PlanId,
						TaskId:     task.Id,
						ActivityId: activity.Id,
					},
					Selector: &provider.Selector{
						Query:       activity.Selector.Query,
						Labels:      activity.Selector.Labels,
						Expressions: expressions,
						Ids:         activity.Selector.Ids,
					},
				}
				subjects, err := p.Evaluate(input)

				if err != nil {
					log.WithFields(log.Fields{
						"provider": activity.Provider.Name,
						"error":    err,
					}).Error("failed to evaluate selector")
					return nil, err
				}

				return subjects.Subjects, nil
			}
		}
	}

	err := fmt.Errorf("activity %s not found", activityId)
	log.WithField("activity", activityId).Error(err)
	return nil, err
}

func (r *Runner) execute(name string, input *provider.ExecuteInput) (*provider.ExecuteResult, error) {
	p, err := r.provider(name)
	if err != nil {
		log.WithFields(log.Fields{
			"provider": name,
			"error":    err,
		}).Error("failed to get provider")
		return nil, err
	}

	result, err := p.Execute(input)
	if err != nil {
		log.WithFields(log.Fields{
			"plugin": name,
			"error":  err,
		}).Error("failed to execute plugin")
		return nil, err
	}

	log.WithFields(log.Fields{
		"plugin": name,
		"result": result,
	}).Info("provider executed successfully")

	return result, nil
}

func (r *Runner) Run(ctx context.Context) map[string]*provider.ExecuteResult {
	outputs := make(map[string]*provider.ExecuteResult)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, task := range r.spec.Tasks {
		for _, activity := range task.Activities {
			subjects, err := r.subjects(activity.Id)
			if err != nil {
				log.WithField("activity", activity.Id).Error(err)
				continue
			}

			if len(subjects) == 0 {
				log.WithField("activity", activity.Id).Info("no subjects found")
				continue
			}

			for _, subject := range subjects {
				wg.Add(1)

				go func(subject *provider.Subject, activity model.Activity) {
					defer wg.Done()

					pluginConfig := activity.Provider
					pluginName := pluginConfig.Name

					select {
					case <-ctx.Done():
						// TODO: Propagate cancellation to GRPC plugins
						log.WithField("plugin", pluginName).Info("execution cancelled")
						mu.Lock()
						outputs[pluginName] = &provider.ExecuteResult{}
						mu.Unlock()
						return
					default:
						params := make(map[string]string)
						for name, value := range activity.Params {
							params[name] = value
						}

						// TODO: Add missing information to the input: ComponentId, ControlId, etc.
						input := provider.ExecuteInput{
							Plan: &provider.Plan{
								Id:         r.spec.PlanId,
								TaskId:     task.Id,
								ActivityId: activity.Id,
							},
						}

						output, err := r.execute(pluginName, &input)
						mu.Lock()
						if err != nil {
							// TODO: Add error information
							outputs[pluginName] = &provider.ExecuteResult{}
							log.WithField("plugin", pluginName).Error(err)
						} else {
							outputs[pluginName] = output
						}
						mu.Unlock()
					}
				}(subject, activity)
			}

		}
	}

	wg.Wait()

	return outputs
}

func (r *Runner) Stop() {
	log.Info("unloading providers")

	var wg sync.WaitGroup

	for _, client := range r.clients {
		wg.Add(1)
		go func(c *goplugin.Client) {
			defer wg.Done()
			c.Kill()
		}(client)
	}

	wg.Wait()
}
