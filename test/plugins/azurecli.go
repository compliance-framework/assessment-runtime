package main

import (
	"fmt"
	"time"
	"context"
	"os"
	"strings"

	. "github.com/compliance-framework/assessment-runtime/internal/provider"
	"github.com/google/uuid"
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
    "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
)

type AzureCliProvider struct {
	message string
}

func (p *AzureCliProvider) Evaluate(input *EvaluateInput) (*EvaluateResult, error) {
	var vmIds []string

    // Get environment variables and configuration data (not sure why subId is not env variable)
    clientId := os.Getenv("AZURE_CLIENT_ID")
    clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
    tenantId := os.Getenv("AZURE_TENANT_ID")
	subscriptionId, ok := input.Configuration["subscriptionId"]

    if clientId == "" || clientSecret == "" || tenantId == "" {
		return nil, fmt.Errorf("One or more environment variables are not set")
    }
	if !ok {
		return nil, fmt.Errorf("subscriptionId parameter is missing")
	}

    // Create a credential using Azure identity
    cred, err := azidentity.NewDefaultAzureCredential(nil)
    if err != nil {
		return nil, fmt.Errorf("failed to obtain a credential: %v", err)
    }

    // Create a VM Client
    vmClient, err := armcompute.NewVirtualMachinesClient(subscriptionId, cred, nil)
    if err != nil {
		return nil, fmt.Errorf("failed to create virtual machines client: %v", err)
    }

    // Create a context
    ctx := context.Background()

    // List all VMs in a subscription
    pager := vmClient.NewListAllPager(nil)
    for pager.More() {
        page, err := pager.NextPage(ctx)
        if err != nil {
			return nil, fmt.Errorf("failed to get next page of VMs: %v", err)
        }
        for _, vm := range page.Value {
			vmIds = append(vmIds, *vm.ID)
        }
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

	var obs *Observation
	var fndngs *Finding

    // Get environment variables and config variable.
    clientId := os.Getenv("AZURE_CLIENT_ID")
    clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
    tenantId := os.Getenv("AZURE_TENANT_ID")
	subscriptionId, ok := input.Configuration["subscriptionId"]

    if clientId == "" || clientSecret == "" || tenantId == "" {
		return nil, fmt.Errorf("One or more environment variables are not set")
    }

	if !ok {
		return nil, fmt.Errorf("subscriptionId parameter is missing")
	}

    // Create a context
    ctx := context.Background()

    // Create a credential using client ID, client secret, and tenant ID
    cred, err := azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret, nil)
    if err != nil {
		return nil, fmt.Errorf("failed to obtain a credential: %v", err)
    }

	// Create a VM client
    vmClient, err := armcompute.NewVirtualMachinesClient(subscriptionId, cred, nil)
    if err != nil {
		return nil, fmt.Errorf("failed to create virtual machines client: %v", err)
    }


	// Retrieve the VM ID from the subject properties
	vmId, ok := input.Subject.Props["id"]
	if !ok {
		return nil, fmt.Errorf("Vm Id is missing in subject properties")
	}
	start_time := time.Now().Format(time.RFC3339)

    // Extract resource group and VM name from VM ID
    resourceGroup, vmName, err := extractResourceGroupAndVMName(vmId)
    if err != nil {
		return nil, fmt.Errorf("failed to extract resource group and VM name: %v", err)
    }

    // Get the VM
    vm, err := vmClient.Get(ctx, resourceGroup, vmName, nil)
    if err != nil {
		return nil, fmt.Errorf("failed to get virtual machine: %v", err)
    }

    // Retrieve the tags
    tags := vm.Tags

	// Initialize variables to store the results
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

// extractResourceGroupAndVMName extracts the resource group and VM name from the VM ID
func extractResourceGroupAndVMName(vmID string) (string, string, error) {
    parts := strings.Split(vmID, "/")
    if len(parts) < 9 {
        return "", "", fmt.Errorf("invalid VM ID format")
    }
    return parts[4], parts[8], nil
}

func main() {
	Register(&AzureCliProvider{
		message: "Azure CLI provider completed",
	})
}
