package routes

import (
	"event_trigger/handlers"
	"event_trigger/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine, v string) {
	u := r.Group(v + "/users")
	u.Use(middleware.PostgresMiddleware())
	u.GET("/", handlers.ListUsers)
	u.GET("/:user_id", handlers.GetUserByID)
	u.POST("/create", handlers.CreateUser)
	u.PUT("/:user_id", handlers.EditUser)
	u.DELETE("/:user_id", handlers.DeleteUser)
}
