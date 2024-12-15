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

func CreateUser(c *gin.Context) {
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
		"message": fmt.Sprintf("New User Created with id %d\n", id),
	})
}

func DeleteUser(c *gin.Context) {
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

func GetUserByID(c *gin.Context) {
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

func ListUsers(c *gin.Context) {

	pgClient := c.MustGet("postgresClient").(*sql.DB)

	usrs, err := repo.ListUsers(pgClient)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to get messages",
		})
		return
	}

	c.JSON(http.StatusOK, usrs)
}

func EditUser(c *gin.Context) {
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
