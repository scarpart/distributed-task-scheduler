package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/scarpart/distributed-task-scheduler/remote-server/db/sqlc"
)

type CreateMappingRequest struct {
	TaskID int64 `json:"task_id"`
	NodeID int64 `json:"node_id"`
}

// Binds a Node to a Task and vice versa
func (server *Server) CreateMapping(ctx *gin.Context) {
	var req CreateMappingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateTaskNodeMappingParams{
		TaskID: req.TaskID,
		NodeID: req.NodeID,
	}

	mapping, err := server.store.CreateTaskNodeMapping(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, mapping)
}

