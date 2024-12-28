package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"nest/models"
)

func CreateEvent(ctx context.Context, event *models.Event) (*models.Event, error) {
	query := `
		INSERT INTO events (group_id, created_by, name, description, start_time, end_time)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	err := Pool.QueryRow(
		ctx,
		query,
		event.GroupID,
		event.CreatedBy,
		event.Name,
		event.Description,
		event.StartTime,
		event.EndTime,
	).Scan(&event.ID, &event.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return event, nil
}

func GetEventByID(ctx context.Context, eventID int) (*models.Event, error) {
	var event models.Event
	query := `
        SELECT id, group_id, created_by, name, description, start_time, end_time, created_at
        FROM events 
        WHERE id = $1
    `
	err := Pool.QueryRow(ctx, query, eventID).Scan(
		&event.ID,
		&event.GroupID,
		&event.CreatedBy,
		&event.Name,
		&event.Description,
		&event.StartTime,
		&event.EndTime,
		&event.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("event not found")
		}
		return nil, fmt.Errorf("query error: %w", err)
	}
	return &event, nil
}

func GetAllEventsByUser(ctx context.Context, userID int) ([]models.Event, error) {
	query := `
		SELECT id, group_id, created_by, name, description, start_time, end_time, created_at
		FROM events
		WHERE created_by = $1
	`

	rows, err := Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get events for user %d: %w", userID, err)
	}
	defer rows.Close()

	var events []models.Event

	for rows.Next() {
		var event models.Event
		err = rows.Scan(
			&event.ID,
			&event.GroupID,
			&event.Name,
			&event.Description,
			&event.StartTime,
			&event.EndTime,
			&event.CreatedBy,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event row: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating event rows: %w", err)
	}

	return events, nil
}

func GetAllEventsByGroup(ctx context.Context, groupID int) ([]models.Event, error) {
	query := `
		SELECT id, group_id, created_by, name, description, start_time, end_time, created_at
		FROM events
		WHERE group_id = $1
	`

	rows, err := Pool.Query(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get events for group %d: %w", groupID, err)
	}
	defer rows.Close()

	var events []models.Event

	for rows.Next() {
		var event models.Event
		err = rows.Scan(
			&event.ID,
			&event.GroupID,
			&event.Name,
			&event.Description,
			&event.StartTime,
			&event.EndTime,
			&event.CreatedBy,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event row: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating event rows: %w", err)
	}

	return events, nil
}
