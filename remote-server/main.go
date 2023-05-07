package main 

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/lib/pq"
	api "github.com/scarpart/distributed-task-scheduler/remote-server/api-server"
	db "github.com/scarpart/distributed-task-scheduler/remote-server/db/sqlc"
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
	server := api.NewServer(store, net.IP(config.SERVER_ADDRESS))

	err = server.Start()
	if err != nil {
		log.Fatal("could not start server:", err)
	}
}
