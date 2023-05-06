package loadbalancer

import (
	lbheap "github.com/scarpart/distributed-task-scheduler/load-balancer/lb-heap"
	"github.com/scarpart/distributed-task-scheduler/load-balancer/server"
)

type ServerMonitor struct {
	Heap lbheap.Heap
}

func NewServerMonitor() ServerMonitor {
	return ServerMonitor{
		Heap: lbheap.NewHeap(),
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


