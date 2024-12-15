package goroutines

import (
	"database/sql"
	"event_trigger/utils"
	"fmt"
	"log"
	"time"
)

func ProcessTriggers(triggerChannel <-chan string, db *sql.DB) {
	for triggerID := range triggerChannel {
		go func(id string) {
			log.Printf("Processing trigger ID: %s", id)

			// Fetch trigger details
			var scheduleTime time.Time
			var intervalSeconds int
			var isRecurring bool
			err := db.QueryRow("SELECT schedule_time, interval_seconds, is_recurring FROM triggers WHERE id = $1", id).
				Scan(&scheduleTime, &intervalSeconds, &isRecurring)
			if err != nil {
				log.Printf("Failed to fetch trigger details for ID %s: %v", id, err)
				return
			}

			// Ensure schedule time is in UTC
			scheduleTime = scheduleTime.UTC()
			currentTime := time.Now().UTC()

			fmt.Println("Schedule TIme is", scheduleTime)
			fmt.Println(currentTime)

			// Calculate sleep duration
			sleepDuration := time.Until(scheduleTime)
			if sleepDuration > 0 {
				log.Printf("Sleeping for %v to match schedule time", sleepDuration)
				time.Sleep(sleepDuration)
			}

			// Execute trigger
			log.Printf("Executing trigger ID: %s", id)
			_, err = db.Exec("INSERT INTO events (id, trigger_id, event_time, status) VALUES ($1, $2, $3, $4)",
				utils.GenerateID(), id, time.Now().UTC(), "ACTIVE")
			if err != nil {
				log.Printf("Failed to log event for trigger ID %s: %v", id, err)
				return
			}

			// Reschedule if recurring
			if isRecurring {
				nextExecution := scheduleTime.Add(time.Second * time.Duration(intervalSeconds))
				_, err = db.Exec("UPDATE triggers SET schedule_time = $1 WHERE id = $2", nextExecution, id)
				if err != nil {
					log.Printf("Failed to reschedule recurring trigger ID %s: %v", id, err)
				} else {
					log.Printf("Trigger ID %s rescheduled for %s", id, nextExecution)
				}
			}
		}(triggerID)
	}
}
