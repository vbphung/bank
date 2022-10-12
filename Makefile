MIGRATE_PATH=db/migrate
DB_URL=postgresql://root:postgres_password@localhost:5432/simplebank?sslmode=disable

init-postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=postgres_password -d postgres:14-alpine

create-db:
	docker exec -it postgres createdb --username=root --owner=root simplebank

drop-db:
	docker exec -it postgres dropdb simplebank

migrate-up:
	migrate -path $(MIGRATE_PATH) -database $(DB_URL) -verbose up

migrate-down:
	migrate -path $(MIGRATE_PATH) -database $(DB_URL) -verbose down

.PHONY: init-postgres create-db drop-db migrate-up migrate-down
