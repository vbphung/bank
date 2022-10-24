package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/vbph/bank/api"
	db "github.com/vbph/bank/db/sqlc"
	"github.com/vbph/bank/gapi"
	"github.com/vbph/bank/pb"
	"github.com/vbph/bank/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

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

	switch cfg.ServerType {
	case "http":
		httpServer(cfg, store)
	case "grpc":
		grpcServer(cfg, store)
	default:
		go gatewayServer(cfg, store)
		grpcServer(cfg, store)
	}
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

func gatewayServer(cfg utils.Config, store *db.Store) {
	server, err := gapi.CreateServer(cfg, store)
	if err != nil {
		log.Fatalf("create server failed: %+v\n", err)
	}

	jsonOpts := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	gMux := runtime.NewServeMux(jsonOpts)

	ctx, ccl := context.WithCancel(context.Background())
	defer ccl()

	err = pb.RegisterBankHandlerServer(ctx, gMux, server)
	if err != nil {
		log.Fatalf("register server handler failed: %+v\n", err)
	}

	hMux := http.NewServeMux()
	hMux.Handle("/", gMux)

	listener, err := net.Listen("tcp", cfg.HTTPServerAddr)
	if err != nil {
		log.Fatalf("create listener failed: %+v\n", err)
	}

	err = http.Serve(listener, hMux)
	if err != nil {
		log.Fatalf("start server failed: %+v\n", err)
	}
}
