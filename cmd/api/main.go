package main

import (
	"log"

	"github.com/abedsully/golang-microservice/api/handler"
	"github.com/abedsully/golang-microservice/api/server"
	"github.com/abedsully/golang-microservice/api/storer"
	"github.com/abedsully/golang-microservice/db"
	"github.com/ianschenck/envflag"
)

const minSecretKeySize = 32

func main() {
	var secretKey = envflag.String("SECRET_KEY", "01234567890123456789012345678901", "secret key for jwt signing")
	db, err := db.NewDatabase()

	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	defer db.Close()

	log.Println("Successfully connected to database")

	st := storer.NewMySqlStorer(db.GetDB())

	srv := server.NewServer(st)
	hdl := handler.NewHandler(srv, *secretKey)
	handler.RegisterRoutes(hdl)
	handler.Start(":8080")
}