package config

import (
	"encoding/json"
	"fmt"
	"github.com/compliance-framework/assessment-runtime/internal/bus"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type ConfigurationManager struct {
	config     Config
	jobConfigs []JobConfig
	client     *resty.Client
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

	jobs, err := cm.getJobConfigs()
	if err != nil {
		log.Warn("failed to get job configurations from control plane. loading jobs from disk")
		err = cm.loadJobConfigs(assessmentPath)
		if err != nil {
			return nil, err
		}
	} else {
		cm.jobConfigs = jobs

		err = cm.writeJobConfigs(jobs)
		if err != nil {
			return nil, err
		}
	}
	return cm, nil
}

func (cm *ConfigurationManager) Listen() {
	// Subscribe to job configuration updates
	ch, err := bus.Subscribe[JobConfig]("runtime.jobs")
	if err != nil {
		log.Errorf("failed to subscribe to job configuration updates: %s", err)
	}

	go func() {
		for {
			select {
			case jobConfig := <-ch:
				log.Infof("received job configuration update for assessment %s", jobConfig.AssessmentId)
				cm.jobConfigs = append(cm.jobConfigs, jobConfig)
			}
		}
	}()
}

func (cm *ConfigurationManager) getJobConfigs() ([]JobConfig, error) {
	resp, err := cm.client.R().Get(cm.config.ControlPlaneURL + "/runtime/jobs")
	if err != nil {
		return nil, err
	}

	var jobs []JobConfig
	err = json.Unmarshal(resp.Body(), &jobs)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (cm *ConfigurationManager) writeJobConfigs(jobConfigs []JobConfig) error {
	execDir, err := getExecutableDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(execDir, "assessments")

	for _, jobConfig := range jobConfigs {
		data, err := yaml.Marshal(jobConfig)
		if err != nil {
			return fmt.Errorf("failed to marshal yaml data: %w", err)
		}

		err = os.WriteFile(filepath.Join(configPath, jobConfig.AssessmentId+".yaml"), data, 0644)
		if err != nil {
			return fmt.Errorf("failed to write file: %w", err)
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

func (cm *ConfigurationManager) loadJobConfigs(path string) error {
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

			var config JobConfig
			err = yaml.Unmarshal(data, &config)
			if err != nil {
				return fmt.Errorf("failed to unmarshal yaml data: %w", err)
			}

			cm.jobConfigs = append(cm.jobConfigs, config)
		}
	}

	return nil
}

func (cm *ConfigurationManager) Config() Config {
	return cm.config
}

func (cm *ConfigurationManager) Packages() []Package {
	pluginInfoMap := make(map[string]Package)

	for _, config := range cm.jobConfigs {
		for _, plugin := range config.Plugins {
			key := plugin.Package + plugin.Version
			if _, exists := pluginInfoMap[key]; !exists {
				info := Package{
					Name:    plugin.Package,
					Version: plugin.Version,
				}
				pluginInfoMap[key] = info
			}
		}
	}

	pluginInfos := make([]Package, 0, len(pluginInfoMap))
	for _, info := range pluginInfoMap {
		pluginInfos = append(pluginInfos, info)
	}

	return pluginInfos
}

func (cm *ConfigurationManager) JobConfigs() []JobConfig {
	return cm.jobConfigs
}
