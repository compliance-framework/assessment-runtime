package config

// Config represents the entire configuration loaded from the Yaml file.
type Config struct {
	RuntimeId         string `yaml:"runtimeId" json:"runtimeId"`
	ControlPlaneURL   string `yaml:"controlPlaneURL" json:"controlPlaneURL"`
	PluginRegistryURL string `yaml:"pluginRegistryURL" json:"pluginRegistryURL"`
	EventBusURL       string `yaml:"eventBusURL" json:"eventBusURL"`
}

type AssessmentConfig struct {
	AssessmentID string         `yaml:"assessment-id" json:"assessmentId"`
	SspID        string         `yaml:"ssp-id" json:"sspId"`
	ControlID    string         `yaml:"control-id" json:"controlId"`
	ComponentID  string         `yaml:"component-id" json:"componentId"`
	Schedule     string         `yaml:"schedule" json:"schedule"`
	Plugins      []PluginConfig `yaml:"plugins" json:"plugins"`
}

type PluginConfig struct {
	Name          string  `yaml:"name" json:"name"`
	Package       string  `yaml:"package" json:"package"`
	Version       string  `yaml:"version" json:"version"`
	Configuration []Entry `yaml:"configuration" json:"configuration"`
	Parameters    []Entry `yaml:"parameters" json:"parameters"`
}

type Entry struct {
	Name  string `yaml:"name" json:"name"`
	Value string `yaml:"value" json:"value"`
}
