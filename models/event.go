package models

import "time"

type Event struct {
	ID          int64     `json:"id"`
	GroupID     int64     `json:"group_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	CreatedByID int64     `json:"created_by_id"`
	CreatedAt   time.Time `json:"created_at"`
}
