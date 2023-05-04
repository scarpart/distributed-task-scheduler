package server

type RemoteServer struct {
	IpAddr    string  `json:"ip_addr"`
	CPU_Usage float32 `json:"cpu_usage"`
	MEM_Usage float32 `json:"mem_usage"`
	Weight 	  int32   `json:"weight"`
	BaseUrl   string  `json:"base_url"`
	ApiKey    string  `json:"api_key"`
	// This is probably not needed. Need to rethink it, might be useless or there is a better way to do it
	TotalValue     string  `json:"total_value"`
}

func (server *RemoteServer) MoreUsageThan(other *RemoteServer) bool {
	left := (server.CPU_Usage + server.MEM_Usage) * float32(server.Weight)
	right := (other.CPU_Usage + other.MEM_Usage) * float32(other.Weight)
	return left > right
}

func (server *RemoteServer) Value() float32 {
	return float32(server.Weight) * (server.CPU_Usage + server.MEM_Usage)
}



