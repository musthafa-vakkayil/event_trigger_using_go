package repo

import (
	"database/sql"
	"event_trigger/model"
	"fmt"
	"log"
)

func CreateEvent(db *sql.DB, eve model.Event) (string, error) {
	stmt := "INSERT INTO public.events(id, trigger_id, event_time, status, is_manual) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"

	var id string

	if err := db.QueryRow(stmt, eve.Id, eve.TriggerId, eve.EventTime, eve.Status, eve.Manual).Scan(&id); err != nil {
		log.Print(err)
		return "", err
	}

	fmt.Println("New trigger successfully created, id:", id)
	return id, nil
}

func ListEvents(db *sql.DB, includeArchived bool, includeAll bool) ([]model.Event, error) {
	// Query to fetch events
	var stmt string
	if includeAll {
		// Select all events (both ACTIVE and ARCHIVE)
		stmt = "SELECT id, trigger_id, event_time, status, is_manual FROM public.events"
	} else if includeArchived {
		// Select only archived events
		stmt = "SELECT id, trigger_id, event_time, status, is_manual FROM public.events WHERE status = 'ARCHIVE'"
	} else {
		// Select only active events
		stmt = "SELECT id, trigger_id, event_time, status, is_manual FROM public.events WHERE status = 'ACTIVE'"
	}

	// Execute the query
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse the results into a slice of Event
	var events []model.Event
	for rows.Next() {
		var event model.Event
		err := rows.Scan(&event.Id, &event.TriggerId, &event.EventTime, &event.Status, &event.Manual)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	// Return the list of events
	return events, nil
}
