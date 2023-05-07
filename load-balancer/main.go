package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	requestdistribution "github.com/scarpart/distributed-task-scheduler/load-balancer/request-distribution"
	"github.com/scarpart/distributed-task-scheduler/load-balancer/server"
)

func main() {
	router := *gin.Default()

	sm := requestdistribution.ServerMonitor{
		Heap: requestdistribution.NewHeap(),
		Servers: [] server.RemoteServer{
			{
				IpAddr: "127.0.0.1:8081",
				BaseUrl: "localhost:8081",	
			},
			{
				IpAddr: "127.0.0.1:8080",
				BaseUrl: "localhost:8080",	
			},
		},
	}

	router.GET("/test", sm.WatchServers)
}
