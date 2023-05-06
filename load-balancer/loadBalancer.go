package loadbalancer

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/scarpart/distributed-task-scheduler/enums"
	lbheap "github.com/scarpart/distributed-task-scheduler/load-balancer/lb-heap"
	"github.com/scarpart/distributed-task-scheduler/load-balancer/server"
)

type LoadBalancer struct {
	IpAddr  string
	Port    string
	Servers lbheap.Heap 
	Status  enums.NodeStatus  
	BaseUrl string
}

func NewLoadBalancer() LoadBalancer {
	return LoadBalancer{
		IpAddr: "127.0.0.1",
		Port: "8080",
		Servers: lbheap.NewHeap(),
		Status: 1, // TODO: fix this 
		BaseUrl: "127.0.0.1:8080/",
	}
}

func (lb LoadBalancer) WithIpAddr(ip string) LoadBalancer {
	lb.IpAddr = ip
	return lb
}	

func (lb LoadBalancer) WithPort(port string) LoadBalancer {
	lb.Port = port
	return lb
}

func (lb LoadBalancer) WithBaseUrl(url string) LoadBalancer {
	lb.BaseUrl = url
	return lb
}

func (lb *LoadBalancer) DistributeRequest(ctx *gin.Context) {
	client := &http.Client{}

	server := lb.SelectServer()
	req, err := http.NewRequest(ctx.Request.Method, server.BaseUrl + ctx.Request.URL.Path, ctx.Request.Body)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Connects using the individual server's API key, which is read from a config file 
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", server.ApiKey))	
	
	resp, err := client.Do(req)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

// Probably going to be replaced by the heap monitor
// The actual algorithm behind the load balancing
//
// Uses a Heap data structure to store the servers, updating their priorities
// according to CPU and Memory usage, and (maybe) default and hand-picked weights
// depending on their performance. 
func (lb *LoadBalancer) SelectServer() server.RemoteServer {
	// TODO
	return server.RemoteServer{}
}


