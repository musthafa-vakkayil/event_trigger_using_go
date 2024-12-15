package handlers

import (
	"database/sql"
	"event_trigger/repo"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListEventsHandler(c *gin.Context) {
	pgClient := c.MustGet("postgresClient").(*sql.DB)

	// Parse query parameters
	includeArchived := c.Query("includeArchived") == "true"
	includeAll := c.Query("includeAll") == "true"

	// Fetch events from the database
	events, err := repo.ListEvents(pgClient, includeArchived, includeAll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch events: " + err.Error(),
		})
		return
	}

	// Return the list of events
	c.JSON(http.StatusOK, gin.H{
		"events": events,
	})
}
