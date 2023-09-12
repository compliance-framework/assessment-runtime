package events

type JobConfig struct {
	Uuid         string     `yaml:"uuid" json:"uuid" query:"uuid"`
	RuntimeUuid  string     `yaml:"runtime-id" json:"runtime-id"`
	SspId        string     `yaml:"ssp-id,omitempty" json:"ssp-id,omitempty"`
	AssessmentId string     `yaml:"assessment-id" json:"assessment-id"`
	TaskId       string     `yaml:"task-id" json:"task-id"`
	ActivityId   string     `yaml:"activity-id,omitempty" json:"activity-id,omitempty"`
	Schedule     string     `yaml:"schedule" json:"schedule"`
	Activities   []Activity `yaml:"activities,omitempty" json:"activities,omitempty"`
}

type Pair struct {
	Name  string `yaml:"name" json:"name"`
	Value string `yaml:"value" json:"value"`
}

type Activity struct {
	Id         string `yaml:"id" json:"id"`
	ControlId  string `yaml:"control-id,omitempty" json:"control-id,omitempty"`
	Selector   Selector
	Plugins    []PluginConfig `yaml:"plugins,omitempty" json:"plugins,omitempty"`
	Parameters []Pair         `yaml:"parameters,omitempty" json:"parameters,omitempty"`
}

type Selector struct {
	Query       string            `yaml:"query" json:"query"`
	Labels      map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`
	Expressions []MatchExpression `yaml:"expressions,omitempty" json:"expressions,omitempty"`
	Ids         []string          `yaml:"ids,omitempty" json:"ids,omitempty"`
}

type MatchExpression struct {
	Key      string   `yaml:"key" json:"key"`
	Operator string   `yaml:"operator" json:"operator"`
	Values   []string `yaml:"values" json:"values"`
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

// ConfigChanged represents the configuration updates sent to the event bus.
type ConfigChanged struct {
	Type string    `yaml:"type" json:"type"`
	Uuid string    `yaml:"uuid" json:"uuid"`
	Data JobConfig `yaml:"data" json:"data"`
}
