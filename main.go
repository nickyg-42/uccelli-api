package main

import (
	"fmt"
	"log"
	"nest/db"
	"nest/routes"
	"os"

	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dbConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db.InitDB(dbConnStr)

	router := routes.RegisterRoutes()

	log.Println("Server is running on port 5000...")
	log.Fatal(http.ListenAndServe("localhost:5000", router))
}
