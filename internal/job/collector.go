package job

import (
	"github.com/compliance-framework/assessment-runtime/internal/event"
)

type Collector struct {
}

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) Process(results []Result) {
	// For now, we just publish the event to the event bus without any processing
	for _, r := range results {
		// Not handling the error case for now and depending on NATS retry mechanism
		_ = event.Publish[Result](r, `job.result`)
	}
}
