package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/vbph/bank/api"
	db "github.com/vbph/bank/db/sqlc"
	"github.com/vbph/bank/gapi"
	"github.com/vbph/bank/pb"
	"github.com/vbph/bank/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

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

	if cfg.ServerType == "http" {
		httpServer(cfg, store)
		return
	}

	grpcServer(cfg, store)
}

func httpServer(cfg utils.Config, store *db.Store) {
	server, err := api.CreateServer(cfg, store)
	if err != nil {
		log.Fatalf("create server failed: %+v\n", err)
	}

	err = server.Start(cfg.HTTPServerAddr)
	if err != nil {
		log.Fatalf("start server failed: %+v\n", err)
	}
}

func grpcServer(cfg utils.Config, store *db.Store) {
	server, err := gapi.CreateServer(cfg, store)
	if err != nil {
		log.Fatalf("create server failed: %+v\n", err)
	}

	gServer := grpc.NewServer()
	pb.RegisterBankServer(gServer, server)

	reflection.Register(gServer)

	listener, err := net.Listen("tcp", cfg.GrpcServerAddr)
	if err != nil {
		log.Fatalf("create listener failed: %+v\n", err)
	}

	err = gServer.Serve(listener)
	if err != nil {
		log.Fatalf("start server failed: %+v\n", err)
	}
}
