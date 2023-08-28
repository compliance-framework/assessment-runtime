package plugins

import (
	"github.com/compliance-framework/assessment-runtime/config"
	log "github.com/sirupsen/logrus"

	"github.com/robfig/cron/v3"
)

type JobFunc func()

type Scheduler struct {
	c *cron.Cron
}

func NewScheduler(assessmentConfigs []config.AssessmentConfig) *Scheduler {
	s := &Scheduler{
		c: cron.New(cron.WithSeconds()),
	}

	for _, assessmentConfig := range assessmentConfigs {
		_, err := s.c.AddFunc(assessmentConfig.Schedule, func() {
			assessment, err := NewAssessment(assessmentConfig)
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
			go assessment.Run()
		})
		if err != nil {
			log.Fatal("Failed to add job:", err)
		}
	}

	return s
}

func (s *Scheduler) Start() {
	s.c.Start()
}
