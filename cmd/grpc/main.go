package main

import (
	"log"
	"net"

	"github.com/abedsully/golang-microservice/db"
	"github.com/abedsully/golang-microservice/grpc/pb"
	"github.com/abedsully/golang-microservice/grpc/server"
	"github.com/abedsully/golang-microservice/grpc/storer"
	"github.com/ianschenck/envflag"
	"google.golang.org/grpc"
)

func main() {

	var (
		svcAddr = envflag.String("SVC_ADDR", "0.0.0.0:9091", "address where grpc service is listening on")
	)

	// instantiate db
	db, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	log.Println("Successfully connected to database")

	// instantiate server
	st := storer.NewMySqlStorer(db.GetDB())
	srv := server.NewServer(st)

	// register server with gRPC server
	grpcSrv := grpc.NewServer()
	pb.RegisterGolangMicroserviceServer(grpcSrv, srv)

	listener, err := net.Listen("tcp", *svcAddr)

	if err != nil {
		log.Fatalf("listener failed: %v", err)
	}

	log.Printf("server is listening on %s", *svcAddr)
	err = grpcSrv.Serve(listener)

	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
