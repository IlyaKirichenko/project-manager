package controllers

import (
	"my-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateTaskRequest struct {
	IdProject   int    `json:"id_project" binding:"required"`
	IdUser      int    `json:"id_user" binding:"required"`
	NameTask    string `json:"name_task" binding:"required"`
	IdState     int    `json:"id_state" binding:"required"`
	Description string `json:"description"`
}

func CreateTaskHandler(c *gin.Context) {
	var newTask CreateTaskRequest

	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createNewTask, err := models.CreateTask(newTask.IdProject, newTask.IdUser, newTask.IdState, newTask.NameTask, newTask.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createNewTask)
}
