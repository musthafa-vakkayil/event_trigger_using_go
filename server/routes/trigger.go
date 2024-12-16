package routes

import (
	"event_trigger/handlers"
	"event_trigger/middleware"

	"github.com/gin-gonic/gin"
)

func TriggerRoutes(r *gin.Engine, v string) {
	u := r.Group(v + "/triggers")
	u.Use(middleware.PostgresMiddleware())
	u.POST("/create", handlers.CreateTriggerHandler)
	u.GET("/", handlers.ListTriggersHandler)
	u.GET("/:trigger_id", handlers.GetTriggerByIDHandler)
	u.DELETE("/:trigger_id", handlers.DeleteTriggerHandler)
	u.PUT("/:trigger_id", handlers.EditTriggerHandler)
	u.GET("/api/:trigger_id", handlers.TriggerAPIHandler)
}
