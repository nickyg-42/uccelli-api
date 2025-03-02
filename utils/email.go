package utils

import (
	"context"
	"log"
	"nest/db"
	"net/smtp"
	"os"
	"sync"
)

func SendEmail(to string, subject string, body string) error {
	from := os.Getenv("SMTP_ADDRESS")
	password := os.Getenv("SMTP_PASSWORD")

	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")

	msg := []byte("From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body)

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		from,
		[]string{to},
		msg,
	)
	return err
}

func NotifyAllUsersInGroup(groupID int, subject string, body string) {
	users, err := db.GetAllMembersForGroup(context.Background(), groupID)
	if err != nil {
		log.Printf("ERROR: Failed to get users for group %d: %v", groupID, err)
		return
	}

	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		go func(email string) {
			defer wg.Done()
			NotifyUser(email, subject, body)
		}(user.Email)
	}
	wg.Wait()
}

func NotifyUser(email string, subject string, body string) {
	err := SendEmail(email, subject, body)
	if err != nil {
		log.Printf("ERROR: Failed to send email to user %s: %v", email, err)
	} else {
		log.Printf("INFO: Sent email to user %s", email)
	}
}
