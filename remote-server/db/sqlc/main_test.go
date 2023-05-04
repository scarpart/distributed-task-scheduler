package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/scarpart/distributed-task-scheduler/util"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Could not read the config file:", err) 
	}

	conn, err := sql.Open(config.DB_DRIVER, config.DB_SOURCE)

	if err != nil {
		log.Fatal("cannot connect to the database:", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
