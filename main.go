package main

import (
	"database/sql"
	"log"

	"github.com/herbi-dino/bank/api"
	db "github.com/herbi-dino/bank/db/sqlc"
	"github.com/herbi-dino/bank/utils"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := utils.LoadConfig("./")
	if err != nil {
		log.Fatalf("load config failed: %+v\n", err)
	}

	conn, err := sql.Open(cfg.DbDriver, cfg.DbSource)
	if err != nil {
		log.Fatalf("connect db failed: %+v\n", err)
	}

	store := db.CreateStore(conn)

	server := api.CreateServer(store)

	err = server.Start(cfg.ServerAddr)
	if err != nil {
		log.Fatalf("start server failed: %+v\n", err)
	}
}
