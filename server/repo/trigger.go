package repo

import (
	"database/sql"
	"event_trigger/model"
	"fmt"
	"log"
)

func CreateTrigger(db *sql.DB, trr model.Trigger) (string, error) {
	stmt := "INSERT INTO public.triggers(id, name, type, schedule_time, interval_seconds, api_endpoint, api_payload, api_method, is_recurring) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id"

	var id string

	if err := db.QueryRow(stmt, trr.Id, trr.Name, trr.Type, trr.ScheduleTime, trr.Interval, trr.ApiEndpoint, trr.ApiPayload, trr.ApiMethod, trr.Recurring).Scan(&id); err != nil {
		log.Print(err)
		return "", err
	}

	fmt.Println("New trigger successfully created, id:", id)
	return id, nil
}

func ListTriggers(db *sql.DB) ([]model.Trigger, error) {
	rows, err := db.Query("SELECT id, name, type, schedule_time, interval_seconds, api_endpoint, api_payload, api_method, is_recurring, created_at, updated_at  FROM public.triggers")
	if err != nil {
		return []model.Trigger{}, nil
	}

	var triggers []model.Trigger

	// Loop Through each item from the response and create a array of objects
	for rows.Next() {
		var triger model.Trigger
		if err := rows.Scan(
			&triger.Id, &triger.Name, &triger.Type, &triger.ScheduleTime, &triger.Interval,
			&triger.ApiEndpoint, &triger.ApiPayload, &triger.ApiMethod, &triger.Recurring, &triger.CreatedAt, &triger.UpdatedAt,
		); err != nil {
			return nil, err
		}

		triggers = append(triggers, triger)
	}

	return triggers, nil
}

func GetTriggerByID(db *sql.DB, trigger_id string) (model.Trigger, error) {
	var triger model.Trigger
	if err := db.QueryRow("SELECT id, name, type, schedule_time, interval_seconds, api_endpoint, api_payload, api_method, is_recurring, created_at, updated_at  FROM public.triggers WHERE id=$1", trigger_id).Scan(
		&triger.Id, &triger.Name, &triger.Type, &triger.ScheduleTime, &triger.Interval,
		&triger.ApiEndpoint, &triger.ApiPayload, &triger.ApiMethod, &triger.Recurring, &triger.CreatedAt, &triger.UpdatedAt,
	); err != nil {
		log.Print(err)
		{
			return model.Trigger{}, err
		}
	}

	return triger, nil
}

func DeleteTrigger(db *sql.DB, trigger_id string) error {
	_, err := db.Exec("DELETE FROM public.triggers WHERE id = $1", trigger_id)
	return err
}

func UpdateTrigger(db *sql.DB, trigger model.Trigger) error {
	stmt := `
		UPDATE public.triggers
		SET name = $1, type = $2, schedule_time = $3, interval_seconds = $4,
		    api_endpoint = $5, api_payload = $6, api_method = $7 is_recurring = $8, updated_at = $9
		WHERE id = $10
	`

	_, err := db.Exec(stmt, trigger.Name, trigger.Type, trigger.ScheduleTime,
		trigger.Interval, trigger.ApiEndpoint, trigger.ApiPayload, trigger.ApiMethod, trigger.Recurring, trigger.UpdatedAt, trigger.Id)

	if err != nil {
		log.Printf("Error updating trigger: %v", err)
		return err
	}

	fmt.Println("Trigger updated successfully")
	return nil
}
