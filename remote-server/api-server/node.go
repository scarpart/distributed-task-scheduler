package apiserver

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/scarpart/distributed-task-scheduler/remote-server/db/sqlc"
)

type CreateNodeRequest struct {
	Hostname string `json:"hostname" binding:"required"`
	IpAddr   string `json:"ip_addr" binding:"required"`
}

func (server *Server) CreateNode(ctx *gin.Context) {
	fmt.Println("On CreateNode")

	var req CreateNodeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Println("Could not bind json (1):", err)
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
		fmt.Println("Err (2):", err)
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
	fmt.Println("On GetAllNodes")

	// Parsing query parameters, since this is a GET request
	limitParam := ctx.Query("limit")
	if limitParam == "" {
		limitParam = "10"
	}
	offsetParam := ctx.Query("offset")
	if offsetParam == "" {
		offsetParam = "0"
	}

	limit, err := strconv.ParseInt(limitParam, 10, 32)
	if err != nil {
		fmt.Println("Invalid limit parameter", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	offset, err := strconv.ParseInt(offsetParam, 10, 32)
	if err != nil {
		fmt.Println("Invalid offset parameter", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fmt.Println("limit:", limit)
	fmt.Println("offset:", offset)

	arg := db.GetAllNodesParams{
		Limit: int32(limit), 
		Offset: int32(offset),
	}
	
	nodes, err := server.store.GetAllNodes(ctx, arg)
	if err != nil {
		fmt.Println("on second error check", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	fmt.Println(nodes)
	ctx.JSON(http.StatusOK, nodes)
}

