package db

import (
	"context"
	"errors"
	"nest/models"
)

func GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	query := "SELECT id, name, email FROM users WHERE id = $1"
	err := Pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func CreateUser(username string, hashedPassword []byte) error {
	// var user models.User
	// query := "INSERT INTO id, name, email FROM users WHERE id = $1"
	// err := Pool.Inse(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email)
	// if err != nil {
	// 	return nil, errors.New("user not found")
	// }
	return nil
}

func GetUserByUsername(username string) (*models.User, error) {
	return nil, errors.New("")
}
