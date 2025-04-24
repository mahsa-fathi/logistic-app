package domain

import "time"

type PeriodicTask struct {
	ID               uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	JobName          string     `json:"job_name" gorm:"unique:not null"`
	IntervalInMinute uint       `json:"interval_in_minute" gorm:"not null"`
	LastRunTime      *time.Time `json:"last_run_time"`
	Failed           bool       `json:"failed" gorm:"not null;default:false"`
	Error            *string    `json:"error"`
	CreatedAt        time.Time  `json:"created_at" gorm:"not null"`
}
