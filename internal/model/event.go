package model

type ConfigChanged struct {
	Type string      `yaml:"type" json:"type"`
	Uuid string      `yaml:"uuid" json:"uuid"`
	Data JobTemplate `yaml:"data" json:"data"`
}
