package db

import (
	"context"
	"fmt"
	"nest/models"
)

func FetchAllEventsForGroup(ctx context.Context, groupID int) ([]models.Event, error) {
	query := `
		SELECT id, group_id, created_by, name, description, start_time, end_time, created_at
		FROM events
		WHERE group_id = $1
	`

	rows, err := Pool.Query(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch events for group %d: %w", groupID, err)
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

func FetchAllEventsByUser(ctx context.Context, userID int) ([]models.Event, error) {
	query := `
		SELECT id, group_id, created_by, name, description, start_time, end_time, created_at
		FROM events
		WHERE created_by = $1
	`

	rows, err := Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch events for user %d: %w", userID, err)
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

func FetchAllEventsByGroup(ctx context.Context, groupID int) ([]models.Event, error) {
	query := `
		SELECT id, group_id, created_by, name, description, start_time, end_time, created_at
		FROM events
		WHERE group_id = $1
	`

	rows, err := Pool.Query(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch events for group %d: %w", groupID, err)
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
