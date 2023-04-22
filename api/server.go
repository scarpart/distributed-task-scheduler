package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/scarpart/distributed-task-scheduler/db/sqlc"
)

type Server struct {
	store *db.Store
	router *gin.Engine
}

// Constructs the server and sets up the routing
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/task", server.CreateTask)
	router.POST("/node", server.CreateNode)

	server.router = router
	return server
}

// Run the HTTP server on the input address to listen to requests
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
