package main

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	docker "github.com/fsouza/go-dockerclient"
)

type WorkerNode struct {
	name 			  string 
	nodeID            uint32
	addr              string
	taskTopic         string
	errorTopic        string 
	heartbeatTopic    string
	maxRetryCount     int32
	heartbeatInterval time.Duration 
	taskConsumer      *kafka.Consumer
	heartbeatProducer *kafka.Producer
	errorProducer     *kafka.Producer
}

func NewNode(name string) *WorkerNode {
	h := fnv.New32a()
	h.Write([]byte(name + time.Now().String()))

	return &WorkerNode{
		name: name,
		nodeID: h.Sum32(),
	}
}

func (node *WorkerNode) WithAddr(addr string) *WorkerNode {
	node.addr = addr
	return node
}

func (node *WorkerNode) WithTopics(taskTopic, errorTopic, heartbeatTopic string) *WorkerNode {
	node.taskTopic = taskTopic
	node.errorTopic = errorTopic
	node.heartbeatTopic = heartbeatTopic
	return node
}

func (node *WorkerNode) WithHeartbeatInterval(heartbeatInterval time.Duration) *WorkerNode {
	node.heartbeatInterval = heartbeatInterval
	return node
}

func (node *WorkerNode) WithMaxRetryCount(maxRetryCount int32) *WorkerNode {
	node.maxRetryCount = maxRetryCount
	return node
}

func (node *WorkerNode) Build() error {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "something-blank-for-now",
		"group.id":			 "some-group-id",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return err
	}

	consumer.SubscribeTopics([]string{
		node.heartbeatTopic, node.taskTopic, node.errorTopic,
	}, nil)
	node.taskConsumer = consumer

	hbProducer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "something-blank-for-now",
	}) 
	if err != nil {
		return err 
	}
	node.heartbeatProducer = hbProducer

	errProducer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "something-blank-for-now",
	})
	if err != nil {
		return err 
	}
	node.errorProducer = errProducer

	return nil
}

func (node *WorkerNode) Run(ctx context.Context) error {
	ticker := time.NewTicker(node.heartbeatInterval)	
	
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			node.taskConsumer.Close()
			node.errorProducer.Close()
			node.heartbeatProducer.Close()
			return ctx.Err()
		
		case <-ticker.C:
			node.heartbeatProducer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &node.heartbeatTopic, Partition: kafka.PartitionAny},
				Value:          []byte(fmt.Sprintf("Heartbeat from node %d at %s", node.nodeID, time.Now().String())),
			}, nil)
				
		default:
			msg, err := node.taskConsumer.ReadMessage(-1)
			if err != nil {
				node.errorProducer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &node.errorTopic, Partition: kafka.PartitionAny},
					Value		  : []byte(err.Error()),
				}, nil)
				continue
			}
			
			var task Task
			if err := json.Unmarshal(msg.Value, &task); err != nil {
				node.errorProducer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &node.errorTopic, Partition: kafka.PartitionAny},
					Value:          []byte(err.Error()),
				}, nil)
			}

			client, err := docker.NewClientFromEnv()
			if err != nil {
				node.errorProducer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &node.errorTopic, Partition: kafka.PartitionAny},
					Value:          []byte(err.Error()),
				}, nil)
			}
			
			container, err := client.CreateContainer(docker.CreateContainerOptions{
				Config: &docker.Config{
					Image: task.Image,
					Cmd:   []string{task.Command, task.Args}, 
				},
				HostConfig: &docker.HostConfig{
					NetworkMode: "none",	
				},
			})
			if err != nil {
				return err 
			}

			client.StartContainer(container.ID, nil)
			exitCode, err := client.WaitContainer(container.ID)
			if err != nil {
				node.errorProducer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &node.errorTopic, Partition: kafka.PartitionAny},
					Value:          []byte(err.Error()),
				}, nil)
			}

			logOptions := docker.LogsOptions{
				Container:    container.ID,
				OutputStream: os.Stdout,
				ErrorStream:  os.Stderr, 
				Stdout:       true,
				Stderr:       true,
			}

			err = client.Logs(logOptions)
			if err != nil {
				node.errorProducer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &node.errorTopic, Partition: kafka.PartitionAny},
					Value:          []byte(err.Error()),
				}, nil)
			}

			if exitCode != 0 {
				node.errorProducer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &node.errorTopic, Partition: kafka.PartitionAny},
					Value:          []byte(err.Error()),
				}, nil)
			}

			client.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID, Force: true})
		}
	}
}

