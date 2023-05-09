package loadbalancer

import (
	"net/http"
	"sync"
	"time"
)

type RemoteServer struct {
	URL 		string 
	ApiKey      string `json:"api_key"`
	Connections int32  `json:"connections"`
	IsAvailable bool  
	Mutex       sync.Mutex 
	// This is probably not needed. Need to rethink it, might be useless or there is a better way to do it
}

func (rs *RemoteServer) HealthCheck() {
	for {
		resp, err := http.Get(rs.URL + "/health")
		rs.Mutex.Lock()
		rs.IsAvailable = err == nil && resp.StatusCode == http.StatusOK
		rs.Mutex.Unlock()
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(5 * time.Second)
	}
}

