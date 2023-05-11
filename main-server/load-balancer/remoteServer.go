package loadbalancer

import (
	"sync"
)

type RemoteServer struct {
	URL 		string 
	ApiKey      string `json:"api_key"`
	Connections int32  `json:"connections"`
	IsAvailable bool  
	Mutex       sync.Mutex 
	// This is probably not needed. Need to rethink it, might be useless or there is a better way to do it
}

