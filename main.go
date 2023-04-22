package main

import (
	"database/sql"
	"log" 

	_ "github.com/lib/pq"
	"github.com/scarpart/distributed-task-scheduler/api"
	db "github.com/scarpart/distributed-task-scheduler/db/sqlc"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://taskmanager:distributed-tasks@localhost:5432/distributed_task_scheduler_db?sslmode=disable"
	serverAddress = "127.0.0.1:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to the database:", err)
	}

	store := db.NewStore(conn)
	server := db.NewServer(store)

	err := server.Start(serverAddress)
	if err != nil {
		log.Fatal("could not start server:", err)
	}
}
