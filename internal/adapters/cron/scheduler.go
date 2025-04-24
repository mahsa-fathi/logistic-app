package cron

import (
	"log"
	"logistic-app/internal/app/ports"
)

type Scheduler struct {
	service ports.Service
}

func NewScheduler(service ports.Service) *Scheduler {
	return &Scheduler{
		service: service,
	}
}

func (s *Scheduler) Run() {
	log.Println("Running scheduler for order update: ")
	s.service.ScheduleUpdateOrderStatus()
}
