package main

import (
	"fmt"
	"log"
	"golang.org/x/crypto/ssh"

	"time"

	. "github.com/compliance-framework/assessment-runtime/provider"
	"github.com/google/uuid"
)

type SSHCommandProvider struct {
	message string
}

// SSHConfig contains the SSH connection configuration
type SSHConfig struct {
	User     string
	Password string
	Host     string
	Port     string
}

func (p *SSHCommandProvider) Evaluate(input *EvaluateInput) (*EvaluateResult, error) {
	username, ok := input.Configuration["username"]
	if !ok {
		return nil, fmt.Errorf("username parameter is missing")
	}
	// TODO: Get this from an env var, as it gets stored in the mongodb otherwise.
	_, ok = input.Configuration["password"]
	if !ok {
		return nil, fmt.Errorf("password parameter is missing")
	}
	host, ok := input.Configuration["host"]
	if !ok {
		return nil, fmt.Errorf("host parameter is missing")
	}
	command, ok := input.Configuration["command"]
	if !ok {
		return nil, fmt.Errorf("command parameter is missing")
	}
	port, ok := input.Configuration["port"]
	if !ok {
		// default to 22 if no port supplied
		port = "22"
	}
	// There is only one subject, so create one
	subjects := make([]*Subject, 0)
	ssh_target_id := fmt.Sprintf("%s@%s:%s %s", username, host, port, command)
	subjects = append(subjects, &Subject{
		Id:    ssh_target_id,
		Type:  SubjectType_INVENTORY_ITEM,
		Title: fmt.Sprintf("SSH target ssh %s", ssh_target_id),
		Props: map[string]string{
			"id": ssh_target_id,
		},
	})

	// Return the result with subjects and additional props if necessary
	return &EvaluateResult{
		Subjects: subjects,
	}, nil
}

func (p SSHCommandProvider) Execute(input *ExecuteInput) (*ExecuteResult, error) {
	start_time := time.Now().Format(time.RFC3339)
	username, ok := input.Configuration["username"]
	if !ok {
		return nil, fmt.Errorf("username parameter is missing")
	}
	password, ok := input.Configuration["password"]
	if !ok {
		return nil, fmt.Errorf("password parameter is missing")
	}
	host, ok := input.Configuration["host"]
	if !ok {
		return nil, fmt.Errorf("host parameter is missing")
	}
	command, ok := input.Configuration["command"]
	if !ok {
		return nil, fmt.Errorf("command parameter is missing")
	}
	port, ok := input.Configuration["port"]
	if !ok {
		// default to 22 if no port supplied
		port = "22"
	}
	config := SSHConfig{
		User:     username,
		Password: password,
		Host:     host,
		Port:     port,
	}

	var obs *Observation
	var fndngs *Finding

	observations := []*Observation{}
	findings := []*Finding{}

	obs_id := uuid.New().String()
	ssh_target_command := fmt.Sprintf("ssh -p %s %s@%s %s", port, username, host, command)

	// Run the command and get the output
	output, exit_code, err := RunCommand(config, command)
	if err != nil {
		log.Fatalf("Failed to run command: %v", err)
	}

	if (exit_code != 0) {
		// observation and finding
		obs = &Observation{
			Id:               obs_id,
			Title:            "SSH Command Did Not Succeed",
			Description:      fmt.Sprintf("The command: %s did not succeed.", ssh_target_command),
			Collected:        time.Now().Format(time.RFC3339),
			Expires:          time.Now().AddDate(0, 1, 0).Format(time.RFC3339), // Add one month for the expiration
			Links:            []*Link{},
			Props:            []*Property{
				{
					Name:  "Command",
					Value: fmt.Sprintf("%s", ssh_target_command),
				},
			},
			RelevantEvidence: []*Evidence{
				{
					Description: fmt.Sprintf("The command returned an exit code of %d for the command: %s", exit_code, ssh_target_command),
				},
			},
			Remarks:          fmt.Sprintf("The command: '%s' should return a zero exit code.", ssh_target_command),
		}
		fndngs = &Finding{
			Id:                  uuid.New().String(),
			Title:               "SSH Command Failure",
			Description:         fmt.Sprintf("The command %s did not succeed, and produced output: %s.", ssh_target_command, output),
			Remarks:             fmt.Sprintf("Correct the command %s.", ssh_target_command),
			RelatedObservations: []string{obs_id},
		}
		observations = append(observations, obs)
		findings = append(findings, fndngs)
	} else {
		// observation only
		obs = &Observation{
			Id:          obs_id,
			Title:       "SSH Command Succeeded",
			Description: fmt.Sprintf("The command: %s succeeded.", ssh_target_command),
			Collected:   time.Now().Format(time.RFC3339),
			Expires:     time.Now().AddDate(0, 1, 0).Format(time.RFC3339), // Add one month for the expiration
			Links:       []*Link{},
			Props: []*Property{
				{
					Name:  "Command",
					Value: fmt.Sprintf("%s", ssh_target_command),
				},
			},
			RelevantEvidence: []*Evidence{
				{
					Description: fmt.Sprintf("The command returned an exit code of %d for the command: %s", exit_code, ssh_target_command),
				},
			},
			Remarks: "All OK.",
		}
		observations = append(observations, obs)
	}

	// Log that the check has successfully run
	logEntry := &LogEntry{
		Title:       "SSH Command Check",
		Description: "SSH command check has run successfully",
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

// RunCommand executes a command on the remote server over SSH and returns the output
func RunCommand(config SSHConfig, command string) (string, int, error) {
	// Define the SSH client configuration
	sshConfig := &ssh.ClientConfig{
		User: config.User,
		Auth: []ssh.AuthMethod{
		    ssh.Password(config.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // For simplicity, ignore host key verification
	}

	// Establish the SSH connection
	address := fmt.Sprintf("%s:%s", config.Host, config.Port)
	client, err := ssh.Dial("tcp", address, sshConfig)
	if err != nil {
		return "", -1, fmt.Errorf("failed to dial: %v", err)
	}
	defer client.Close()

	// Create a session for the command execution
	session, err := client.NewSession()
	if err != nil {
		return "", -1, fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	// Execute the command and capture the output
	output, err := session.CombinedOutput(command)
	exit_code := -1
	if err != nil {
		if exitErr, ok := err.(*ssh.ExitError); ok {
		    exit_code = exitErr.ExitStatus()
		} else {
			return "", -1, fmt.Errorf("failed to execute command: %v", err)
		}
	} else {
		exit_code = 0
	}

	return string(output), exit_code, nil
}

func main() {
	Register(&SSHCommandProvider{
		message: "Azure CLI provider completed",
	})
}
