package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/compliance-framework/assessment-runtime/internal"
)

type ConfigurationError string

type ConfigurationManager struct {
	config            Config
	assessmentConfigs []AssessmentConfig
}

func NewConfigurationManager() *ConfigurationManager {
	return &ConfigurationManager{}
}

func (cm *ConfigurationManager) LoadConfig(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	err = yaml.Unmarshal(data, &cm.config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal yaml data: %w", err)
	}

	err = cm.validate()
	if err != nil {
		return cm.config, fmt.Errorf("config validation failed: %w", err)
	}

	return cm.config, nil
}

func (cm *ConfigurationManager) LoadAssessmentConfigs(path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".yaml" || filepath.Ext(file.Name()) == ".yml" {
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

func (cm *ConfigurationManager) Packages() ([]internal.PackageInfo, error) {
	pluginInfoMap := make(map[string]internal.PackageInfo)

	for _, config := range cm.assessmentConfigs {
		for _, plugin := range config.Plugins {
			key := plugin.Package + plugin.Version
			if _, exists := pluginInfoMap[key]; !exists {
				info := internal.PackageInfo{
					Name:    plugin.Package,
					Version: plugin.Version,
				}
				pluginInfoMap[key] = info
			}
		}
	}

	var pluginInfos []internal.PackageInfo
	for _, info := range pluginInfoMap {
		pluginInfos = append(pluginInfos, info)
	}

	return pluginInfos, nil
}

func (cm *ConfigurationManager) AssessmentConfigs() []AssessmentConfig {
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
