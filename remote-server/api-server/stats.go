package apiserver

type ServerStats struct {	
	CPU_Usage  float32 `json:"cpu_usage"`
	MEM_Usage  float32 `json:"mem_usage"`
	BaseUrl    string  `json:"base_url"`
}
