package plugins

import (
	"context"
	"github.com/compliance-framework/assessment-runtime/internal/config"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"sync"
)

type JobFunc func()

type Scheduler struct {
	c                  *cron.Cron
	configs            []config.AssessmentConfig
	runningAssessments sync.Map
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
			s.runningAssessments.Store(assessmentConfig.AssessmentId, assessment)
			assessment.Run(ctx)
			assessment.Stop()
			s.runningAssessments.Delete(assessmentConfig.AssessmentId)
		})
		if err != nil {
			log.Fatal("failed to add job:", err)
		}
	}

	s.c.Start()
}

func (s *Scheduler) Stop() {
	s.c.Stop()

	log.Info("stopping scheduler")

	var wg sync.WaitGroup

	s.runningAssessments.Range(func(key, value interface{}) bool {
		wg.Add(1)
		go func() {
			assessment := value.(*AssessmentRunner)
			assessment.Stop()
			wg.Done()
		}()
		return true
	})

	wg.Wait()
}
