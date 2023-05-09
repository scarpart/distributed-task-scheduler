package main

import (
	"log"
	"net"

	_ "github.com/lib/pq"
	loadbalancer "github.com/scarpart/distributed-task-scheduler/main-server/load-balancer"
	"github.com/scarpart/distributed-task-scheduler/util"
)

var lb *loadbalancer.LoadBalancer 

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not read .evn config file:", err)
	}
	
	lb := loadbalancer.NewLoadBalancer().
		WithIpAddr(net.IP(config.SERVER_ADDRESS)).
		WithPort(config.PORT)

	server1 := &loadbalancer.RemoteServer{
		IpAddr: "127.0.0.1",
		Port: "8080", 
	}
	server2 := &loadbalancer.RemoteServer{
		IpAddr: "127.0.0.1",
		Port: "8081", 
	}
	
	lb.Servers.Add(server1)
	lb.Servers.Add(server2)

	err = lb.Start()
	if err != nil {
		log.Fatal("Could not start the Load Balancer: ", err)
	}
}

