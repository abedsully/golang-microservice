package main

import (
	"log"

	"github.com/abedsully/golang-microservice/api/handler"
	"github.com/abedsully/golang-microservice/grpc/pb"
	"github.com/ianschenck/envflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const minSecretKeySize = 32

func main() {
	var (
		secretKey = envflag.String("SECRET_KEY", "01234567890123456789012345678901", "secret key for jwt signing")
		svcAddr = envflag.String("GRPC_SVC_ADDR", "0.0.0.0:9091", "address where grpc service is listening on")
	)

	if len(*secretKey) < minSecretKeySize {
		log.Fatalf("SECRET_KEY must be at least %d characters, now: %d", minSecretKeySize, len(*secretKey))
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(*svcAddr, opts...)

	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}

	defer conn.Close()

	client := pb.NewGolangMicroserviceClient(conn)

	hdl := handler.NewHandler(client, *secretKey)
	handler.RegisterRoutes(hdl)
	handler.Start(":8080")
}