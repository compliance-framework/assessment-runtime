package model

type JobTemplate struct {
	Id           string     `yaml:"id" json:"id" query:"id"`
	RuntimeId    string     `yaml:"runtime-id" json:"runtime-id"`
	SspId        string     `yaml:"ssp-id,omitempty" json:"ssp-id,omitempty"`
	AssessmentId string     `yaml:"assessment-id" json:"assessment-id"`
	TaskId       string     `yaml:"task-id" json:"task-id"`
	Schedule     string     `yaml:"schedule" json:"schedule"`
	Activities   []Activity `yaml:"activities,omitempty" json:"activities,omitempty"`
}

type Activity struct {
	Id         string    `yaml:"id" json:"id"`
	Selector   *Selector `json:"selector"`
	ControlId  string    `yaml:"control-id,omitempty" json:"control-id,omitempty"`
	Plugin     *Plugin   `yaml:"plugin,omitempty" json:"plugin,omitempty"`
	Parameters []*Pair   `yaml:"parameters,omitempty" json:"parameters,omitempty"`
}

type Selector struct {
	Query       string            `yaml:"query" json:"query"`
	Labels      map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`
	Expressions []Expression      `yaml:"expressions,omitempty" json:"expressions,omitempty"`
	Ids         []string          `yaml:"ids,omitempty" json:"ids,omitempty"`
}

type Plugin struct {
	Id            string  `json:"id"`
	Name          string  `json:"name"`
	Package       string  `json:"package"`
	Version       string  `json:"version"`
	Configuration []*Pair `json:"configuration"`
}

type Pair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Expression struct {
	Key      string   `json:"key"`
	Operator string   `json:"operator"`
	Values   []string `json:"values"`
}
