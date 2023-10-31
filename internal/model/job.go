package model

// JobSpec is the model used to communicate with the runtime
// It is used to publish a plan to the runtime. The runtime will then
// use the information to execute the activities and publish the results back to the control-plane.
// Here's an example tailored specifically for Azure Cloud:
// Task: "Assess Azure cloud's storage security configuration."
// Activities could include:
// - "Review the Azure Blob storage access policies and Private Endpoint settings."
// - "Check for encryption at rest and in transit for Azure storage accounts."
// - "Evaluate Azure Shared Access Signatures (SAS) and Azure Storage Service Encryption (SSE)."
// One more example:
// Task: "Verify Azure network security settings."
// Activities could include:
// - "Review Azure Network Security Groups (NSGs) to ensure least privilege access."
// - "Assess Virtual Private Network (VPN) and ExpressRoute configurations for secure connectivity."
// - "Check Azure DDoS Protection settings to ensure resilience against DDoS attacks."
// In this scenario, the task provides the overall direction for the assessment (e.g., assessing storage security or network security on Azure),
// while the activities break this task down into smaller, concrete steps to follow.
type JobSpec struct {
	Id          string `json:"id" yaml:"id"`
	Title       string `json:"title" yaml:"title"`
	PlanId      string `json:"assessment-plan-id" yaml:"assessment-plan-id"`
	ComponentId string `json:"component-id" yaml:"component-id"`
	ControlId   string `json:"control-id" yaml:"control-id"`
	Tasks       []Task `json:"tasks" yaml:"tasks"`
}

type Task struct {
	Id         string     `json:"id" yaml:"id"`
	Title      string     `json:"title" yaml:"title"`
	Schedule   string     `json:"schedule" yaml:"schedule"`
	Activities []Activity `json:"activities" yaml:"activities"`
}

type Activity struct {
	Id       string   `json:"id" yaml:"id"`
	Title    string   `json:"title" yaml:"title"`
	Selector Selector `json:"selector" yaml:"selector"`
	Provider Provider `json:"provider" yaml:"provider"`
}

type Selector struct {
	Title       string            `json:"title,omitempty" yaml:"title,omitempty"`
	Description string            `json:"description,omitempty" yaml:"description,omitempty"`
	Query       string            `json:"query,omitempty" yaml:"query,omitempty"`
	Labels      map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Expressions []Expression      `json:"expressions,omitempty" yaml:"expressions,omitempty"`
	Ids         []string          `json:"ids,omitempty" yaml:"ids,omitempty"`
}

type Provider struct {
	Name          string            `json:"name" yaml:"name"`
	Package       string            `json:"package" yaml:"package"`
	Version       string            `json:"version" yaml:"version"`
	Configuration map[string]string `json:"configuration" yaml:"configuration"`
}

type Expression struct {
	Key      string   `json:"key" yaml:"key"`
	Operator string   `json:"operator" yaml:"operator"`
	Values   []string `json:"values" yaml:"values"`
}
