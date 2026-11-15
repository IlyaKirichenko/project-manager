package controllers

import (
	"my-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func CreateUserHandler(c *gin.Context) {
	var newUser CreateUserRequest

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createNewUser, err := models.CreateUser(newUser.UserName, newUser.Email, newUser.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createNewUser)

}
