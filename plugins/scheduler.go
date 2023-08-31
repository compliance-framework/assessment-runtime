package plugins

import (
	"context"
	"github.com/compliance-framework/assessment-runtime/config"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"sync/atomic"
	"time"
)

type JobFunc func()

type Scheduler struct {
	c                  *cron.Cron
	configs            []config.AssessmentConfig
	runningAssessments int32
}

func NewScheduler(assessmentConfigs []config.AssessmentConfig) *Scheduler {
	s := &Scheduler{
		c:       cron.New(cron.WithSeconds()),
		configs: assessmentConfigs,
	}
	return s
}

func (s *Scheduler) Start(ctx context.Context) {
	for _, assessmentConfig := range s.configs {
		_, err := s.c.AddFunc(assessmentConfig.Schedule, func() {
			assessment, err := NewAssessmentRunner(assessmentConfig)
			if err != nil {
				log.WithFields(log.Fields{
					"assessment-id": assessmentConfig.AssessmentId,
					"ssp-id":        assessmentConfig.SSPId,
					"control-id":    assessmentConfig.ControlId,
					"component-id":  assessmentConfig.ComponentId,
				}).Errorf("Failed to create assessment: %s", err)
				// TODO: We should report this back to the control plane.
				return
			}
			atomic.AddInt32(&s.runningAssessments, 1)
			assessment.Run(ctx)
			atomic.AddInt32(&s.runningAssessments, -1)
		})
		if err != nil {
			log.Fatal("Failed to add job:", err)
		}
	}

	s.c.Start()
}

func (s *Scheduler) Stop() {
	// TODO: We should wait for all running assessments to finish before stopping the scheduler.
	log.Info("Stopping scheduler")

	// Add a sleep of 10 seconds
	time.Sleep(10 * time.Second)

	s.c.Stop()
}
