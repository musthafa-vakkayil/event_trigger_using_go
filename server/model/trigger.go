package model

import (
	"time"

	"gorm.io/datatypes"
)

// @Description Trigger object used for input and output
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

// @Description TriggerDto object used for input and output
type TriggerDto struct {
	Name         string         `json:"name" validate:"required" example:"API Trigger"`
	Type         string         `json:"type" validate:"required" example:"API/SCHEDULED"`
	ScheduleTime int            `json:"schedule_time"`
	Interval     int            `json:"interval_seconds"`
	ApiEndpoint  string         `json:"api_endpoint" example:"https://httpbin.org/get"`
	ApiPayload   datatypes.JSON `json:"api_payload"`
	ApiMethod    string         `json:"api_method" example:"GET"`
	Recurring    bool           `json:"is_recurring"`
}

// @Description EditTriggerDto object used for input and output
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
