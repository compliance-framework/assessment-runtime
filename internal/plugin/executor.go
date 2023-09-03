package plugin

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

type Executor struct {
	pluginManager *Pack
}

func NewExecutor(pluginManager *Pack) *Executor {
	return &Executor{
		pluginManager: pluginManager,
	}
}

func (e *Executor) ExecutePlugin(name string, input *ActionInput) (*ActionOutput, error) {
	client, ok := e.pluginManager.Clients[name]
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

	plg := raw.(Plugin)
	output, err := plg.Execute(input)
	if err != nil {
		log.WithFields(log.Fields{
			"plugin": name,
			"error":  err,
		}).Error("Failed to execute plugin")
		return nil, err
	}
	log.WithFields(log.Fields{
		"plugin": name,
		"output": output,
	}).Info("Plugin executed successfully")

	return output, nil
}
