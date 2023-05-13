package loadbalancer

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/scarpart/distributed-task-scheduler/util"
)

type LoadBalancer struct {
	URL         string
	Servers  	*Heap 
	mutex		sync.Mutex
	client      http.Client
	Router      *gin.Engine
	config 		util.Config
}

func NewLoadBalancer() *LoadBalancer {
	config, err := util.LoadConfig("./main-server")
	if err != nil {
		log.Println("Could not read LB Client configs. Using default values.", err)
	}
	// The HTTP Transport ensures that the remote servers have a concurrent connection cap and do not get overwhelmed
	lb := &LoadBalancer{
		Servers: &Heap{},  
		config: config,
		mutex: sync.Mutex{}, 
		client: http.Client{
			Timeout: time.Duration(config.LB_CONN_TIMEOUT) * time.Second,
			Transport: &http.Transport{
				MaxConnsPerHost: int(config.LB_CLIENT_MAX_CONNS),
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

func (lb *LoadBalancer) ServerHealthChecks(interval time.Duration) {
	for {
		lb.mutex.Lock()	
		for _, server := range *lb.Servers {
			url := server.URL + "/health"
			
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Println("Could not perform health check for server", server.URL, err)
				continue
			}

			resp, err := lb.client.Do(req)
			if err != nil {
				log.Println("Error in the response of server", server.URL, err)
				continue
			}

			if resp.StatusCode != http.StatusOK {
				log.Println("Server unavailable", server.URL, err)
				server.Mutex.Lock()
				server.IsAvailable = false
				server.Mutex.Unlock()
			} else {
				server.IsAvailable = true
			}

		}
		lb.mutex.Unlock()
		time.Sleep(interval)
	}
}

func (lb *LoadBalancer) InitRemoteServers(addrToKey map[string]string) {
	for addr := range addrToKey {
		server := &RemoteServer{
			URL: addr,
			Mutex: sync.Mutex{},
			//ApiKey: apiKey,
		}
		lb.Servers.Add(server)	
	}
	go lb.ServerHealthChecks(time.Second * time.Duration(lb.config.LB_HEALTH_CHECK_INTERVAL))
}

func (lb *LoadBalancer) DistributeRequest(ctx *gin.Context) {
	server := lb.Servers.LeastConnections()
	if server == nil {
		return 
	}

	server.Mutex.Lock()
	isAvailable := server.IsAvailable 
	server.Mutex.Unlock()

	fmt.Println(isAvailable)

	if !isAvailable {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Creates the request to be executed by the HTTP Client
	req, err := http.NewRequest(ctx.Request.Method, server.URL + ctx.Request.URL.Path + "?" + ctx.Request.URL.RawQuery, ctx.Request.Body)
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

func (lb *LoadBalancer) Start(certFile, keyFile string) error {
	// Configuring TLS such that we have secure HTTPS connections
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
		},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		},
		PreferServerCipherSuites: true,
	}

	httpsServer := &http.Server{
		Addr: lb.URL,
		Handler: lb.Router,
		TLSConfig: tlsConfig,
	}

	return httpsServer.ListenAndServeTLS(certFile, keyFile)
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

