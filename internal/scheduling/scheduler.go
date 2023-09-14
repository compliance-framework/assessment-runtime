package scheduling

import (
	"context"
	"github.com/compliance-framework/assessment-runtime/internal/job"
	"github.com/compliance-framework/assessment-runtime/internal/model"
	"github.com/compliance-framework/assessment-runtime/internal/pubsub"
	"sync"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

// JobFunc represents a function to be executed by the scheduler.
type JobFunc func()

// Scheduler represents a scheduler service.
type Scheduler struct {
	c            *cron.Cron
	jobTemplates []model.JobTemplate
	runners      sync.Map
	collector    *job.Collector
}

func NewScheduler(jobTemplates []model.JobTemplate) *Scheduler {
	s := &Scheduler{
		c:            cron.New(cron.WithSeconds()),
		jobTemplates: jobTemplates,
		collector:    job.NewCollector(),
	}
	return s
}

// Start starts the scheduler and runs the assessments based on the configured schedule.
func (s *Scheduler) Start(ctx context.Context) {
	for _, assessmentConfig := range s.jobTemplates {
		err := s.addJob(ctx, assessmentConfig)
		if err != nil {
			log.WithFields(log.Fields{
				"assessment-id": assessmentConfig.AssessmentId,
				"ssp-id":        assessmentConfig.SspId,
			}).Errorf("Failed to add assessment job: %s", err)
			// TODO: We should report this back to the control plane.
			continue
		}
	}

	s.c.Start()
}

// Stop stops the scheduler and running assessments.
func (s *Scheduler) Stop() {
	s.c.Stop()

	log.Info("Stopping scheduler")

	var wg sync.WaitGroup

	s.runners.Range(func(key, value interface{}) bool {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runner := value.(*job.Runner)
			runner.Stop()
		}()
		return true
	})

	wg.Wait()
}

// addJob adds an assessment job to the scheduler.
func (s *Scheduler) addJob(ctx context.Context, jobTemplate model.JobTemplate) error {
	jobFn := func() {
		runner, err := job.NewRunner(jobTemplate)
		if err != nil {
			log.WithFields(log.Fields{
				"assessment-id": jobTemplate.AssessmentId,
				"ssp-id":        jobTemplate.SspId,
			}).Errorf("Failed to create assessment: %s", err)

			pubsub.Publish(pubsub.Event{
				Type: pubsub.AssessmentFailed,
				Data: jobTemplate.AssessmentId,
			})
			return
		}

		defer runner.Stop()

		s.runners.Store(jobTemplate.AssessmentId, runner)
		result := runner.Run(ctx)
		s.collector.Process(job.Result{
			AssessmentId: jobTemplate.AssessmentId,
			Outputs:      result,
		})
		s.runners.Delete(jobTemplate.AssessmentId)
	}

	_, err := s.c.AddFunc(jobTemplate.Schedule, jobFn)
	return err
}
