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
	Id     string `json:"id" yaml:"id"`
	PlanId string `json:"assessment-plan-id" yaml:"assessment-plan-id"`
	Title  string `json:"title" yaml:"title"`
	Tasks  []Task `json:"tasks" yaml:"tasks"`
}

type Task struct {
	Id         string     `json:"id" yaml:"id"`
	Title      string     `json:"title" yaml:"title"`
	Schedule   string     `json:"schedule" yaml:"schedule"`
	Activities []Activity `json:"activities" yaml:"activities"`
}

type Activity struct {
	Id       string            `json:"id" yaml:"id"`
	Title    string            `json:"title" yaml:"title"`
	Selector Selector          `json:"selector" yaml:"selector"`
	Provider Provider          `json:"provider" yaml:"provider"`
	Params   map[string]string `json:"params" yaml:"params"`
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
	Name    string            `json:"name" yaml:"name"`
	Package string            `json:"package" yaml:"package"`
	Version string            `json:"version" yaml:"version"`
	Params  map[string]string `json:"params" yaml:"params"`
}

type Pair struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

type Expression struct {
	Key      string   `json:"key" yaml:"key"`
	Operator string   `json:"operator" yaml:"operator"`
	Values   []string `json:"values" yaml:"values"`
}

type Observation struct {
	SubjectId string `json:"subject-id" yaml:"subject-id"`
	Remarks   string `json:"description" yaml:"description"` // Holds the observation text (couldn't find a better name)
}

type Risk struct {
	SubjectId string `json:"subject-id" yaml:"subject-id"`
	Remarks   string `json:"description" yaml:"description"` // Holds the risk text
	Impact    string `json:"impact" yaml:"impact"`           // Holds the impact text
}

type JobResult struct {
	Id           string         `json:"id" yaml:"id"`
	RuntimeId    string         `json:"runtime-id" yaml:"runtime-id"` // only if the control-plane doesn't listen to runtime specific topic
	AssessmentId string         `json:"assessment-id" yaml:"assessment-id"`
	ActivityId   string         `json:"activity-id" yaml:"activity-id"`
	Observations []*Observation `json:"observations" yaml:"observations"`
	Risks        []*Risk        `json:"risks" yaml:"risks"`
}
