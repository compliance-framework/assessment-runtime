package job

import (
	"github.com/compliance-framework/assessment-runtime/internal/events"
	log "github.com/sirupsen/logrus"
)

type Collector struct {
}

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) Process(result Result) {
	log.WithFields(log.Fields{
		"assessment-id": result.AssessmentId,
	}).Infof("Processing result")

	// For now, we just publish the event to the event bus without any processing
	err := events.Publish[Result](result, `job.result`)
	if err != nil {
		return
	}
}
