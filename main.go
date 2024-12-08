package main

import (
	"log"
	"nest/db"
	"nest/routes"

	"net/http"
)

func main() {
	db.InitDB("postgres://birdman:Blueberry42@192.168.0.207:5432/master")

	routes.RegisterRoutes()

	//Start the server
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
