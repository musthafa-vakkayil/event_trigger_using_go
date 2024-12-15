package routes

import (
	"event_trigger/handlers"
	"event_trigger/middleware"

	"github.com/gin-gonic/gin"
)

func EventRoutes(r *gin.Engine, v string) {
	u := r.Group(v + "/events")
	u.Use(middleware.PostgresMiddleware())
	u.GET("/", handlers.ListEventsHandler)
}
