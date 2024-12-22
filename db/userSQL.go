package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"nest/models"
)

func GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	query := "SELECT id, username, email FROM users WHERE id = $1"
	err := Pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	query := `
        SELECT id, username, email, first_name, last_name, password_hash, role, created_at
        FROM users 
        WHERE username = $1
    `
	err := Pool.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("query error: %w", err)
	}
	return &user, nil
}

func CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := `
		INSERT INTO users (first_name, last_name, username, email, password_hash)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	err := Pool.QueryRow(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.Username,
		user.Email,
		user.PasswordHash,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func DeleteUser(ctx context.Context, userID int) error {
	query := `
		DELETE FROM users
		WHERE user_id = $1;
	`
	_, err := Pool.Exec(
		ctx,
		query,
		userID,
	)

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
