package goroutines

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/lib/pq"
)

// This function will be actively listening the the database for any incoming trigger
// It will Skip the API Trigger
func ListenForTriggers(db *sql.DB, triggerChannel chan<- string) {
	dbString := os.Getenv("DB_CONNECTION")
	listener := pq.NewListener(dbString,
		10*time.Second, time.Minute, func(ev pq.ListenerEventType, err error) {
			if err != nil {
				log.Printf("Listener error: %v", err)
			}
		})

	err := listener.Listen("trigger_channel")
	if err != nil {
		log.Fatalf("Error listening to trigger channel: %v", err)
	}

	log.Println("Listening for triggers...")

	for {
		select {
		case n := <-listener.Notify:
			if n != nil {
				triggerID := n.Extra

				// Check trigger type before adding to the processing channel
				var triggerType string
				err := db.QueryRow("SELECT type FROM triggers WHERE id = $1", triggerID).Scan(&triggerType)
				if err != nil {
					log.Printf("Failed to fetch trigger type for ID %s: %v", triggerID, err)
					continue
				}

				if triggerType == "API" {
					log.Printf("Trigger ID %s is of type 'API'. Ignoring.", triggerID)
					continue
				}

				log.Printf("Trigger received: %s", triggerID)
				triggerChannel <- triggerID
			}
		}
	}
}
