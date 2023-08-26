package config

// Config represents the entire configuration loaded from the Yaml file.
type Config struct {
	ControlPlaneURL   string         `yaml:"controlPlaneURL" json:"controlPlaneURL"`
	PluginRegistryURL string         `yaml:"pluginRegistryURL" json:"pluginRegistryURL"`
	Plugins           []PluginConfig `yaml:"plugins" json:"plugins"`
}

// PluginConfig represents a plugins configuration.
type PluginConfig struct {
	Name     string `yaml:"name" json:"name"`
	Version  string `yaml:"version" json:"version"`
	Package  string `yaml:"package" json:"package"`
	Schedule string `yaml:"schedule" json:"schedule"`
}

type PluginPackageConfig struct {
	Name    string         `yaml:"name" json:"name"`
	Version string         `yaml:"version" json:"version"`
	Author  string         `yaml:"author" json:"author"`
	Plugins []PluginConfig `yaml:"plugins" json:"plugins"`
}
