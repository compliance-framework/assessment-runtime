package scheduling

import (
	"context"
	"github.com/compliance-framework/assessment-runtime/internal/pubsub"
	"sync"

	"github.com/compliance-framework/assessment-runtime/internal/assessment"
	"github.com/compliance-framework/assessment-runtime/internal/config"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

// JobFunc represents a function to be executed by the scheduler.
type JobFunc func()

// Scheduler represents a scheduler service.
type Scheduler struct {
	c                  *cron.Cron
	configs            []config.JobConfig
	runningAssessments sync.Map
	collector          *assessment.Collector
}

func NewScheduler(assessmentConfigs []config.JobConfig) *Scheduler {
	s := &Scheduler{
		c:         cron.New(cron.WithSeconds()),
		configs:   assessmentConfigs,
		collector: assessment.NewCollector(),
	}
	return s
}

// Start starts the scheduler and runs the assessments based on the configured schedule.
func (s *Scheduler) Start(ctx context.Context) {
	for _, assessmentConfig := range s.configs {
		err := s.addJob(ctx, assessmentConfig)
		if err != nil {
			log.WithFields(log.Fields{
				"assessment-id": assessmentConfig.AssessmentId,
				"ssp-id":        assessmentConfig.SspId,
				"control-id":    assessmentConfig.ControlId,
				"component-id":  assessmentConfig.ControlId,
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

	s.runningAssessments.Range(func(key, value interface{}) bool {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runner := value.(*assessment.Runner)
			runner.Stop()
		}()
		return true
	})

	wg.Wait()
}

// addJob adds an assessment job to the scheduler.
func (s *Scheduler) addJob(ctx context.Context, assessmentConfig config.JobConfig) error {
	job := func() {
		runner, err := assessment.NewRunner(assessmentConfig)
		if err != nil {
			log.WithFields(log.Fields{
				"assessment-id": assessmentConfig.AssessmentId,
				"ssp-id":        assessmentConfig.SspId,
				"control-id":    assessmentConfig.ControlId,
				"component-id":  assessmentConfig.ControlId,
			}).Errorf("Failed to create assessment: %s", err)

			pubsub.Publish(pubsub.Event{
				Type: pubsub.AssessmentFailed,
				Data: assessmentConfig.AssessmentId,
			})
			return
		}

		defer runner.Stop()

		s.runningAssessments.Store(assessmentConfig.AssessmentId, runner)
		result := runner.Run(ctx)
		s.collector.Process(assessment.Result{
			AssessmentId: assessmentConfig.AssessmentId,
			Outputs:      result,
		})
		s.runningAssessments.Delete(assessmentConfig.AssessmentId)
	}

	_, err := s.c.AddFunc(assessmentConfig.Schedule, job)
	return err
}
