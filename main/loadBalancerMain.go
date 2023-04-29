package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/scarpart/distributed-task-scheduler/api"
	db "github.com/scarpart/distributed-task-scheduler/db/sqlc"
	loadbalancer "github.com/scarpart/distributed-task-scheduler/load_balancer"
	"github.com/scarpart/distributed-task-scheduler/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not read .env config file:", err)
	}

	conn, err := sql.Open(config.DB_DRIVER, config.DB_SOURCE)
	if err != nil {
		log.Fatal("cannot connect to the database:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	test := loadbalancer.Server{}
	test.GetServerStatus()

	err = server.Start(config.SERVER_ADDRESS)
	if err != nil {
		log.Fatal("could not start server:", err)
	}


}
