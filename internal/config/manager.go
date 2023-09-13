package config

import (
	"encoding/json"
	"fmt"
	"github.com/compliance-framework/assessment-runtime/internal/event"
	"github.com/compliance-framework/assessment-runtime/internal/model"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the entire configuration loaded from the Yaml file.
type Config struct {
	RuntimeId         string `yaml:"runtimeId" json:"runtimeId"`
	ControlPlaneURL   string `yaml:"controlPlaneURL" json:"controlPlaneURL"`
	PluginRegistryURL string `yaml:"pluginRegistryURL" json:"pluginRegistryURL"`
	EventBusURL       string `yaml:"eventBusURL" json:"eventBusURL"`
}

type ConfigurationManager struct {
	config       Config
	jobTemplates []model.JobTemplate
	client       *resty.Client
}

func getExecutableDir() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}
	return filepath.Dir(execPath), nil
}

func NewConfigurationManager() (*ConfigurationManager, error) {
	execDir, err := getExecutableDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(execDir, "config.yml")
	assessmentPath := filepath.Join(execDir, "assessments")

	cm := &ConfigurationManager{
		client: resty.New().SetRetryCount(3).SetRetryWaitTime(5 * time.Second).SetRetryMaxWaitTime(20 * time.Second),
	}

	err = cm.loadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	jobs, err := cm.getJobTemplates()
	if err != nil {
		log.Warn("failed to get job configurations from control plane. loading jobs from local config")
		err = cm.loadJobTemplates(assessmentPath)
		if err != nil {
			return nil, err
		}
	} else {
		cm.jobTemplates = jobs

		err = cm.writeJobTemplates(jobs)
		if err != nil {
			return nil, err
		}
	}

	cm.Listen()

	return cm, nil
}

func (cm *ConfigurationManager) Listen() {
	topic := "runtime.configuration" //fmt.Sprintf(, cm.config.RuntimeId)

	// Subscribe to job configuration updates
	ch, err := event.Subscribe[[]model.ConfigChanged](topic)
	if err != nil {
		log.Errorf("failed to subscribe to job configuration updates: %s", err)
	}

	// Listen for job configuration updates
	// CM only manages config files, it doesn't work on running jobs
	go func() {
		for {
			select {
			case changes := <-ch:
				for _, change := range changes {
					log.Infof("received job configuration change: %s", change.Type)
					if change.Type == "created" || change.Type == "updated" {
						err := cm.writeJobTemplate(change.Data)
						if err != nil {
							log.Errorf("failed to write job config: %s for job: %s", err, change.Data.Uuid)
						}
					} else if change.Type == "delete" {
						err := os.Remove(filepath.Join("assessments", change.Data.Uuid+".yaml"))
						if err != nil {
							log.Errorf("failed to delete job config: %s for job: %s", err, change.Data.Uuid)
						}
					}
				}
			}
		}
	}()
}

func (cm *ConfigurationManager) getJobTemplates() ([]model.JobTemplate, error) {
	resp, err := cm.client.R().Get(cm.config.ControlPlaneURL + "/runtime/jobs")
	if err != nil {
		return nil, err
	}

	var jobs []model.JobTemplate
	err = json.Unmarshal(resp.Body(), &jobs)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (cm *ConfigurationManager) writeJobTemplate(jobConfig model.JobTemplate) error {
	execDir, err := getExecutableDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(execDir, "assessments")

	data, err := yaml.Marshal(jobConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal yaml data: %w", err)
	}

	err = os.WriteFile(filepath.Join(configPath, jobConfig.Uuid+".yaml"), data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (cm *ConfigurationManager) writeJobTemplates(jobConfigs []model.JobTemplate) error {
	for _, jobConfig := range jobConfigs {
		err := cm.writeJobTemplate(jobConfig)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cm *ConfigurationManager) loadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	err = yaml.Unmarshal(data, &cm.config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal yaml data: %w", err)
	}

	return nil
}

func (cm *ConfigurationManager) loadJobTemplates(path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, file := range files {
		fileExt := filepath.Ext(file.Name())
		if fileExt == ".yaml" || fileExt == ".yml" {
			data, err := os.ReadFile(filepath.Join(path, file.Name()))
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}

			var config model.JobTemplate
			err = yaml.Unmarshal(data, &config)
			if err != nil {
				return fmt.Errorf("failed to unmarshal yaml data: %w", err)
			}

			cm.jobTemplates = append(cm.jobTemplates, config)
		}
	}

	return nil
}

func (cm *ConfigurationManager) Config() Config {
	return cm.config
}

func (cm *ConfigurationManager) Packages() []model.Package {
	pluginInfoMap := make(map[string]model.Package)

	for _, template := range cm.jobTemplates {
		for _, activity := range template.Activities {
			for _, plugin := range activity.Plugins {
				key := plugin.Package + plugin.Version
				if _, exists := pluginInfoMap[key]; !exists {
					info := model.Package{
						Name:    plugin.Package,
						Version: plugin.Version,
					}
					pluginInfoMap[key] = info
				}
			}
		}
	}

	packages := make([]model.Package, 0, len(pluginInfoMap))
	for _, info := range pluginInfoMap {
		packages = append(packages, info)
	}

	return packages
}

func (cm *ConfigurationManager) JobTemplates() []model.JobTemplate {
	return cm.jobTemplates
}
