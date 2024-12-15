package handlers

import (
	"database/sql"
	"event_trigger/model"
	"event_trigger/repo"
	"event_trigger/utils"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateTrigger godoc
// @Summary Create a new Trigger. API or Scheduled
// @Description Adds a new Trigger to the system
// @Tags Triggers
// @Accept json
// @Produce json
// @Param user body model.TriggerDto true "Trigger details"
// @Success 200 {string} string "New Trigger Created with ID"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /triggers/create [post]
func CreateTrigger(c *gin.Context) {
	var req model.TriggerDto

	// Parse and validate the input payload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad payload provided",
		})
		return
	}

	// Validate Trigger Type
	if req.Type != "API" && req.Type != "SCHEDULED" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Trigger Type",
		})
		return
	}

	// Validation for API Trigger
	if req.Type == "API" {
		if req.ApiEndpoint == "" || req.ApiMethod == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "API Endpoint and Method are required for API Trigger",
			})
			return
		}
	}

	// Validation for Scheduled Trigger
	if req.Type == "SCHEDULED" {
		if req.ScheduleTime == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Schedule Time is required for SCHEDULED Trigger",
			})
			return
		}
	}

	// Get PostgreSQL client
	pgClient := c.MustGet("postgresClient").(*sql.DB)

	// Initialize the Trigger model
	var trigger model.Trigger
	trigger.Id = uuid.New().String()
	trigger.Name = req.Name
	trigger.Type = req.Type
	trigger.Recurring = req.Recurring

	if req.Type == "API" {
		// Set fields for API Trigger
		trigger.ApiEndpoint = req.ApiEndpoint
		trigger.ApiPayload = req.ApiPayload
		trigger.ApiMethod = req.ApiMethod
		trigger.Interval = 0
		trigger.ScheduleTime = nil
	} else if req.Type == "SCHEDULED" {
		// Set fields for Scheduled Trigger
		scheduledTime := time.Now().UTC().Add(time.Duration(req.ScheduleTime) * time.Second)
		trigger.ScheduleTime = &scheduledTime
		trigger.ApiEndpoint = ""
		trigger.ApiPayload = nil
		trigger.Interval = req.Interval
	}

	// Insert the Trigger into the database
	id, err := repo.CreateTrigger(pgClient, trigger)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create Trigger",
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message":    "success",
		"trigger_id": id,
	})
}

// DeleteTrigger godoc
// @Summary Delete a trigger by ID
// @Description Removes a trigger from the system
// @Tags Triggers
// @Param trigger_id path string true "Trigger ID"
// @Success 200 {string} string "Trigger deleted successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /trigegrs/{trigger_id} [delete]
func DeleteTrigger(c *gin.Context) {
	id := c.Param("trigger_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "trigger_id query param is required",
		})
		return
	}

	pgClient := c.MustGet("postgresClient").(*sql.DB)

	if err := repo.DeleteTrigger(pgClient, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to delete Trigger",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Trigger with id %s Deleted\n", id),
	})
}

// GetTriggerByID godoc
// @Summary Get trigger by ID
// @Description Retrieves a trigger based on the given ID
// @Tags Triggers
// @Param trigger_id path string true "Trigger ID"
// @Produce json
// @Success 200 {object} model.Trigger "Trigger data"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /triggers/{trigger_id} [get]
func GetTriggerByID(c *gin.Context) {
	id := c.Param("trigger_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "trigger_id query param is required",
		})
		return
	}

	pgClient := c.MustGet("postgresClient").(*sql.DB)

	msgdata, err := repo.GetTriggerByID(pgClient, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to get Trigger",
		})
		return
	}

	c.JSON(http.StatusOK, msgdata)
}

// ListTriggers godoc
// @Summary List all trigegrs
// @Description Retrieves a list of all triggers
// @Tags Triggers
// @Produce json
// @Success 200 {array} model.Trigger "List of triggers"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /triggers/ [get]
func ListTriggers(c *gin.Context) {

	pgClient := c.MustGet("postgresClient").(*sql.DB)

	usrs, err := repo.ListTriggers(pgClient)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to get triggers",
		})
		return
	}

	c.JSON(http.StatusOK, usrs)
}

