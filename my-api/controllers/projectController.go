package controllers

import (
	"my-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateProjectRequest struct {
	IdOwner     int    `json:"id_owner" binding:"required"`
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

	createNewProject, err := models.CreateProject(newProject.IdOwner, newProject.NameProject, newProject.Description, newProject.Deadline)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createNewProject)
}
