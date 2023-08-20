package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	content := `
controlPlaneAPI: "http://localhost:1234"
plugins:
  - name: plugin1
    version: v1.2.3
`
	tmpfile, err := os.CreateTemp("", "config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer func(name string) {
		_ = os.Remove(name)
	}(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	cm := NewConfigurationManager()
	_, err = cm.LoadConfig(tmpfile.Name())
	assert.Nil(t, err)

	_, err = cm.LoadConfig("nonexistingfile.yaml")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to open config file")
}
