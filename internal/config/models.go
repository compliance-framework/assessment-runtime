package config

// Config represents the entire configuration loaded from the Yaml file.
type Config struct {
	RuntimeId         string `yaml:"runtimeId" json:"runtimeId"`
	ControlPlaneURL   string `yaml:"controlPlaneURL" json:"controlPlaneURL"`
	PluginRegistryURL string `yaml:"pluginRegistryURL" json:"pluginRegistryURL"`
	EventBusURL       string `yaml:"eventBusURL" json:"eventBusURL"`
}

type Pair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type JobConfig struct {
	Uuid         string         `json:"uuid" query:"uuid"`
	RuntimeUuid  string         `json:"runtime-id"`
	SspId        string         `json:"ssp-id,omitempty"`
	AssessmentId string         `json:"assessment-id"`
	TaskId       string         `json:"task-id"`
	ActivityId   string         `json:"activity-id,omitempty"`
	SubjectId    string         `json:"subject-id,omitempty"`
	ControlId    string         `json:"control-id,omitempty"`
	Schedule     string         `json:"schedule"`
	Plugins      []PluginConfig `json:"plugins,omitempty"`
	Parameters   []Pair         `json:"parameters,omitempty"`
}

type PluginConfig struct {
	Uuid          string `yaml:"uuid"`
	Name          string `yaml:"name" json:"name"`
	Package       string `yaml:"package" json:"package"`
	Version       string `yaml:"version" json:"version"`
	Configuration []Pair `yaml:"configuration" json:"configuration,omitempty"`
}

type Package struct {
	Name    string `yaml:"name" json:"name"`
	Version string `yaml:"version" json:"version"`
}
