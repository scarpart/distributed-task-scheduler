package apiserver

import (
	"crypto/tls"
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

	// In a real application, there would be no "public" group, 
	// as it would only be accessible for people that paid for the service,
	// but this is simple enough for the project, I don't want to overcomplicate
	// things, since the core of it is the actual distribution of requests.
	publicGroup := router.Group("/public")
	{
		publicGroup.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, nil)
		})
		//publicGroup.POST("/user", server.CreateUser)
	}

	// Private Group uses API Key validation as custom a middleware 
	privateGroup := router.Group("/private")
	privateGroup.Use(server.ValidateKeys())
	{
		// POST
		privateGroup.POST("/task", server.CreateTask)
		privateGroup.POST("/node", server.CreateNode)
		privateGroup.POST("/mapping", server.CreateMapping)

		// GET
		//privateGroup.GET("/task", server.GetTask)
		privateGroup.GET("/node/:node_id", server.GetNode)
		privateGroup.GET("/node", server.GetAllNodes)
	}

	server.router = router
	return server
}

// Run the HTTPS server on the input address to listen to requests
func (server *Server) Start(certFile, keyFile string) error {
	// Configuring TLS such that we have secure HTTPS connections
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
		},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		},
		PreferServerCipherSuites: true,
	}

	httpsServer := &http.Server{
		Addr: server.addr,
		Handler: server.router,
		TLSConfig: tlsConfig,
	}

	return httpsServer.ListenAndServeTLS(certFile, keyFile)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}

