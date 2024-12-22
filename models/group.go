package models

import "time"

type Group struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Owner     User      `json:"owner"`
	CreatedBy User      `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}
