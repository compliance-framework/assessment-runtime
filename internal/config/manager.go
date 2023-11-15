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
	config   Config
	jobSpecs []model.JobSpec
	client   *resty.Client
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
		client: resty.New().SetRetryCount(0).SetRetryWaitTime(5 * time.Second).SetRetryMaxWaitTime(20 * time.Second),
	}

	err = cm.loadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	jobs, err := cm.getJobSpecs()
	if err != nil {
		log.Warn("failed to get job configurations from control plane. loading jobs from local config")
		err = cm.loadJobSpecs(assessmentPath)
		if err != nil {
			return nil, err
		}
	} else {
		cm.jobSpecs = jobs

		err = cm.writeJobSpecs(jobs)
		if err != nil {
			return nil, err
		}
	}
	return cm, nil
}

func (cm *ConfigurationManager) Listen() {
	topic := "runtime.configuration" //fmt.Sprintf(, cm.config.RuntimeId)

	// Subscribe to job configuration updates
	ch, err := event.Subscribe[model.PlanEvent](topic)
	if err != nil {
		log.Errorf("failed to subscribe to job configuration updates: %s", err)
	}

	// Listen for job configuration updates
	// CM only manages config files, it doesn't work on running jobs
	go func() {
		for {
			select {
			case planEvent := <-ch:
				log.Infof("received job configuration change: %s", planEvent.Type)
				if planEvent.Type == "activated" || planEvent.Type == "updated" {
					err := cm.writeJobSpec(planEvent.Data)
					if err != nil {
						log.Errorf("failed to write job config: %s for job: %s", err, planEvent.Data.Id)
					}
				} else if planEvent.Type == "delete" {
					err := os.Remove(filepath.Join("assessments", planEvent.Data.Id+".yaml"))
					if err != nil {
						log.Errorf("failed to delete job config: %s for job: %s", err, planEvent.Data.Id)
					}
				}
			}
		}
	}()
}

func (cm *ConfigurationManager) getJobSpecs() ([]model.JobSpec, error) {
	resp, err := cm.client.R().Get(cm.config.ControlPlaneURL + "/runtime/jobs")
	if err != nil {
		return nil, err
	}

	var jobs []model.JobSpec
	err = json.Unmarshal(resp.Body(), &jobs)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (cm *ConfigurationManager) writeJobSpec(jobConfig model.JobSpec) error {
	execDir, err := getExecutableDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(execDir, "assessments")

	data, err := yaml.Marshal(jobConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal yaml data: %w", err)
	}

	err = os.WriteFile(filepath.Join(configPath, jobConfig.Id+".yaml"), data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (cm *ConfigurationManager) writeJobSpecs(jobConfigs []model.JobSpec) error {
	for _, jobConfig := range jobConfigs {
		err := cm.writeJobSpec(jobConfig)
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

func (cm *ConfigurationManager) loadJobSpecs(path string) error {
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

			var config model.JobSpec
			err = yaml.Unmarshal(data, &config)
			if err != nil {
				return fmt.Errorf("failed to unmarshal yaml data: %w", err)
			}

			cm.jobSpecs = append(cm.jobSpecs, config)
		}
	}

	return nil
}

func (cm *ConfigurationManager) Config() Config {
	return cm.config
}

func (cm *ConfigurationManager) Packages() []model.Package {
	pluginInfoMap := make(map[string]model.Package)

	for _, jobSpec := range cm.jobSpecs {
		for _, task := range jobSpec.Tasks {
			for _, activity := range task.Activities {
				key := activity.Provider.Package + activity.Provider.Version
				if _, exists := pluginInfoMap[key]; !exists {
					info := model.Package{
						Name:    activity.Provider.Package,
						Version: activity.Provider.Version,
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

func (cm *ConfigurationManager) JobSpecs() []model.JobSpec {
	return cm.jobSpecs
}
