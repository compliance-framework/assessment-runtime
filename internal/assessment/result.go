package assessment

import "github.com/compliance-framework/assessment-runtime/internal/plugin"

type Result struct {
	AssessmentId string
	Outputs      map[string]*plugin.ActionOutput
}
