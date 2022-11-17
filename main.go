package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/DarrelASandbox/go-simple-bank/api_gin"
	"github.com/DarrelASandbox/go-simple-bank/db/util"
	"github.com/DarrelASandbox/go-simple-bank/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/DarrelASandbox/go-simple-bank/api"
	db "github.com/DarrelASandbox/go-simple-bank/db/sqlc"
	_ "github.com/lib/pq"
)

func main() {
	// "." refers to current folder
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	// runGinServer(config, store);
	runGrpcServer(config, store)
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := api_gin.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server")
	}
}
