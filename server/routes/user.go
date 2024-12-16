package routes

import (
	"event_trigger/handlers"
	"event_trigger/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine, v string) {
	u := r.Group(v + "/users")
	u.Use(middleware.PostgresMiddleware())
	u.GET("/", handlers.ListUsersHandler)
	u.GET("/:user_id", handlers.GetUserByIDHandler)
	u.POST("/create", handlers.CreateUserHandler)
	u.PUT("/:user_id", handlers.EditUserHandler)
	u.DELETE("/:user_id", handlers.DeleteUserHandler)
}
