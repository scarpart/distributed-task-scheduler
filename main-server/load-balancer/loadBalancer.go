package loadbalancer

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
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
	// -------- ATTENTION: This is just for testing purposes, it should never be used in production.
	// What I'm doing is essentially just trusting my self-signed TLS certificates that are in the remote servers,
	// Since I don't really want to pay for domains and I'm just testing things in localhost. 
	cert, err := ioutil.ReadFile("remote-server/certs/cert.pem")
	if err != nil {
		log.Fatalf("Couldn't load server certifitate: %v", err)
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		log.Fatalf("Failed to append certificate")
	}
	// -------- ATTENTION: read the above 

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
				// ATTENTION: view code above 
				TLSClientConfig: &tls.Config{
					RootCAs: certPool,
				},
			},
		},
	}

	router := gin.Default()

	// This is just here temporarily, to facilitate development 
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true 
	router.Use(cors.New(corsConfig))

	// POST 
	router.POST("/task", lb.DistributeRequest)
	router.POST("/node", lb.DistributeRequest)
	router.POST("/mapping", lb.DistributeRequest)
	
	// GET 
	//router.GET("/task", lb.DistributeRequest)
	router.GET("/node/:node_id", lb.DistributeRequest)
	router.GET("/node", lb.DistributeRequest)

	// Sign in endpoints
	router.POST("/user", lb.DistributeRequest)

	// TESTING -- SHOULD BE REMOVED SOON
	// This is here so I can code the frontend up before finishing the backend.
	// The reason for that is simple: for the backend, I have to think and that takes time,
	// But for the fronted, I can just code whenever, so in idle times I can pick this up and do some work.
	// ie, it's for the sake of efficiency, since I know what the backend is going to look like.
	type Task struct {
		Status           string    `json:"status"`
		CpuUsage         float64   `json:"cpuUsage"`
		MemUsage         float64   `json:"memUsage"`
		TaskName         string    `json:"taskName"`
		TaskDescription  string    `json:"taskDescription"`
		TaskID           int       `json:"taskID"`
		UserID           int       `json:"userID"`
		ScheduledTime    time.Time `json:"scheduledTime"`
		NThreads         int       `json:"nThreads"`
		Priority         int       `json:"priority"`
		RetryCount       int       `json:"retryCount"`
		MaxRetries       int       `json:"maxRetries"`
		Dependencies     []int     `json:"dependencies"`
		CreatedTime      time.Time `json:"createdTime"`
		LastUpdatedTime  time.Time `json:"lastUpdatedTime"`
		ErrorMessage     string    `json:"errorMessage"`
		MachineID        string    `json:"machineID"`
		Output           string    `json:"output"`
	}

	statuses := []string{"ready", "running", "done"}
	var tasks []Task
	for i := 0; i < 5; i++ {
		tasks = append(tasks, Task{
			Status:           statuses[rand.Intn(len(statuses))],
			CpuUsage:         rand.Float64() * 100,
			MemUsage:         rand.Float64() * 100,
			TaskName:         "Task " + string(i),
			TaskDescription:  "This is task " + string(i),
			TaskID:           i,
			UserID:           rand.Intn(100) + 1,
			ScheduledTime:    time.Now().Add(time.Duration(rand.Intn(24))*time.Hour),
			NThreads:         rand.Intn(10) + 1,
			Priority:         rand.Intn(10) + 1,
			RetryCount:       0,
			MaxRetries:       3,
			Dependencies:     []int{},
			CreatedTime:      time.Now(),
			LastUpdatedTime:  time.Now(),
			ErrorMessage:     "",
			MachineID: 		  fmt.Sprintf("%d", rand.Intn(200)),			
			Output:           "output",
		})
	}

	router.GET("/scheduled-tasks", func(ctx *gin.Context) {
		ctx.JSON(200, tasks)
	})	

	lb.Router = router
	return lb
}

func (lb LoadBalancer) WithServerAddr(addr string) LoadBalancer {
	lb.URL = addr 
	return lb
}	

func (lb *LoadBalancer) ServerHealthChecks(interval time.Duration) {
	for {
		lb.mutex.Lock()	
		for _, server := range *lb.Servers {
			fmt.Println("in server health check for")

			url := server.URL + "/public/health"
			
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
				fmt.Println(server.URL, " should be available: ", server.IsAvailable)
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
	fmt.Println("Inside distribute request")

	server := lb.Servers.LeastConnections()
	if server == nil {
		return 
	}

	server.Mutex.Lock()
	isAvailable := server.IsAvailable 
	server.Mutex.Unlock()

	fmt.Println("server is available:", isAvailable)

	if !isAvailable {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Creates the request to be executed by the HTTP Client
	// TODO: change the logic for private servers here
	var group string
	switch ctx.Request.URL.Path {
	case "/user":
		group = "/public"
	default:
		group = "/private"
	}

	req, err := http.NewRequest(ctx.Request.Method, server.URL + group + ctx.Request.URL.Path + "?" + ctx.Request.URL.RawQuery, ctx.Request.Body)
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

