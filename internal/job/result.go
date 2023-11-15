package job

import "github.com/compliance-framework/assessment-runtime/internal/provider"

// Result represents the result of a runner execution.
type Result struct {
	AssessmentId  string                  `json:"assessmentId"`
	ComponentId   string                  `json:"componentId"`
	ControlId     string                  `json:"controlId"`
	TaskId        string                  `json:"taskId"`
	ActivityId    string                  `json:"activityId"`
	Error         error                   `json:"error"`
	ExecuteResult *provider.ExecuteResult `json:"results"`
}
