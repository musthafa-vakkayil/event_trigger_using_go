package swagger

import "gorm.io/datatypes"

// User represents a user in the system
// @Description UserDto object used for input and output
type UserDto struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password"`
}

// Trigger represents a Trigger in the system
// @Description TriggerDto object used for input and output
type TriggerDto struct {
	Name         string         `json:"name" example:"API Trigger"`
	Type         string         `json:"type" example:"API/SCHEDULED"`
	ScheduleTime int            `json:"schedule_time"`
	Interval     int            `json:"interval_seconds"`
	ApiEndpoint  string         `json:"api_endpoint" example:"https://httpbin.org/get"`
	ApiPayload   datatypes.JSON `json:"api_payload"`
	ApiMethod    string         `json:"api_method" example:"GET"`
	Recurring    bool           `json:"is_recurring"`
}
