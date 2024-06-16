package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	. "github.com/compliance-framework/assessment-runtime/internal/provider"
	"github.com/google/uuid"
)

type AzureCliProvider struct {
	message string
}

// type VirtualMachine struct {
// 	ID   string `json:"id"`
// 	Name string `json:"name"`
// 	VmID string `json:"vmId"`
// }

func (p *AzureCliProvider) Evaluate(input *EvaluateInput) (*EvaluateResult, error) {
	// Extract Azure Subscription ID from the parameters
	subscriptionId, ok := input.Configuration["subscriptionId"]
	if !ok {
		return nil, fmt.Errorf("subscriptionId parameter is missing")
	}
	clientIdb, _ := os.ReadFile("/run/secrets/clientId")
	clientSecretb, _ := os.ReadFile("/run/secrets/clientSecret")
	tenantIdb, _ := os.ReadFile("/run/secrets/tenantId")
	clientId := strings.Replace(string(clientIdb), "\n", "", -1)
	clientSecret := strings.Replace(string(clientSecretb), "\n", "", -1)
	tenantId := strings.Replace(string(tenantIdb), "\n", "", -1)
	// Login to Azure CLI
	cmd := exec.Command("az", "login", "--service-principal", "-u", clientId, "-p", clientSecret, "--tenant", tenantId)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("List VMs: failed to login on Azure: %s\n\n%s", out, err)
	}
	// Setup Subscription
	cmd = exec.Command("az", "account", "set", "-s", subscriptionId)
	out, err = cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("List VMs: failed to login on Azure: %s\n\n%s", out, err)
	}

	// Construct the Azure CLI command to list all VM IDs in the subscription
	cmd = exec.Command("az", "vm", "list", "--subscription", subscriptionId, "--query", "[].id", "--output", "json")
	out, err = cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("List VMs: failed to execute Azure CLI command: %s\n\n%s", out, err)
	}

	// Parse the output into a slice of VirtualMachine structs
	var vmIds []string
	if err := json.Unmarshal(out, &vmIds); err != nil {
		return nil, fmt.Errorf("Parse VmIds: failed to parse Azure CLI command output: %s", err)
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
	start_time := time.Now().Format(time.RFC3339)

	if !ok {
		return nil, fmt.Errorf("Vm Id is missing in subject properties")
	}

	// Construct the Azure CLI command to retrieve the tags of the specific VM
	cmd := exec.Command("az", "vm", "show", "--ids", vmId, "--query", "tags", "--output", "json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Find Vm Tags: failed to execute Azure CLI command: %s", err)
	}

	// Parse the output into a map of tags
	var tags map[string]string
	// If there are no tags on any Vms, not parse the output
	if len(out) == 0 {
		fmt.Println("Parse Vm Tags: No tags found")
	} else {
		if err := json.Unmarshal(out, &tags); err != nil {
			return nil, fmt.Errorf("Parse Vm Tags: failed to parse Azure CLI command output: %s", err)
		}
	}

	// Initialize variables to store the results
	var obs *Observation
	var fndngs *Finding
	observations := []*Observation{}
	findings := []*Finding{}

	// Check if the "dataclassification" tag exists
	_, hasTag := tags["dataclassification"]
	obs_id := uuid.New().String()
	// Create an observation if the tag is either missing, or there.
	if !hasTag {
		obs = &Observation{
			Id:          obs_id,
			Title:       "Missing Data Classification Tag",
			Description: fmt.Sprintf("The virtual machine %s does not have a 'dataclassification' tag.", vmId),
			Collected:   time.Now().Format(time.RFC3339),
			Expires:     time.Now().AddDate(0, 1, 0).Format(time.RFC3339), // Add one month for the expiration
			Links:       []*Link{},
			Props: []*Property{
				{
					Name:  "VmId",
					Value: vmId,
				},
			},
			RelevantEvidence: []*Evidence{
				{
					Description: fmt.Sprintf("az cli command did not find any 'dataclassification' tag for the vm %s",vmId),
				},
			},
			Remarks: "The 'dataclassification' tag is required for compliance.",
		}
		fndngs = &Finding{
			Id:          uuid.New().String(),
			Title:       "Missing Data Classification Tag",
			Description: fmt.Sprintf("The virtual machine %s does not have a 'dataclassification' tag.", vmId),
			Remarks:     fmt.Sprintf("Give the virtual machine %s a 'dataclassification' tag.", vmId),
			RelatedObservations: []string{obs_id},
		}
		observations = append(observations, obs)
		findings = append(findings, fndngs)
	} else {
		obs = &Observation{
			Id:          obs_id,
			Title:       "Data Classification Tag Present",
			Description: fmt.Sprintf("The virtual machine %s has a 'dataclassification' tag.", vmId),
			Collected:   time.Now().Format(time.RFC3339),
			Expires:     time.Now().Format(time.RFC3339),
			Links:       []*Link{},
			Props: []*Property{
				{
					Name:  "VmId",
					Value: vmId,
				},
			},
			RelevantEvidence: []*Evidence{
				{
					Description: fmt.Sprintf("az cli command found a 'dataclassification' tag for the vm: %s", vmId),
				},
			},
			Remarks: "All OK.",
		}
		observations = append(observations, obs)
	}


	// Log that the check has successfully run
	logEntry := &LogEntry{
		Title:       "Data classification check",
		Description: "Data classification check has run successfully",
		Start:       start_time,
		End:         time.Now().Format(time.RFC3339),
	}

	// Return the result
	return &ExecuteResult{
		Status:       ExecutionStatus_SUCCESS,
		Observations: observations,
		Findings:     findings,
		Logs:         []*LogEntry{logEntry},
	}, nil
}

func main() {
	Register(&AzureCliProvider{
		message: "Azure CLI provider completed",
	})
}
