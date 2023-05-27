package main

import (
	"log"
	"time"
)

func main() {
	node := NewNode("node1").
		WithAddr("localhost:9093").
		WithTopics("tasks", "task-errors", "heartbeats").
		WithMaxRetryCount(3).// TODO: make this configurable (env)
		WithHeartbeatInterval(time.Second * 10) // TODO: make this configurable also

	err := node.Build()
	if err != nil {
		log.Fatal(err)
	}

	node.Run()
}
