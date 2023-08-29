package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/compliance-framework/assessment-runtime/internal"
	"gopkg.in/yaml.v3"
)

type ConfigurationError string

type ConfigurationManager struct {
	config            Config
	assessmentConfigs []AssessmentConfig
}

func NewConfigurationManager() (*ConfigurationManager, error) {
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}
	execDir := filepath.Dir(execPath)

	configPath := filepath.Join(execDir, "config.yml")
	assessmentPath := filepath.Join(execDir, "assessments")

	cm := &ConfigurationManager{}

	if err := cm.loadConfig(configPath); err != nil {
		return nil, err
	}

	if err := cm.loadAssessmentConfigs(assessmentPath); err != nil {
		return nil, err
	}

	return cm, nil
}

func (cm *ConfigurationManager) loadConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	err = yaml.Unmarshal(data, &cm.config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal yaml data: %w", err)
	}

	err = cm.validate()
	if err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	return nil
}

func (cm *ConfigurationManager) loadAssessmentConfigs(path string) error {
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

			var config AssessmentConfig
			err = yaml.Unmarshal(data, &config)
			if err != nil {
				return fmt.Errorf("failed to unmarshal yaml data: %w", err)
			}

			cm.assessmentConfigs = append(cm.assessmentConfigs, config)
		}
	}

	return nil
}

func (cm *ConfigurationManager) Config() Config {
	return cm.config
}

func (cm *ConfigurationManager) Packages() ([]runtime.PackageInfo, error) {
	pluginInfoMap := make(map[string]runtime.PackageInfo)

	for _, config := range cm.assessmentConfigs {
		for _, plugin := range config.Plugins {
			key := plugin.Package + plugin.Version
			if _, exists := pluginInfoMap[key]; !exists {
				info := runtime.PackageInfo{
					Name:    plugin.Package,
					Version: plugin.Version,
				}
				pluginInfoMap[key] = info
			}
		}
	}

	pluginInfos := make([]runtime.PackageInfo, 0, len(pluginInfoMap))
	for _, info := range pluginInfoMap {
		pluginInfos = append(pluginInfos, info)
	}

	return pluginInfos, nil
}

func (cm *ConfigurationManager) Assessments() []AssessmentConfig {
	return cm.assessmentConfigs
}

func (cm *ConfigurationManager) validate() error {
	if cm.config.RuntimeId == "" {
		return ConfigurationError("runtimeId is empty")
	}

	if cm.config.ControlPlaneURL == "" {
		return ConfigurationError("controlPlaneAPI is empty")
	}

	if cm.config.PluginRegistryURL == "" {
		return ConfigurationError("pluginRegistryURL is empty")
	}

	if cm.config.EventBusURL == "" {
		return ConfigurationError("eventBusURL is empty")
	}

	return nil
}

func (e ConfigurationError) Error() string {
	return string(e)
}
