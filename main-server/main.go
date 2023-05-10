package main

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
	loadbalancer "github.com/scarpart/distributed-task-scheduler/main-server/load-balancer"
	"github.com/scarpart/distributed-task-scheduler/util"
)

func main() {
	addrToKey := make(map[string]string)

	config, err := util.LoadConfig("./main-server/")
	if err != nil {
		log.Fatal("Could not read .env config file:", err)
	}
	
	lb := loadbalancer.NewLoadBalancer().
		WithServerAddr(config.SERVER_ADDRESS)

	addrToKey["http://127.0.0.1:8080"] = ""
 	lb.InitRemoteServers(addrToKey)

	fmt.Printf("here in main after add: %v\n", lb.Servers)

	err = lb.Start()
	if err != nil {
		log.Fatal("Could not start the Load Balancer: ", err)
	}
}

