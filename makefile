postgres:
	docker run --name postgresDTS -p 5432:5432 -e POSTGRES_DB=distributed_task_scheduler_db -e POSTGRES_USER=taskmanager -e POSTGRES_PASSWORD=distributed-tasks -d postgres:latest

createdb:
	docker exec -it postgresDTS createdb --username=taskmanager --owner=taskmanager distributed_task_scheduler_db

dropdb:
	docker exec -it postgresDTS dropdb --username=taskmanager distributed_task_scheduler_db

migrateup:
	migrate -path db/migration -database "postgresql://taskmanager:distributed-tasks@localhost:5432/distributed_task_scheduler_db?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://taskmanager:distributed-tasks@localhost:5432/distributed_task_scheduler_db?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

main:
	go run main-server/main.go 

server:
	go run remote-server/main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test main server 
