package main

import (
	"database/sql"
	"log"

	"github.com/herbi-dino/bank/api"
	db "github.com/herbi-dino/bank/db/sqlc"

	_ "github.com/lib/pq"
)

const (
	dbDriver   = "postgres"
	dbSource   = "postgresql://root:postgres_password@localhost:5432/bank?sslmode=disable"
	serverAddr = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("connect db failed: %+v\n", err)
	}

	store := db.CreateStore(conn)

	server := api.CreateServer(store)

	err = server.Start(serverAddr)
	if err != nil {
		log.Fatalf("start server failed: %+v\n", err)
	}
}
