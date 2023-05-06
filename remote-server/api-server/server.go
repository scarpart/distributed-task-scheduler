package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/gops/agent"
	"github.com/prometheus/client_golang/prometheus"
	db "github.com/scarpart/distributed-task-scheduler/remote-server/db/sqlc"
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
	router.GET("/node", server.GetAllNodes)
	router.GET("/stats", server.GetStats)

	// PROMETHEUS METRICS 
	cpuUsage := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cpu_usage",
			Help: "CPU usage of a remote server in percentage",
		},
		[]string{"endpoint"},
	)
	memUsage := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mem_usage",
			Help: "Memory usage of a remote server in percentage",
		},
		[]string{"endpoint"},
	)
	prometheus.MustRegister(cpuUsage, memUsage)

	// PROMETHEUS ROUTES
	router.GET("/metrics", func(ctx *gin.Context) {
		cpu, mem, err := getUsageStats()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		cpuUsage.Set(cpu)
		memUsage.Set(mem)	
		prometheus.Handler().ServeHTTP(ctx.Writer, ctx.Request)
	})

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

func getUsageStats() (float64, float64, error) {
	pid := agent.ProcessID()

	cpu, err := agent.CPUPercent(time.Second)
	if err != nil {
		return 0, 0, err
	}

	mem, err := agent.MemoryUsage(pid)
	if err != nil {
		return 0, 0, err
	}

	return cpu, mem, nil
}
