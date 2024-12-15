package handlers

import (
	"database/sql"
	"event_trigger/repo"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListEventsHandler godoc
// @Summary List events
// @Description Retrieves a list of events with filtering options
// @Tags Events
// @Accept json
// @Produce json
// @Param includeArchived query bool false "Include archived events only (default: false)"
// @Param includeAll query bool false "Include all events, both active and archived (default: false)"
// @Success 200 {object} model.Event "List of events"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /events/ [get]
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
