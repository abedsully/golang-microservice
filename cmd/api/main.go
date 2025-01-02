package main

import (
	"log"

	"github.com/abedsully/golang-microservice/db"
)

func main() {
	db, err := db.NewDatabase()

	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	defer db.Close()

	log.Println("Successfully connected to database")
}