package loadbalancer 

type RemoteServer struct {
	Port      	string  `json:"port"`
	IpAddr    	string  `json:"ip_addr"`
	ApiKey    	string  `json:"api_key"`
	Connections int32   `json:"connections"`
	// This is probably not needed. Need to rethink it, might be useless or there is a better way to do it
}

// temporary
type ServerStats struct {
	MEM_Usage float32 
	CPU_Usage float32
	Weight    int32
}	

func (rm *RemoteServer) URL() string {
	return rm.IpAddr + ":" + rm.Port
}

func (ss *ServerStats) Value() float32 {
	return float32(1)
}



