package handlers

import (
	"database/sql"
	"event_trigger/model"
	"event_trigger/repo"
	"event_trigger/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateUser godoc
// @Summary Create a new user
// @Description Adds a new user to the system
// @Tags Users
// @Accept json
// @Produce json
// @Param user body model.UserDto true "User details"
// @Success 200 {string} string "New User Created with ID"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/create [post]
func CreateUserHandler(c *gin.Context) {
	var req model.User

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad payload pprovided",
		})
		return
	}

	pgClient := c.MustGet("postgresClient").(*sql.DB)

	// Hash the password using bcrypt
	passWordHash, err := utils.GeneratePasswordHash(req.PasswordHash)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	usr := model.User{
		UserId:       uuid.New().String(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passWordHash,
	}

	id, err := repo.CreateUser(pgClient, usr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("New User Created with id %s\n", id),
	})
}

// DeleteUser godoc
// @Summary Delete a user by ID
// @Description Removes a user from the system
// @Tags Users
// @Param user_id path string true "User ID"
// @Success 200 {string} string "User deleted successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{user_id} [delete]
func DeleteUserHandler(c *gin.Context) {
	id := c.Param("user_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id query param is required",
		})
		return
	}

	pgClient := c.MustGet("postgresClient").(*sql.DB)

	if err := repo.DeleteUser(pgClient, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to delete user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("User with id %s Deleted\n", id),
	})
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Retrieves a user based on the given ID
// @Tags Users
// @Param user_id path string true "User ID"
// @Produce json
// @Success 200 {object} model.User "User data"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{user_id} [get]
func GetUserByIDHandler(c *gin.Context) {
	id := c.Param("user_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id query param is required",
		})
		return
	}

	pgClient := c.MustGet("postgresClient").(*sql.DB)

	msgdata, err := repo.GetUserByID(pgClient, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to get user",
		})
		return
	}

	c.JSON(http.StatusOK, msgdata)
}

// ListUsers godoc
// @Summary List all users
// @Description Retrieves a list of all users
// @Tags Users
// @Produce json
// @Success 200 {array} model.User "List of users"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/ [get]
func ListUsersHandler(c *gin.Context) {

	pgClient := c.MustGet("postgresClient").(*sql.DB)

	usrs, err := repo.ListUsers(pgClient)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to get users",
		})
		return
	}

	c.JSON(http.StatusOK, usrs)
}

// EditUser godoc
// @Summary Edit user details
// @Description Updates the information of a user
// @Tags Users
// @Param user_id path string true "User ID"
// @Param user body model.UserDto true "Updated user details"
// @Success 200 {string} string "User updated successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{user_id} [put]
func EditUserHandler(c *gin.Context) {
	id := c.Param("user_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id query param is required",
		})
		return
	}
	var req model.User

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad payload pprovided",
		})
		return
	}

	pgClient := c.MustGet("postgresClient").(*sql.DB)

	// Hash the password using bcrypt
	passWordHash, err := utils.GeneratePasswordHash(req.PasswordHash)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	msg := model.User{
		UserId:       uuid.New().String(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passWordHash,
	}

	if err := repo.UpdateUser(pgClient, msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User Updated Succesfully",
	})
}
