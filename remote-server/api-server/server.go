package apiserver

import (
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	db "github.com/scarpart/distributed-task-scheduler/remote-server/db/sqlc"
	"github.com/shirou/gopsutil/process"
)

type Server struct {
	ipAddr net.IP
	store *db.Store
	router *gin.Engine
}

// Constructs the server and sets up the routing
func NewServer(store *db.Store, ipAddr net.IP) *Server {
	server := &Server{store: store, ipAddr: ipAddr}
	router := gin.Default()

	// POST
	router.POST("/task", server.CreateTask)
	router.POST("/node", server.CreateNode)
	router.POST("/mapping", server.CreateMapping)

	// GET 
	//router.GET("/task", server.GetTask)
	router.GET("/node/:node_id", server.GetNode)
	router.GET("/node", server.GetAllNodes)

	// Prometheus metrics to be read from the main server (Load Balancer)
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

	// GET - Prometheus 
	router.GET("/metrics", func(ctx *gin.Context) {
		cpu, mem, err := getUsageStats()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		cpuUsage.With(prometheus.Labels{"endpoint": server.ipAddr.String()}).Set(cpu)
		memUsage.With(prometheus.Labels{"endpoint": server.ipAddr.String()}).Set(mem)	
		promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
	})

	server.router = router
	return server
}

// Run the HTTP server on the input address to listen to requests
func (server *Server) Start() error {
	return server.router.Run(server.ipAddr.String())
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}

func getUsageStats() (float64, float64, error) {
	proc, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		panic("aaa")
	}

	cpu, err := proc.CPUPercent()
	if err != nil {
		return 0, 0, err
	}

	mem, err := proc.MemoryPercent()
	if err != nil {
		return 0, 0, err
	}

	return cpu, float64(mem), nil
}
