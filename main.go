package main

import (
	"context"
	"fmt"
	"log"
	"nest/db"
	"nest/models"
	"nest/routes"
	"nest/utils"
	"os"
	"sync"
	"time"

	"net/http"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	logFile, err := os.OpenFile("/var/log/uccelli-api.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Warning: Failed to open log file, defaulting to stdout: %v", err)
	} else {
		defer logFile.Close()
		log.SetOutput(logFile)
		log.Println("Log file attached...")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dbConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db.InitDB(dbConnStr)

	router := routes.RegisterRoutes()

	location := time.Local
	eastern, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Println("Failed to load EST:", err)
	} else {
		location = eastern
	}
	s := gocron.NewScheduler(location)

	_, err = s.Every(1).Day().At("09:00").Do(func() {
		ctx := context.Background()
		events, err := db.GetEventsForTomorrow(ctx, time.Now().In(location))
		if err != nil {
			log.Printf("Error fetching events for tomorrow: %v", err)
			return
		}

		var wg sync.WaitGroup
		for _, event := range events {
			wg.Add(1)
			go func(e models.Event) {
				defer wg.Done()
				subject := fmt.Sprintf("Upcoming Event: %s", e.Name)
				body := fmt.Sprintf("**%s** is starting tomorrow at %s\n\nDescription: %s\n\nYou can view it here: %s",
					e.Name,
					e.StartTime.In(location).Format("3:04 PM"),
					e.Description,
					"https://uccelliapp.duckdns.org",
				)
				utils.NotifyAllUsersInGroup(int(e.GroupID), subject, body)
			}(event)
		}
		wg.Wait()
	})

	if err != nil {
		log.Printf("Error scheduling event notifications: %v", err)
	} else {
		log.Println("Event notifications scheduled successfully...")
	}

	s.StartAsync()

	log.Println("Server is running on port 5000...")
	log.Fatal(http.ListenAndServe("127.0.0.1:5000", router))
}
