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

// TODO: keep adding things here for the constructor pattern 
// TODO: fix enums (both individual enums and their declarations)
func (lb LoadBalancer) WithIpAddr(ip string) LoadBalancer {
	lb.IpAddr = ip
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

// The actual algorithm behind the load balancing
//
// Uses a Heap data structure to store the servers, updating their priorities
// according to CPU and Memory usage, and (maybe) default and hand-picked weights
// depending on their performance. 
func (lb *LoadBalancer) SelectServer() server.RemoteServer {
	// TODO
	return server.RemoteServer{}
}

func (lb *LoadBalancer) GetServerStatus(server *server.RemoteServer) {
	cmd := exec.Command("top", "-b", "-n", "1")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal("Could not get the system statistics:", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()	
		
		if strings.HasPrefix(line, "%Cpu(s): ") {
			fmt.Println(line)
			num := strings.Split(line, "ni, ")[1][:4]
			idleCpu, _ := strconv.ParseFloat(num, 32)
			server.CPU_Usage = float32(100) - float32(idleCpu)
			fmt.Println(server.CPU_Usage)
		} else if strings.HasPrefix(line, "MiB Mem :") {
			fmt.Println(line)

			split := strings.Split(line, ":  ")
			used := split[1][:7] 
			idle := split[1][16:23]
			
			idleMem, _ := strconv.ParseFloat(idle, 32)
			usedMem, _ := strconv.ParseFloat(used, 32)

			server.MEM_Usage = float32(usedMem) - float32(idleMem)
			fmt.Println(server.MEM_Usage)
		}
	}
}


