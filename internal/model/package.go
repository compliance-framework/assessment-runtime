package model

// Package represents a plugin package.
type Package struct {
	Name    string `yaml:"name" json:"name"`
	Tag     string `yaml:"tag" json:"tag"`
	Image   string `yaml:"image" json:"image"`
}
