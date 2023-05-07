package requestdistribution

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"github.com/scarpart/distributed-task-scheduler/load-balancer/server"
)

type ServerMonitor struct {
	Heap    Heap
	Servers []server.RemoteServer
}

func NewServerMonitor() ServerMonitor {
	return ServerMonitor{
		Heap: NewHeap(),
		Servers: []server.RemoteServer{},
	}
}

// TODO: Make it so ServerMonitor listens to the /metrics endpoint and updates
// the heap accordingly, in an interval of 5 seconds. There should be a method
// to return the root node to the main thread of the LoadBalancer, which will 
// then distribute the request to that remote server.

// I think Poll removes the node. Make a method to just get it. 
func (sm *ServerMonitor) GetLeastBusy() (*server.RemoteServer, error) {
	return sm.Heap.Poll()
}

func (sm *ServerMonitor) WatchServers(ctx *gin.Context, cpuGauge prometheus.GaugeVec, memGauge prometheus.GaugeVec) {
	timer := time.NewTicker(time.Second * 5)
	client := &http.Client{
		Timeout: time.Second * 10,
	}	

	for {
		select {
		case <-timer.C:
			sm.GetMetrics(ctx, client, cpuGauge, memGauge)					
			os.Exit(0)
		}
	}
}

func (sm *ServerMonitor) GetMetrics(ctx *gin.Context, client *http.Client, cpuGauge prometheus.GaugeVec, memGauge prometheus.GaugeVec) {
	for _, server := range sm.Servers {
		resp, err := client.Get(server.BaseUrl + "/metrics")		
		if err != nil {
			// log or print something here.
			continue
		}
		defer resp.Body.Close()

		var parser expfmt.TextParser
		metrics, err := parser.TextToMetricFamilies(resp.Body)
		if err != nil {
			// log or print something
			continue
		}

		var cpuValueToPrint, memValueToPrint float64

		if cpuValue, ok := metrics["node_cpu_stats"]; ok {
			cpuValueToPrint = float64(cpuValue.GetMetric()[0].GetGauge().GetValue())
			cpuGauge.WithLabelValues("something").Set(float64(cpuValue.GetMetric()[0].GetGauge().GetValue()))
		}
	
		if memValue, ok := metrics["node_mem_stats"]; ok {
			memValueToPrint = float64(memValue.GetMetric()[0].GetGauge().GetValue())
			memGauge.WithLabelValues("endpoint").Set(float64(memValue.GetMetric()[0].GetGauge().GetValue()))
		}

		fmt.Printf("Got here. CPU = %.2f, Memory = %.2f\n", cpuValueToPrint, memValueToPrint)
	}
}



