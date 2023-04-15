package main

import (
	"github.com/gin-gonic/gin"
	"github.com/scarpart/distributed-task-scheduler/src/controllers"
	"github.com/scarpart/distributed-task-scheduler/src/models"
)

func main() {
	r := gin.Default()

	models.ConnectDatabase()

	r.GET("/tasks", controllers.FindTasks)
	r.POST("/tasks", controllers.AddTask)

	r.Run()
}
