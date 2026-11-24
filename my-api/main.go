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

	router.Static("/static", "./static")
	router.StaticFile("/", "./static/Login/login.html")

	router.POST("/user", controllers.CreateUserHandler)
	router.POST("/login", controllers.LoginHandler)

	protected := router.Group("/")
	protected.Use(middleware.RequireAuth())
	{
		protected.POST("/logout", controllers.LogoutHandler)
		protected.GET("/me", controllers.MeHandler)
		protected.POST("/project", controllers.CreateProjectHandler)
		protected.GET("/project", controllers.GetUserProjectsHandler)
		protected.POST("/tasks", controllers.CreateTaskHandler)
		protected.GET("/tasks/:id", controllers.GetTaskByIDHandler)
	}

	router.Run(":8000")
}
