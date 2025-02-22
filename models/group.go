package models

import "time"

type Group struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	CreatedByID int64     `json:"created_by_id"`
	CreatedAt   time.Time `json:"created_at"`
	Code        string    `json:"code"`
}
