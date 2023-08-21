package config

// Config represents the entire configuration loaded from the Yaml file.
type Config struct {
	ControlPlaneURL   string   `yaml:"controlPlaneURL"`
	PluginRegistryURL string   `yaml:"pluginRegistryURL"`
	Plugins           []Plugin `yaml:"plugins"`
}

// Plugin represents a plugins configuration.
type Plugin struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}
