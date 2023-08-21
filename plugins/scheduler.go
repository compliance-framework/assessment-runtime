package plugins

import (
	"log"

	"github.com/robfig/cron/v3"
)

type JobFunc func()

type Scheduler struct {
	c *cron.Cron
}

func NewScheduler() *Scheduler {
	s := &Scheduler{
		c: cron.New(cron.WithSeconds()),
	}
	return s
}

func (s *Scheduler) Start() {
	s.c.Start()
}

func (s *Scheduler) AddJob(cronExpr string, job JobFunc) {
	_, err := s.c.AddFunc(cronExpr, func() {
		go job()
	})
	if err != nil {
		log.Fatal("Failed to add job:", err)
	}
}
