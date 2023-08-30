package plugins

import (
	"github.com/compliance-framework/assessment-runtime/config"
	"testing"
)

func TestStart(t *testing.T) {
	cfg := config.Config{
		ControlPlaneURL:   "localhost:0",
		PluginRegistryURL: "localhost:0",
		Plugins: []config.PluginConfig{
			{
				Name:    "testPlugin",
				Package: "testPackage",
				Version: "1.0",
			},
		},
	}

	pm := NewAssessmentRunner(cfg)
	err := pm.Init()

	if err != nil {
		t.Errorf("Init method returned an error: %v", err)
	}

	if len(pm.clients) == 0 {
		t.Errorf("No clients were started")
	}
}
