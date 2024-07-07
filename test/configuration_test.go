package test

import (
	"encoding/json"
	"github.com/compliance-framework/assessment-runtime/internal/config"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
)

func TestPublishEventConfigChanged(t *testing.T) {
	// Connect to NATS server
	nc, err := nats.Connect(nats.DefaultURL)
	require.NoError(t, err)
	defer nc.Close()

	// Sample EventConfigChanged array
	events := []config.EventConfigChanged{
		{
			Type: "created",
			Uuid: "uuid1",
			Data: config.JobConfig{
				Uuid:         "b6ad3c3e-c636-4c4c-8a73-3066c4e5e9bd",
				SspId:        "6b9f56e0-8ae9-45d0-9fdd-3db4e8a5bdad",
				AssessmentId: "c758890a-e6e8-4e7f-88fd-9a2e5ea26d86",
				TaskId:       "6b9f56e0-8ae9-45d0-9fdd-3db4e8a5bd22",
				ActivityId:   "555f56e0-8ae9-45d0-9fdd-3db4e8a5bdad",
				SubjectId:    "6b94eae0-8ae9-45d0-9fdd-3db4e8a5bdad",
				ControlId:    "6b9f56e0-8ae9-45d0-9fdd-3d5328a5bdad",
				Schedule:     "*/10 * * * * *",
				Plugins: []config.PluginConfig{
					{
						Name:    "busy-plugin",
						Version: "1.0.0",
						Configuration: []config.Pair{
							{
								Name:  "compliance-level",
								Value: "high",
							},
						},
					},
				},
			},
		},
	}

	// Convert the events array to JSON
	data, err := json.Marshal(events)
	require.NoError(t, err)

	// Publish the data to the runtime.configuration topic
	err = nc.Publish("runtime.configuration", data)
	require.NoError(t, err)
}
