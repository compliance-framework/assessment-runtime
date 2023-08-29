package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	configYml := `
runtimeId: "123e4567-e89b-12d3-a456-426614174000"
controlPlaneAPI: "http://localhost:1234"
pluginRegistryURL: "http://plugin-registry"
eventBusURL: "nats://nats:4222"
`

	assessmentYml := `
assessment-id: "assess-1234"
ssp-id: "ssp-5678"
control-id: "ctrl-9101"
component-id: "comp-1121"
schedule: "*/10 * * * * *"
plugins:
  - name: "do-nothing"
    package: "sample"
    version: "1.0.0"
    configuration:
		config1: "value1"
		config2: "value2"
		config3: "value3"
	parameters:
		param1: "value1"
		param2: "value2"
		param3: "value3"
`

	tmpCfg, err := os.CreateTemp("", "config.yml")
	if err != nil {
		t.Fatal(err)
	}

	tmpAssessCfg, err := os.CreateTemp("assessments", "config.yml")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = os.Remove(tmpCfg.Name())
		_ = os.Remove(tmpAssessCfg.Name())
	}()

	if _, err := tmpCfg.Write([]byte(configYml)); err != nil {
		t.Fatal(err)
	}
	if err := tmpCfg.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := tmpAssessCfg.Write([]byte(assessmentYml)); err != nil {
		t.Fatal(err)
	}
	if err := tmpAssessCfg.Close(); err != nil {
		t.Fatal(err)
	}

	_, err = NewConfigurationManager()
	assert.Nil(t, err)
}
