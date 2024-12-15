package handlers

import (
	"database/sql"
	"event_trigger/model"
	"event_trigger/repo"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
		if req.ApiEndpoint == "" || req.ApiPayload == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "API Endpoint and Payload are required for API Trigger",
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
	now := time.Now()

	if req.Type == "API" {
		// Set fields for API Trigger
		trigger.ApiEndpoint = req.ApiEndpoint
		trigger.ApiPayload = req.ApiPayload
		trigger.Interval = nil
		trigger.ScheduleTime = nil
	} else if req.Type == "SCHEDULED" {
		// Set fields for Scheduled Trigger
		scheduledTime := now.Add(time.Duration(req.ScheduleTime) * time.Second)
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
		if req.ApiEndpoint != nil && *req.ApiEndpoint == "" {
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
		scheduledTime := time.Now().Add(time.Duration(*req.ScheduleTime) * time.Second)
		trigger.ScheduleTime = &scheduledTime
	}
	if req.Interval != nil {
		trigger.Interval = req.Interval
	}
	if req.ApiEndpoint != nil {
		trigger.ApiEndpoint = *req.ApiEndpoint
	}
	if req.ApiPayload != nil {
		trigger.ApiPayload = *req.ApiPayload
	}
	if req.Recurring != nil {
		trigger.Recurring = *req.Recurring
	}
	trigger.UpdatedAt = time.Now() // Update the `updated_at` field

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
