package routes

import (
	"event_trigger/handlers"
	"event_trigger/middleware"

	"github.com/gin-gonic/gin"
)

func TriggerRoutes(r *gin.Engine, v string) {
	u := r.Group(v + "/triggers")
	u.Use(middleware.PostgresMiddleware())
	u.POST("/create", handlers.CreateTrigger)
	u.GET("/", handlers.ListTriggers)
	u.GET("/:trigger_id", handlers.GetTriggerByID)
	u.DELETE("/:trigger_id", handlers.DeleteTrigger)
	u.PUT("/:trigger_id", handlers.EditTrigger)
}
