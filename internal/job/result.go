package job

import "github.com/compliance-framework/assessment-runtime/internal/provider"

type Result struct {
	AssessmentId string
	Outputs      map[string]*provider.ExecuteResult
}
