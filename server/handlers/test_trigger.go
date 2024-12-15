package handlers

import (
	"database/sql"
	"event_trigger/model"
	"event_trigger/repo"
	"event_trigger/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TestTriggerWithoutSaving godoc
// @Summary excute a trigger without saving
// @Description Executes API trigegr or One Time Scheduled Trigger Without Saving, but can see in logs
// @Tags Test
// @Accept json
// @Produce json
// @Param user body model.TriggerDto true "Trigger details"
// @Success 200 {string} string "API response / Trigger Scheduled Successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /test/trigger [post]
func TestTriggerHandler(c *gin.Context) {
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

	testTriggerId := utils.GenerateID()

	var response string
	var err error
	if req.Type == "API" {
		// trigger api and return response
		response, err = handleTestApiTrigger(c, testTriggerId, req.ApiMethod, req.ApiEndpoint, []byte(req.ApiPayload))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
	} else {
		// start a go routine and return response
		go handleTestScheduledTrigger(c, testTriggerId, req.ScheduleTime)
		response = "Trigger Scheduled Successfully"
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": response,
	})
}

func handleTestApiTrigger(c *gin.Context, trigger_id string, ApiMethod string, ApiEndpoint string, payload []byte) (string, error) {
	resp, err := utils.MakeAPICall(ApiMethod, ApiEndpoint, payload)
	if err != nil {
		return "", fmt.Errorf("failed to execute API trigger: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read HTTP response")
	}

	// Create Event
	event := model.Event{
		Id:        uuid.New().String(),
		TriggerId: trigger_id,
		EventTime: time.Now(),
		Status:    "ACTIVE",
		Manual:    true,
	}
	pgClient := c.MustGet("postgresClient").(*sql.DB)

	_, evErr := repo.CreateEvent(pgClient, event)
	if evErr != nil {
		return "", fmt.Errorf("error While Creating Event Log")
	}

	return string(body), nil
}

func handleTestScheduledTrigger(c *gin.Context, trigger_id string, scheduleTime int) {

	// Ensure schedule time is in UTC
	actualscheduleTime := time.Now().UTC().Add(time.Duration(scheduleTime) * time.Second)

	// Calculate sleep duration
	sleepDuration := time.Until(actualscheduleTime)
	if sleepDuration > 0 {
		log.Printf("Sleeping for %v to match schedule time", sleepDuration)
		time.Sleep(sleepDuration)
	}

	log.Printf("Executing trigger ID: %s", trigger_id)

	// Create Event
	event := model.Event{
		Id:        uuid.New().String(),
		TriggerId: trigger_id,
		EventTime: time.Now(),
		Status:    "ACTIVE",
		Manual:    true,
	}
	pgClient := c.MustGet("postgresClient").(*sql.DB)

	_, evErr := repo.CreateEvent(pgClient, event)
	if evErr != nil {
		log.Print("Error While Creating Event Log")
	}

}
