package main

import (
	"fmt"
	. "github.com/compliance-framework/assessment-runtime/provider"
	"os"
	"strconv"
)

type BusyPlugin struct {
	message string
}

func (p *BusyPlugin) Evaluate(*EvaluateInput) (*EvaluateResult, error) {
	// Normally we execute some queries here based on the data inside the EvaluateInput.Selector and return a subject list.

	subjects := make([]*Subject, 0)
	for i := 0; i < 3; i++ {
		// We fill the Props with the required information to be able to identify the subject later.
		// We can get this information inside the Execute function - ExecuteInput.Subject.Props["id"]
		subjects = append(subjects, &Subject{
			Id:    strconv.Itoa(i),
			Type:  SubjectType_INVENTORY_ITEM,
			Title: fmt.Sprintf("Busy Virtual Machine %d", i),
			Props: map[string]string{
				"id": strconv.Itoa(i),
			},
		})
	}

	// We can also add some props to the result. These props will be available in the Execute function to all the subjects.
	// ExecuteInput.Props["namespace"]
	return &EvaluateResult{
		Subjects: subjects,
		Props: map[string]string{
			"namespace": "busy",
		},
	}, nil
}

func (p *BusyPlugin) Execute(*ExecuteInput) (*ExecuteResult, error) {
	obs := &Observation{
		Id:          "123e4567-e89b-12d3-a456-426614174000",
		Title:       "Unencrypted Data Transmission",
		Description: "The automated assessment tool detected that the application transmits sensitive data without encryption.",
		Collected:   "2022-01-01T00:00:00Z",
		Expires:     "2022-12-31T23:59:59Z",
		Links: []*Link{
			{
				Rel:  "related",
				Href: "https://example.com/related-link",
			},
		},
		Props: []*Property{
			{
				Name:  "Risk Level",
				Value: "High",
			},
			{
				Name:  "Recommendation",
				Value: "Implement encryption methods for all data transmissions.",
			},
		},
		RelevantEvidence: []*Evidence{
			{
				Description: "Automated tool log indicating lack of encryption in data transmission",
			},
		},
		Remarks: "Immediate action required to mitigate potential data breaches.",
	}

	logs := []*LogEntry{
		{
			Title:       "Env test value",
			Description: os.Getenv("TEST"),
		},
		{
			Title:       "Sensitive Data Transmission",
			Description: "The automated assessment tool detected that the application transmits sensitive data without encryption.",
		},
	}

	return &ExecuteResult{
		Status:       ExecutionStatus_SUCCESS,
		Observations: []*Observation{obs},
		Logs:         logs,
	}, nil
}

func main() {
	Register(&BusyPlugin{
		message: "busy provider completed",
	})
}
