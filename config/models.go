package config

// Config represents the entire configuration loaded from the Yaml file.
type Config struct {
	ControlPlaneURL string   `yaml:"controlPlaneAPI"`
	Plugins         []Plugin `yaml:"plugins"`
}

// Plugin represents a plugin configuration.
type Plugin struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}
