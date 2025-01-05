package main

import (
	"log"

	"github.com/abedsully/golang-microservice/api/handler"
	"github.com/abedsully/golang-microservice/api/server"
	"github.com/abedsully/golang-microservice/api/storer"
	"github.com/abedsully/golang-microservice/db"
)

func main() {
	db, err := db.NewDatabase()

	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	defer db.Close()

	log.Println("Successfully connected to database")

	st := storer.NewMySqlStorer(db.GetDB())

	srv := server.NewServer(st)
	hdl := handler.NewHandler(srv)
	handler.RegisterRoutes(hdl)
	handler.Start(":8080")
}