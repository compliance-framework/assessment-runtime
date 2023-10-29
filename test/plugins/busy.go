package main

import (
	. "github.com/compliance-framework/assessment-runtime/internal/provider"
	"strconv"
)

type BusyPlugin struct {
	message string
}

func (p *BusyPlugin) Evaluate(*EvaluateInput) (*EvaluateResult, error) {
	subjects := make([]*Subject, 0)
	for i := 0; i < 3; i++ {
		subjects = append(subjects, &Subject{Id: strconv.Itoa(i)})
	}

	return &EvaluateResult{
		Subjects: subjects,
	}, nil
}

func (p *BusyPlugin) Execute(*ExecuteInput) (*ExecuteResult, error) {
	obs := &Observation{
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
		Uuid:    "123e4567-e89b-12d3-a456-426614174000",
	}

	return &ExecuteResult{
		Status:       ExecutionStatus_SUCCESS,
		Observations: []*Observation{obs},
	}, nil
}

func main() {
	Register(&BusyPlugin{
		message: "busy provider completed",
	})
}
