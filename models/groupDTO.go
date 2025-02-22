package models

type GroupDTO struct {
	Name        string `json:"name"`
	CreatedByID int64  `json:"created_by_id"`
}
