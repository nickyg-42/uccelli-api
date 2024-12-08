package db

import (
	"context"
	"errors"
)

type User struct {
	ID    int
	Name  string
	Email string
}

func GetUserByID(ctx context.Context, id int) (*User, error) {
	var user User
	query := "SELECT id, name, email FROM users WHERE id = $1"
	err := Pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, errors.New("User not found")
	}
	return &user, nil
}
