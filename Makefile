MIGRATE_PATH=db/migrate
DB_URL=postgresql://root:postgres_password@localhost:5432/bank?sslmode=disable

init-postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=postgres_password -d postgres:14-alpine

create-db:
	docker exec -it postgres createdb --username=root --owner=root bank

drop-db:
	docker exec -it postgres dropdb bank

migrate-up:
	migrate -path $(MIGRATE_PATH) -database $(DB_URL) -verbose up

migrate-down:
	migrate -path $(MIGRATE_PATH) -database $(DB_URL) -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

proto:
	rm -rf pb
	mkdir pb
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto

.PHONY: init-postgres create-db drop-db migrate-up migrate-down sqlc test server proto
