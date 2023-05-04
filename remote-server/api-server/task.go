package apiserver

import (
	"database/sql"
	"net/http"
	
	db "github.com/scarpart/distributed-task-scheduler/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/scarpart/distributed-task-scheduler/api/enums"
)

type CreateTaskRequest struct {
	TaskName        string        `json:"task_name" binding:"required"`
	TaskDescription string        `json:"task_description" binding:"required"`
	// `Status` should be "Ready" at first, so there's no point in having this as an input
	// Status          int32         `json:"status"`
	Priority        sql.NullInt32 `json:"priority"`
	Command         string        `json:"command" binding:"required"`
}

func (server *Server) CreateTask(ctx *gin.Context) {
	var req CreateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateTaskParams{
		TaskName: req.TaskName,
		TaskDescription: req.TaskDescription,
		Status: enums.Ready, 
		Priority: req.Priority,
		Command: req.Command,
	}

	task, err := server.store.CreateTask(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, task)
}
