package model

// Package represents a plugin package.
type Package struct {
	Name    string `yaml:"name" json:"name"`
	Version string `yaml:"version" json:"version"`
}