// EditTrigegr godoc
// @Summary Edit trigger details
// @Description Updates the information of a trigger
// @Tags Triggers
// @Param trigger_id path string true "Trigegr ID"
// @Param user body model.EditTriggerDto true "Updated trigger details"
// @Success 200 {string} string "Trigger updated successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /triggers/{trigger_id} [put]
func EditTrigger(c *gin.Context) {
	var req model.EditTriggerDto

	// Parse query parameter `trigger_id`
	triggerID := c.Param("trigger_id")
	if triggerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing query parameter 'trigger_id'",
		})
		return
	}

	// Parse and validate the input payload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid payload provided",
		})
		return
	}

	// Validate Trigger Type if provided
	if req.Type != nil && *req.Type != "API" && *req.Type != "SCHEDULED" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Trigger Type",
		})
		return
	}

	// Validation for API Trigger
	if req.Type != nil && *req.Type == "API" {
		if req.ApiEndpoint != nil && *req.ApiMethod == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "API Endpoint cannot be empty for API Trigger",
			})
			return
		}
	}

	// Validation for Scheduled Trigger
	if req.Type != nil && *req.Type == "SCHEDULED" {
		if req.ScheduleTime != nil && *req.ScheduleTime == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Schedule Time cannot be zero for SCHEDULED Trigger",
			})
			return
		}
	}

	// Get PostgreSQL client
	pgClient := c.MustGet("postgresClient").(*sql.DB)

	// Fetch the existing trigger from the database
	trigger, err := repo.GetTriggerByID(pgClient, triggerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Trigger not found",
		})
		return
	}

	// Update fields if provided in the request
	if req.Name != nil {
		trigger.Name = *req.Name
	}
	if req.Type != nil {
		trigger.Type = *req.Type
	}
	if req.ScheduleTime != nil {
		scheduledTime := time.Now().UTC().Add(time.Duration(*req.ScheduleTime) * time.Second)
		trigger.ScheduleTime = &scheduledTime
	}
	if req.Interval != 0 {
		trigger.Interval = req.Interval
	}
	if req.ApiEndpoint != nil {
		trigger.ApiEndpoint = *req.ApiEndpoint
	}
	if req.ApiPayload != nil {
		trigger.ApiPayload = *req.ApiPayload
	}
	if req.ApiMethod != nil {
		trigger.ApiMethod = *req.ApiMethod
	}
	if req.Recurring != nil {
		trigger.Recurring = *req.Recurring
	}

	// Update the `updated_at` field
	trigger.UpdatedAt = time.Now().UTC()

	// Update the trigger in the database
	if err := repo.UpdateTrigger(pgClient, trigger); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update Trigger",
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Trigger updated successfully",
	})
}

// TriggerApiByID godoc
// @Summary Trigger API By ID
// @Description triggers a stored api trigger in the db
// @Tags Events
// @Param trigger_id path string true "Trigger ID"
// @Produce json
// @Success 200 {object} map[string]string "Response from API"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /triggers/api/{trigger_id} [get]
func TriggerAPI(c *gin.Context) {
	id := c.Param("trigger_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "trigger_id query param is required",
		})
		return
	}

	pgClient := c.MustGet("postgresClient").(*sql.DB)

	trigger, err := repo.GetTriggerByID(pgClient, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Trigger not found",
		})
		return
	}

	if trigger.Type != "API" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "This endpoint is for triggering API triggers",
		})
		return
	}

	// Use the utility function to make the API call
	payload := []byte(trigger.ApiPayload)
	resp, err := utils.MakeAPICall(trigger.ApiMethod, trigger.ApiEndpoint, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to execute API trigger: " + err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read HTTP response",
		})
		return
	}

	// Create Event
	event := model.Event{
		Id:        uuid.New().String(),
		TriggerId: trigger.Id,
		EventTime: time.Now().UTC(),
		Status:    "ACTIVE",
		Manual:    true,
	}

	_, evErr := repo.CreateEvent(pgClient, event)
	if evErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error While Creating Event Log",
		})
		return
	}

	// Return the API response
	c.JSON(http.StatusOK, gin.H{
		"message":  "API trigger executed successfully",
		"response": string(body),
		"status":   resp.StatusCode,
	})
}
