package bus

import (
	"github.com/compliance-framework/assessment-runtime/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSubToConfig(t *testing.T) {
	err := Connect("localhost:4222")
	if err != nil {
		t.Errorf("Failed to connect: %v", err)
	}

	ch := make(chan config.Config, 1)
	err = SubToConfig(ch)
	if err != nil {
		t.Errorf("Failed to subscribe to configuration: %v", err)
	}

	cfg := config.Config{} // Fill this with actual data
	err = PubConfig(cfg)
	if err != nil {
		t.Errorf("Failed to publish configuration: %v", err)
	}

	receivedCfg := <-ch
	assert.Equal(t, cfg, receivedCfg, "Received configuration does not match published one")
}
