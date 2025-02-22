package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"nest/models"

	"github.com/jackc/pgx/v4"
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

func IsUsernameTaken(ctx context.Context, username string) (bool, error) {
	query := `
		SELECT 1
		FROM users
		WHERE username = $1
	`
	var result int

	err := Pool.QueryRow(ctx, query, username).Scan(&result)

	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("query error: %w", err)
	}
	return true, nil
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

func UpdateUserPassword(ctx context.Context, userID int, hashedPassword string) error {
	query := `
		UPDATE users
		SET password = $1
		WHERE user_id = $2
	`
	_, err := Pool.Exec(ctx, query, hashedPassword, userID)
	if err != nil {
		return fmt.Errorf("failed to update user password: %w", err)
	}

	return nil
}

func UpdateUserFirstName(ctx context.Context, userID int, firstName string) error {
	query := `
		UPDATE users
		SET first_name = $1
		WHERE user_id = $2
	`
	_, err := Pool.Exec(ctx, query, firstName, userID)
	if err != nil {
		return fmt.Errorf("failed to update user first name: %w", err)
	}

	return nil
}

func UpdateUserLastName(ctx context.Context, userID int, lastName string) error {
	query := `
		UPDATE users
		SET last_name = $1
		WHERE user_id = $2
	`
	_, err := Pool.Exec(ctx, query, lastName, userID)
	if err != nil {
		return fmt.Errorf("failed to update user last name: %w", err)
	}

	return nil
}

func UpdateUserEmail(ctx context.Context, userID int, email string) error {
	query := `
		UPDATE users
		SET email = $1
		WHERE user_id = $2
	`
	_, err := Pool.Exec(ctx, query, email, userID)
	if err != nil {
		return fmt.Errorf("failed to update user email: %w", err)
	}

	return nil
}
