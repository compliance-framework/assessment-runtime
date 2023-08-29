package config

// Config represents the entire configuration loaded from the Yaml file.
type Config struct {
	RuntimeId         string `yaml:"runtimeId" json:"runtimeId"`
	ControlPlaneURL   string `yaml:"controlPlaneURL" json:"controlPlaneURL"`
	PluginRegistryURL string `yaml:"pluginRegistryURL" json:"pluginRegistryURL"`
	EventBusURL       string `yaml:"eventBusURL" json:"eventBusURL"`
}

type AssessmentConfig struct {
	AssessmentId string         `yaml:"assessment-id" json:"assessmentId"`
	SSPId        string         `yaml:"ssp-id" json:"sspId"`
	ControlId    string         `yaml:"control-id" json:"controlId"`
	ComponentId  string         `yaml:"component-id" json:"componentId"`
	Schedule     string         `yaml:"schedule" json:"schedule"`
	Plugins      []PluginConfig `yaml:"plugins" json:"plugins"`
}

type PluginConfig struct {
	Name          string            `yaml:"name" json:"name"`
	Package       string            `yaml:"package" json:"package"`
	Version       string            `yaml:"version" json:"version"`
	Configuration map[string]string `yaml:"configuration" json:"configuration"`
	Parameters    map[string]string `yaml:"parameters" json:"parameters"`
}
