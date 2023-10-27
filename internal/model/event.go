package model

type PlanPublished struct {
	Type string  `yaml:"type" json:"type"`
	Data JobSpec `yaml:"data" json:"data"`
}
