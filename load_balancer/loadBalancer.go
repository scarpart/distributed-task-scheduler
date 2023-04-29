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
	"github.com/scarpart/distributed-task-scheduler/api"
	db "github.com/scarpart/distributed-task-scheduler/db/sqlc"
	"github.com/scarpart/distributed-task-scheduler/enums"
)

// This should be `LocalServer` and should not be declared here.
//type Server struct {
//	CPU_Usage float32
//	MEM_Usage float32
//	Priority int32
//}

type LoadBalancer struct {
	Port string
	Servers []Server
	Status enums.NodeStatus  
}

func (lb *LoadBalancer) DistributeRequest(ctx *gin.Context, taskRequest *db.CreateTaskParams) {
	client := &http.Client{}

	server := lb.RoundRobin()
	req, err := http.NewRequest("POST", server.Url, taskRequest)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	req.Header.Add("Authorization", "Bearer <token>")	
	
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
func (lb *LoadBalancer) RoundRobin() Server {
	// Implement Round Robin here. It should read the list of servers and pick the one that 
	// is most free in CPU and Memory usage, assigning to it a quantum, such that we have
	// efficient processing capabilities and such. Gotta figure out how to implement this. 
	return Server{}
}

func (server *Server) GetServerStatus() {
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


