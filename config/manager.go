package config

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigurationError string

type ConfigurationManager struct {
	config Config
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

func (cm *ConfigurationManager) validate() error {
	if cm.config.ControlPlaneURL == "" {
		return ConfigurationError("controlPlaneAPI is empty")
	}

	for _, plugin := range cm.config.Plugins {
		if plugin.Name == "" {
			return ConfigurationError("plugin name is empty")
		}
		if plugin.Version == "" {
			return ConfigurationError(fmt.Sprintf("plugin version for '%s' is empty", plugin.Name))
		}
	}

	return nil
}

func (e ConfigurationError) Error() string {
	return string(e)
}
