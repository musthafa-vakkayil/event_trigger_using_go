package routes

import (
	"event_trigger/handlers"
	"event_trigger/middleware"

	"github.com/gin-gonic/gin"
)

func HealthRoutes(r *gin.Engine, v string) {
	u := r.Group(v + "/health")
	u.Use(middleware.PostgresMiddleware())
	u.GET("/", handlers.HealthCheck)
}
