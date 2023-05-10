package loadbalancer 

//
//import (
//	"fmt"
//	"net/http"
//	"os"
//	"time"
//
//	"github.com/gin-gonic/gin"
//	"github.com/prometheus/client_golang/prometheus"
//	"github.com/prometheus/common/expfmt"
//)
//
//type ServerMonitor struct {
//	Heap     Heap
//	Servers  []RemoteServer
//	CpuGauge *prometheus.GaugeVec
//	MemGauge *prometheus.GaugeVec
//}
//
//func NewServerMonitor() ServerMonitor {
//	return ServerMonitor{
//		Heap: NewHeap(),
//		Servers: []RemoteServer{},
//		CpuGauge: prometheus.NewGaugeVec(
//			prometheus.GaugeOpts{
//				Name: "cpu_usage",
//				Help: "CPU usage of a remote server in percentage",
//			},
//			[]string{"endpoint"},
//		),
//		MemGauge: prometheus.NewGaugeVec(
//			prometheus.GaugeOpts{
//				Name: "mem_usage",
//				Help: "Memory usage of a remote server in percentage",
//			},
//			[]string{"endpoint"},
//		),
//	}
//}
//
//// TODO: Make it so ServerMonitor listens to the /metrics endpoint and updates
//// the heap accordingly, in an interval of 5 seconds. There should be a method
//// to return the root node to the main thread of the LoadBalancer, which will 
//// then distribute the request to that remote server.
//
//// I think Poll removes the node. Make a method to just get it. 
//func (sm *ServerMonitor) GetLeastBusy() (*server.ServerStats, error) {
//	return sm.Heap.Poll()
//}
//
//func (sm *ServerMonitor) WatchServers(ctx *gin.Context) {
//	timer := time.NewTicker(time.Second * 15)
//	client := &http.Client{
//		Timeout: time.Second * 10,
//	}	
//
//	for {
//		select {
//		case <-timer.C:
//			sm.GetMetrics(ctx, client)					
//			os.Exit(0)
//		}
//	}
//}
//
//func (sm *ServerMonitor) GetMetrics(ctx *gin.Context, client *http.Client) {
//	for _, server := range sm.Servers {
//		resp, err := client.Get(server.BaseUrl + "/metrics")		
//		if err != nil {
//			// log or print something here.
//			continue
//		}
//		defer resp.Body.Close()
//
//		var parser expfmt.TextParser
//		metrics, err := parser.TextToMetricFamilies(resp.Body)
//		if err != nil {
//			// log or print something
//			continue
//		}
//
//		var cpuValueToPrint, memValueToPrint float64
//
//		if cpuValue, ok := metrics["node_cpu_stats"]; ok {
//			cpuValueToPrint = float64(cpuValue.GetMetric()[0].GetGauge().GetValue())
//			sm.CpuGauge.WithLabelValues("something").Set(float64(cpuValue.GetMetric()[0].GetGauge().GetValue()))
//		}
//	
//		if memValue, ok := metrics["node_mem_stats"]; ok {
//			memValueToPrint = float64(memValue.GetMetric()[0].GetGauge().GetValue())
//			sm.MemGauge.WithLabelValues("endpoint").Set(float64(memValue.GetMetric()[0].GetGauge().GetValue()))
//		}
//
//		fmt.Printf("Got here. CPU = %.2f, Memory = %.2f\n", cpuValueToPrint, memValueToPrint)
//	}
//}
//
//
//
//
