package loadbalancer

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type LoadBalancer struct {
	URL         string
	Servers  	*Heap 
	mutex		sync.Mutex
	client      http.Client
	Router      *gin.Engine
}

func NewLoadBalancer() *LoadBalancer {
	// The HTTP Transport ensures that the remote servers have a concurrent connection cap and do not get overwhelmed
	lb := &LoadBalancer{
		Servers: &Heap{},  
		mutex: sync.Mutex{}, 
		client: http.Client{
			Transport: &http.Transport{
				MaxConnsPerHost: 10,
			},
		},
	}

	router := gin.Default()

	// POST 
	router.POST("/task", lb.DistributeRequest)
	router.POST("/node", lb.DistributeRequest)
	router.POST("/mapping", lb.DistributeRequest)
	
	// GET 
	//router.GET("/task", lb.DistributeRequest)
	router.GET("/node/:node_id", lb.DistributeRequest)
	router.GET("/node", lb.DistributeRequest)

	lb.Router = router
	return lb
}

func (lb *LoadBalancer) InitRemoteServers(addrToKey map[string]string) {
	for addr := range addrToKey {
		server := &RemoteServer{
			URL: addr,
			Mutex: sync.Mutex{},
			//ApiKey: apiKey,
		}
		lb.Servers.Add(server)	
		fmt.Printf("from initRemoteServers: %v\n", server)
		go server.HealthCheck()	
	}
}

func (lb *LoadBalancer) DistributeRequest(ctx *gin.Context) {
	server := lb.Servers.LeastConnections()
	if server == nil {
		fmt.Println("server is nil")
		return 
	}

	fmt.Println("server is not nil")
	fmt.Printf("s.conns = %d\n", server.Connections)
	fmt.Printf("s.url = %s\n",server.URL)

	server.Mutex.Lock()
	isAvailable := server.IsAvailable 
	server.Mutex.Unlock()

	fmt.Println(isAvailable)

	if !isAvailable {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Creates the request to be executed by the HTTP Client
	req, err := http.NewRequest(ctx.Request.Method, server.URL + ctx.Request.URL.Path, ctx.Request.Body)
	fmt.Println("request", req)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	for k, v := range ctx.Request.Header {
		req.Header[k] = v
	}

	// Updates the connection counter for the server and fixes the heap property
	lb.UpdateConnections(server, 1)

	body, statusCode, contentType, err := lb.SendRequest(ctx, server, req)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)	
		lb.UpdateConnections(server, -1)
		return
	}	

	lb.UpdateConnections(server, -1)

	ctx.Data(statusCode, contentType, body)
}

func (lb *LoadBalancer) UpdateConnections(server *RemoteServer, amount int32) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	server.Connections += amount
	lb.Servers.Fix(server)
}

func (lb *LoadBalancer) SendRequest(ctx *gin.Context, server *RemoteServer, req *http.Request) ([]byte, int, string, error)  {
	// Connects using the individual server's API key, which is read from a config file 
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", server.ApiKey))	
	
	resp, err := lb.client.Do(req)
	if err != nil {
		return nil, 0, "", err
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, "", err
	}
	
	return body, resp.StatusCode, resp.Header.Get("Content-Type"), nil
}

func (lb *LoadBalancer) Start() error {
	fmt.Printf("on Start: %v\n", lb.Servers)
	return lb.Router.Run(lb.URL)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}

func (lb LoadBalancer) WithServerAddr(addr string) LoadBalancer {
	lb.URL = addr 
	return lb
}	

