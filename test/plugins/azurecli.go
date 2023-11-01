package main

import (
	"encoding/json"
	"fmt"
	"os/exec"

	. "github.com/compliance-framework/assessment-runtime/internal/provider"
)

type AzureCliProvider struct {
	message string
}

type VM struct {
	Name string            `json:"name"`
	Tags map[string]string `json:"tags"`
}

func (p *AzureCliProvider) Evaluate(input *EvaluateInput) (*EvaluateResult, error) {

	// Extract Azure Subscription ID from the parameters
	subscriptionId, ok := input.Selector.Labels["subscriptionId"]
	if !ok {
		return nil, fmt.Errorf("subscriptionId parameter is missing")
	}

	// Construct the Azure CLI command to list all VM IDs in the subscription
	cmd := exec.Command("az", "vm", "list", "--subscription", subscriptionId, "--query", "[].id", "--output", "json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute Azure CLI command: %s", err)
	}

	// Parse the output into a slice of VM IDs
	var vmIds []string
	if err := json.Unmarshal(out, &vmIds); err != nil {
		return nil, fmt.Errorf("failed to parse Azure CLI command output: %s", err)
	}

	// Create a list of subjects based on the VM IDs
	subjects := make([]*Subject, 0)
	for _, vmId := range vmIds {
		subjects = append(subjects, &Subject{
			Id:    vmId,
			Type:  SubjectType_INVENTORY_ITEM,
			Title: fmt.Sprintf("Azure Virtual Machine %s", vmId),
			Props: map[string]string{
				"id": vmId,
			},
		})
	}

	// Return the result with subjects and additional props if necessary
	return &EvaluateResult{
		Subjects: subjects,
	}, nil
}

func (p *AzureCliProvider) Execute(input *ExecuteInput) (*ExecuteResult, error) {
	// Retrieve the VM ID from the subject properties
	vmId, ok := input.Subject.Props["id"]
	if !ok {
		return nil, fmt.Errorf("VM ID is missing in subject properties")
	}

	// Construct the Azure CLI command to retrieve the tags of the specific VM
	cmd := exec.Command("az", "vm", "show", "--ids", vmId, "--query", "tags", "--output", "json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute Azure CLI command: %s", err)
	}

	// Parse the output into a map of tags
	var tags map[string]string
	if err := json.Unmarshal(out, &tags); err != nil {
		return nil, fmt.Errorf("failed to parse Azure CLI command output: %s", err)
	}

	// Check if the "dataclassification" tag exists
	_, hasTag := tags["dataclassification"]
	if !hasTag {
		// Create an observation if the tag is missing
		obs := &Observation{
			Id:          "123e4567-e89b-12d3-a456-426614174000",
			Title:       "Missing Data Classification Tag",
			Description: "The virtual machine does not have a 'dataclassification' tag.",
			Collected:   "2023-01-01T00:00:00Z",
			Expires:     "2023-12-31T23:59:59Z",
			Links:       []*Link{},
			Props: []*Property{
				{
					Name:  "Recommendation",
					Value: "Add a 'dataclassification' tag to this virtual machine.",
				},
			},
			Remarks: "The 'dataclassification' tag is required for compliance.",
		}

		// Log the observation
		logs := []*LogEntry{
			{
				Timestamp: "2023-01-01T00:00:00Z",
				Type:      LogType_WARNING,
				Details:   "The virtual machine does not have a 'dataclassification' tag.",
			},
		}

		// Return the result
		return &ExecuteResult{
			Status:       ExecutionStatus_SUCCESS,
			Observations: []*Observation{obs},
			Logs:         logs,
		}, nil
	}

	// If the "dataclassification" tag exists, return a successful result without observations
	return &ExecuteResult{
		Status: ExecutionStatus_SUCCESS,
	}, nil
}

func main() {
	Register(&AzureCliProvider{
		message: "Azure CLI provider completed",
	})
}
