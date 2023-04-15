package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scarpart/distributed-task-scheduler/src/models"
	"github.com/scarpart/distributed-task-scheduler/src/types"
)

type NewTask struct {
	TaskName        string            `json:"taskName" binding="required"`
	TaskDescription int               `json:"taskDescription" binding="required"`
	Status          int               `json:"status" binding="required"`
	Priority        int               `json:"priority" binding="required"`
	Dependencies    types.Uint64Array `json:"dependencies" binding="required"`
	NodeID          int               `json:"nodeId" binding="required"`
	Command         string            `json:"command" binding="required"`
}

type TaskUpdate struct {
	TaskName        string            `json:"taskName"`
	TaskDescription int               `json:"taskDescription"`
	Status          int               `json:"status"`
	Priority        int               `json:"priority"`
	Dependencies    types.Uint64Array `json:"dependencies"`
	NodeID          int               `json:"nodeId"`
	Command         string            `json:"command"`
}

// GET /tasks
// GET ALL tasks
func FindTasks(c *gin.Context) {
	var tasks []models.Task
	models.DB.Find(&tasks)
	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

// POST /tasks
// CREATE NEW task
func AddTask(c *gin.Context) {
	var input NewTask
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := models.Task{
		TaskName:        input.TaskName,
		TaskDescription: input.TaskDescription,
		Status:          input.Status,
		Priority:        input.Priority,
		Dependencies:    input.Dependencies,
		NodeID:          input.NodeID,
		Command:         input.Command,
	}
	models.DB.Create(&task)

	c.JSON(http.StatusOK, gin.H{"data": task})
}
