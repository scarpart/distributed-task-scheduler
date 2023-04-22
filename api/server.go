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

	// POST
	router.POST("/task", server.CreateTask)
	router.POST("/node", server.CreateNode)
	router.POST("/mapping", server.CreateMapping)

	// GET 
	//router.GET("/task", server.GetTask)
	router.GET("/node/:node_id", server.GetNode)

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
