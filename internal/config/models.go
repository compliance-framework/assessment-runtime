package config

// Config represents the entire configuration loaded from the Yaml file.
type Config struct {
	RuntimeId         string `yaml:"runtimeId" json:"runtimeId"`
	ControlPlaneURL   string `yaml:"controlPlaneURL" json:"controlPlaneURL"`
	PluginRegistryURL string `yaml:"pluginRegistryURL" json:"pluginRegistryURL"`
	EventBusURL       string `yaml:"eventBusURL" json:"eventBusURL"`
}

// Pair represents a key-value pair.
type Pair struct {
	Name  string `yaml:"name" json:"name"`
	Value string `yaml:"value" json:"value"`
}

// JobConfig represents the configuration of a job (aka Assessment run).
type JobConfig struct {
	Uuid         string         `yaml:"uuid" json:"uuid" query:"uuid"`
	RuntimeUuid  string         `yaml:"runtime-id" json:"runtime-id"`
	SspId        string         `yaml:"ssp-id,omitempty" json:"ssp-id,omitempty"`
	AssessmentId string         `yaml:"assessment-id" json:"assessment-id"`
	TaskId       string         `yaml:"task-id" json:"task-id"`
	ActivityId   string         `yaml:"activity-id,omitempty" json:"activity-id,omitempty"`
	SubjectId    string         `yaml:"subject-id,omitempty" json:"subject-id,omitempty"`
	ControlId    string         `yaml:"control-id,omitempty" json:"control-id,omitempty"`
	Schedule     string         `yaml:"schedule" json:"schedule"`
	Plugins      []PluginConfig `yaml:"plugins,omitempty" json:"plugins,omitempty"`
	Parameters   []Pair         `yaml:"parameters,omitempty" json:"parameters,omitempty"`
}

// PluginConfig represents the configuration of a plugin.
type PluginConfig struct {
	Uuid          string `yaml:"uuid" json:"uuid"`
	Name          string `yaml:"name" json:"name"`
	Package       string `yaml:"package" json:"package"`
	Version       string `yaml:"version" json:"version"`
	Configuration []Pair `yaml:"configuration,omitempty" json:"configuration,omitempty"`
}

// Package represents a plugin package.
type Package struct {
	Name    string `yaml:"name" json:"name"`
	Version string `yaml:"version" json:"version"`
}

// EventConfigChanged represents the configuration updates sent to the event bus.
type EventConfigChanged struct {
	Type string    `yaml:"type" json:"type"`
	Uuid string    `yaml:"uuid" json:"uuid"`
	Data JobConfig `yaml:"data" json:"data"`
}
