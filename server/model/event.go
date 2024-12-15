package model

import "time"

type Event struct {
	Id           string    `json:"id"`
	TriggerId    string    `json:"trigger_id"`
	EventTime    time.Time `json:"event_time"`
	Status       string    `json:"status"`
	ArchivedTime time.Time `json:"archived_time"`
	Manual       bool      `json:"is_manual"`
}
