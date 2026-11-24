package controllers

import (
	"my-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateProjectRequest struct {
	NameProject string `json:"name_project" binding:"required"`
	Description string `json:"description_project"`
	Deadline    string `json:"deadline" binding:"required"`
}

func CreateProjectHandler(c *gin.Context) {
	var newProject CreateProjectRequest

	if err := c.ShouldBindJSON(&newProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(int)

	createNewProject, err := models.CreateProject(userID, newProject.NameProject, newProject.Description, newProject.Deadline)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createNewProject)
}

func GetUserProjectsHandler(c *gin.Context) {
	userID, ok := c.MustGet("userID").(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Вы не авторизованы"})
		return
	}
	projects, err := models.GetProjectsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projects)
}
