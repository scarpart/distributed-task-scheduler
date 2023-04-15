package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://taskmanager:distributed-tasks@localhost:5432/distributed_task_scheduler_db?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("cannot connect to the database:", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
