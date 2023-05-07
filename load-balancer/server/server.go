package server

type RemoteServer struct {
	IpAddr    string  `json:"ip_addr"`
	BaseUrl   string  `json:"base_url"`
	ApiKey    string  `json:"api_key"`
	// This is probably not needed. Need to rethink it, might be useless or there is a better way to do it
}

// temporary
type ServerStats struct {
	MEM_Usage float32 
	CPU_Usage float32
	Weight    int32
}	

func (ss *ServerStats) Value() float32 {
	return float32(1)
}



