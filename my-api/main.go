package main

import (
	"my-api/middleware"
	"my-api/models"

	"my-api/controllers"

	"github.com/gin-gonic/gin"
)

func main() {

	models.InitDB()
	models.InitSchema()
	models.SeedStates()

	router := gin.Default()

	router.POST("/user", controllers.CreateUserHandler)
	router.POST("/login", controllers.LoginHandler)

	protected := router.Group("/")
	protected.Use(middleware.RequireAuth())
	{
		router.POST("/logout", controllers.LogoutHandler)
		router.GET("/me", controllers.MeHandler)
		router.POST("/project", controllers.CreateProjectHandler)
		router.POST("/tasks", controllers.CreateTaskHandler)
		router.GET("/tasks/:id", controllers.GetTaskByIDHandler)
	}

	router.Run(":8000")
}
