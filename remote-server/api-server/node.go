package apiserver

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/scarpart/distributed-task-scheduler/remote-server/db/sqlc"
)

type CreateNodeRequest struct {
	Hostname string `json:"hostname" binding:"required"`
	IpAddr   string `json:"ip_addr" binding:"required"`
}

func (server *Server) CreateNode(ctx *gin.Context) {
	var req CreateNodeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateNodeParams{
		Hostname: req.Hostname,
		IpAddr:   req.IpAddr,
		Status:   0,
	}

	node, err := server.store.CreateNode(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, node)
}

func (server *Server) GetNode(ctx *gin.Context) {
	nodeId, err := strconv.Atoi(ctx.Param("node_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	node, err := server.store.GetNode(ctx, int64(nodeId))
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, node)
}

type GetAllNodesRequest struct {
	Limit int32 `json:"limit" binding:"required"`
	Offset int32 `json:"offset" binding:"required"`
}

func (server *Server) GetAllNodes(ctx *gin.Context) {
	var req GetAllNodesRequest 
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetAllNodesParams{
		Limit: req.Limit,
		Offset: req.Offset,
	}
	
	nodes, err := server.store.GetAllNodes(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nodes)
}
