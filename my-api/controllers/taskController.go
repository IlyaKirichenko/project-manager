package controllers

import (
	"database/sql"
	"my-api/models"
	"net/http"
	"strconv"

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

func GetTaskByIDHandler(c *gin.Context) {

	idTask := c.Param("id")

	id, err := strconv.Atoi(idTask)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный id"})
		return
	}
	task, err := models.GetTaskByID(id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}
