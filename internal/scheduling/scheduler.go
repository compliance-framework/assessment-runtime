package scheduling

import (
	"context"
	"fmt"
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
	c         *cron.Cron
	specs     []model.JobSpec
	runners   sync.Map
	collector *job.Collector
}

func NewScheduler(jobSpecs []model.JobSpec) *Scheduler {
	s := &Scheduler{
		c:         cron.New(cron.WithSeconds()),
		specs:     jobSpecs,
		collector: job.NewCollector(),
	}
	return s
}

// Start starts the scheduler and runs the assessments based on the configured schedule.
func (s *Scheduler) Start(ctx context.Context) {
	s.loadJobs(ctx)

	// Listen for configuration updates
	ch, err := pubsub.Subscribe(pubsub.ConfigurationUpdated)
	if err != nil {
		fmt.Println("error subscribing to configuration updates:", err)
		return
	}

	go func() {
		for event := range ch {
			fmt.Println("received event:", event)
			s.cleanJobs()
			s.specs = event.Data.([]model.JobSpec)
			s.loadJobs(ctx)
		}
	}()

	s.c.Start()
}

// Stop stops the scheduler and running assessments.
func (s *Scheduler) Stop() {
	s.c.Stop()

	log.Info("Stopping scheduler")
}

func (s *Scheduler) loadJobs(ctx context.Context) {
	for _, spec := range s.specs {
		err := s.addJob(ctx, spec)
		if err != nil {
			log.WithFields(log.Fields{
				"id":                 spec.Id,
				"assessment-plan-id": spec.PlanId,
				"title":              spec.Title,
			}).Errorf("Failed to add assessment job: %s", err)
			// TODO: We should report this back to the control plane.
			continue
		}
	}
}

func (s *Scheduler) cleanJobs() {
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

	// Clean up crontab as well, as we will re-add all entries again
	for _, entry := range s.c.Entries() {
		s.c.Remove(entry.ID)
	}

	wg.Wait()
}

// addJob adds an assessment job to the scheduler.
func (s *Scheduler) addJob(ctx context.Context, spec model.JobSpec) error {
	jobFn := func() {
		runner, err := job.NewRunner(spec)
		if err != nil {
			log.WithFields(log.Fields{
				"id":                 spec.Id,
				"assessment-plan-id": spec.PlanId,
				"title":              spec.Title,
			}).Errorf("Failed to create assessment: %s", err)

			pubsub.Publish(pubsub.Event{
				Type: pubsub.AssessmentFailed,
				Data: fmt.Errorf("failed to create job runner: %w", err),
			})
			return
		}

		defer runner.Stop()

		s.runners.Store(spec.PlanId, runner)

		results := runner.Run(ctx)
		log.Info(results)

		s.collector.Process(results)
		s.runners.Delete(spec.PlanId)
	}

	for _, task := range spec.Tasks {
		_, err := s.c.AddFunc(task.Schedule, jobFn)

		if err != nil {
			log.WithFields(log.Fields{
				"id":                 spec.Id,
				"assessment-plan-id": spec.PlanId,
				"title":              spec.Title,
			}).Errorf("Failed to create assessment: %s", err)

			pubsub.Publish(pubsub.Event{
				Type: pubsub.AssessmentFailed,
				Data: fmt.Errorf("failed to add scheduling function: %w", err),
			})

			return err
		}
	}
	return nil
}
