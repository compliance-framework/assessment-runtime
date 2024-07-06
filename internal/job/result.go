package job

import "github.com/compliance-framework/assessment-runtime/provider"

// Result represents the result of a runner execution.
type Result struct {
	Status       provider.ExecutionStatus `json:"status"`
	AssessmentId string                   `json:"assessmentId"`
	ComponentId  string                   `json:"componentId"`
	ControlId    string                   `json:"controlId"`
	TaskId       string                   `json:"taskId"`
	ActivityId   string                   `json:"activityId"`
	Error        error                    `json:"error"`
	Subject      *provider.Subject        `json:"subjects"`
	Observations []*provider.Observation  `json:"observations"`
	Findings     []*provider.Finding      `json:"findings"`
	Risks        []*provider.Risk         `json:"risks"`
	Logs         []*provider.LogEntry     `json:"logs"`
}
