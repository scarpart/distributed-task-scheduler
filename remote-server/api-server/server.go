package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/scarpart/distributed-task-scheduler/remote-server/db/sqlc"
)

type Server struct {
	addr    string
	store   *db.Store
	router  *gin.Engine
}

// Constructs the server and sets up the routing
func NewServer(store *db.Store, addr string) *Server {
	server := &Server{store: store, addr: addr}
	router := gin.Default()

	// Health Check (accessed by the Load Balancer)
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, nil)
	})
	
	// POST
	router.POST("/task", server.CreateTask)
	router.POST("/node", server.CreateNode)
	router.POST("/mapping", server.CreateMapping)

	// GET 
	//router.GET("/task", server.GetTask)
	router.GET("/node/:node_id", server.GetNode)
	router.GET("/node", server.GetAllNodes)

	server.router = router
	return server
}

// Run the HTTP server on the input address to listen to requests
func (server *Server) Start() error {
	return server.router.Run(server.addr)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}

