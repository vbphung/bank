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

.PHONY: init-postgres create-db drop-db migrate-up migrate-down sqlc test
