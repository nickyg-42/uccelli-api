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

	// logFile, err := os.OpenFile("/var/log/uccelli-api.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	log.Fatalf("Failed to open log file: %v", err)
	// }
	// defer logFile.Close()

	// log.SetOutput(logFile)
	// log.Println("log file attached...")

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dbConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db.InitDB(dbConnStr)

	router := routes.RegisterRoutes()

	log.Println("Server is running on port 5000...")
	log.Fatal(http.ListenAndServe("127.0.0.1:5000", router))
}
