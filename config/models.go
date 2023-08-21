package config

// Config represents the entire configuration loaded from the Yaml file.
type Config struct {
	ControlPlaneURL   string         `yaml:"controlPlaneURL"`
	PluginRegistryURL string         `yaml:"pluginRegistryURL"`
	Plugins           []PluginConfig `yaml:"plugins"`
}

// PluginConfig represents a plugins configuration.
type PluginConfig struct {
	Name     string `yaml:"name"`
	Version  string `yaml:"version"`
	Schedule string `yaml:"schedule"`
}
