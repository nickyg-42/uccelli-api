package models

import "time"

type Event struct {
	ID          int64     `json:"id"`
	GroupID     int64     `json:"group_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	CreatedBy   User      `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}
