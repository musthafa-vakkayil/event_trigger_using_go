package model

import "time"

// @Description Event object used for input and output
type Event struct {
	Id        string    `json:"id"`
	TriggerId string    `json:"trigger_id"`
	EventTime time.Time `json:"event_time"`
	Status    string    `json:"status"`
	Manual    bool      `json:"is_manual"`
}
