package routes

import (
	"event_trigger/handlers"
	"event_trigger/middleware"

	"github.com/gin-gonic/gin"
)

func TestTriggerRoutes(r *gin.Engine, v string) {
	u := r.Group(v + "/test")
	u.Use(middleware.PostgresMiddleware())
	u.POST("/trigger", handlers.TestTriggerHandler)
}
