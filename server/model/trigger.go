package model

import (
	"time"

	"gorm.io/datatypes"
)

type Trigger struct {
	Id           string         `json:"id"`
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	ScheduleTime *time.Time     `json:"schedule_time"`
	Interval     int            `json:"interval_seconds"`
	ApiEndpoint  string         `json:"api_endpoint"`
	ApiPayload   datatypes.JSON `json:"api_payload"`
	ApiMethod    string         `json:"api_method"`
	Recurring    bool           `json:"is_recurring"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type TriggerDto struct {
	Name         string         `json:"name" validate:"required"`
	Type         string         `json:"type" validate:"required"`
	ScheduleTime int            `json:"schedule_time"`
	Interval     int            `json:"interval_seconds"`
	ApiEndpoint  string         `json:"api_endpoint"`
	ApiPayload   datatypes.JSON `json:"api_payload"`
	ApiMethod    string         `json:"api_method"`
	Recurring    bool           `json:"is_recurring"`
}

type EditTriggerDto struct {
	Name         *string         `json:"name"`
	Type         *string         `json:"type"`
	ScheduleTime *int            `json:"schedule_time"`
	Interval     int             `json:"interval_seconds"`
	ApiEndpoint  *string         `json:"api_endpoint"`
	ApiPayload   *datatypes.JSON `json:"api_payload"`
	ApiMethod    *string         `json:"api_method"`
	Recurring    *bool           `json:"is_recurring"`
}
