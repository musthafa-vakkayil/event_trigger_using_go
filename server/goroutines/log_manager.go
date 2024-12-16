package goroutines

import (
	"database/sql"
	"log"
	"time"
)

// This function will change the log status to atrchive after 2 hours
func ArchiveLogs(db *sql.DB, archiveChannel chan<- string) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rows, err := db.Query(
			"SELECT id FROM events WHERE status = 'ACTIVE' AND event_time <= NOW() - INTERVAL '2 Hour'",
		)
		if err != nil {
			log.Printf("Error fetching logs for archiving: %v", err)
			continue
		}
		defer rows.Close()

		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err != nil {
				log.Printf("Failed to scan log ID: %v", err)
				continue
			}
			archiveChannel <- id
		}
	}
}

// Update log status to Archive
func ProcessArchiving(db *sql.DB, archiveChannel <-chan string) {
	for logID := range archiveChannel {
		_, err := db.Exec(
			"UPDATE events SET status = 'ARCHIVE', archived_time = NOW() WHERE id = $1", logID,
		)
		if err != nil {
			log.Printf("Error archiving log: %v", err)
		} else {
			log.Printf("Log ID %s archived", logID)
		}
	}
}

// function to Delete The logs after Archiving
func DeleteLogs(db *sql.DB) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		_, err := db.Exec(
			"DELETE FROM events WHERE status = 'ARCHIVE' AND archived_time <= NOW() - INTERVAL '46 Hour'",
		)
		if err != nil {
			log.Printf("Error deleting logs: %v", err)
		}
	}
}
