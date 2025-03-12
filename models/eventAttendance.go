package models

import "time"

type EventAttendance struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	EventID   int       `json:"event_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type AttendanceData struct {
	UserID  int    `json:"user_id"`
	EventID int    `json:"event_id"`
	Status  string `json:"status"`
}
