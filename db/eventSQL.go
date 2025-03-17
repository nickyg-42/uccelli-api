package db

import (
	"context"
	"errors"
	"fmt"
	"nest/models"
	"time"

	"github.com/jackc/pgx/v4"
)

func CreateEvent(ctx context.Context, event *models.Event) (*models.Event, error) {
	query := `
		INSERT INTO events (group_id, created_by, name, description, start_time, end_time, location)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`

	err := Pool.QueryRow(
		ctx,
		query,
		event.GroupID,
		event.CreatedByID,
		event.Name,
		event.Description,
		event.StartTime,
		event.EndTime,
		event.Location,
	).Scan(&event.ID, &event.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return event, nil
}

func DeleteEvent(ctx context.Context, eventID int) error {
	query := `
		DELETE FROM events
		WHERE id = $1;
	`
	_, err := Pool.Exec(
		ctx,
		query,
		eventID,
	)

	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	return nil
}

func GetEventByID(ctx context.Context, eventID int) (*models.Event, error) {
	var event models.Event
	query := `
        SELECT id, group_id, created_by, name, description, start_time, end_time, created_at, location
        FROM events 
        WHERE id = $1
    `
	err := Pool.QueryRow(ctx, query, eventID).Scan(
		&event.ID,
		&event.GroupID,
		&event.CreatedByID,
		&event.Name,
		&event.Description,
		&event.StartTime,
		&event.EndTime,
		&event.CreatedAt,
		&event.Location,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("event not found")
		}
		return nil, fmt.Errorf("query error: %w", err)
	}
	return &event, nil
}

func GetAllEventsByUser(ctx context.Context, userID int) ([]models.Event, error) {
	query := `
		SELECT id, group_id, created_by, name, description, start_time, end_time, created_at, location
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
			&event.CreatedByID,
			&event.Name,
			&event.Description,
			&event.StartTime,
			&event.EndTime,
			&event.CreatedAt,
			&event.Location,
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
		SELECT id, group_id, created_by, name, description, start_time, end_time, created_at, location
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
			&event.CreatedByID,
			&event.Name,
			&event.Description,
			&event.StartTime,
			&event.EndTime,
			&event.CreatedAt,
			&event.Location,
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

func UpdateEventName(ctx context.Context, eventID int, eventName string) error {
	query := `
		UPDATE events
		SET name = $1
		WHERE id = $2;
	`
	_, err := Pool.Exec(
		ctx,
		query,
		eventName,
		eventID,
	)

	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	return nil
}

func UpdateEventDescription(ctx context.Context, eventID int, description string) error {
	query := `
		UPDATE events
		SET description = $1
		WHERE id = $2;
	`
	_, err := Pool.Exec(
		ctx,
		query,
		description,
		eventID,
	)

	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	return nil
}

func UpdateEventStartTime(ctx context.Context, eventID int, startTime time.Time) error {
	query := `
		UPDATE events
		SET start_time = $1
		WHERE id = $2;
	`
	_, err := Pool.Exec(
		ctx,
		query,
		startTime,
		eventID,
	)

	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	return nil
}

func UpdateEventEndTime(ctx context.Context, eventID int, endTime time.Time) error {
	query := `
		UPDATE events
		SET end_time = $1
		WHERE id = $2;
	`
	_, err := Pool.Exec(
		ctx,
		query,
		endTime,
		eventID,
	)

	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	return nil
}

func GetEventsForTomorrow(ctx context.Context, timeToUse time.Time) ([]models.Event, error) {
	tomorrow := timeToUse.Add(24 * time.Hour)
	startOfTomorrow := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())
	endOfTomorrow := startOfTomorrow.Add(24 * time.Hour)

	query := `
		SELECT id, group_id, created_by, name, description, start_time, end_time, created_at, location
		FROM events
		WHERE start_time >= $1 AND start_time < $2
		ORDER BY start_time ASC
	`

	rows, err := Pool.Query(ctx, query, startOfTomorrow, endOfTomorrow)
	if err != nil {
		return nil, fmt.Errorf("failed to query events for tomorrow: %w", err)
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.ID,
			&event.GroupID,
			&event.CreatedByID,
			&event.Name,
			&event.Description,
			&event.StartTime,
			&event.EndTime,
			&event.CreatedAt,
			&event.Location,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event row: %w", err)
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating event rows: %w", err)
	}

	return events, nil
}

func GetEventAttendance(ctx context.Context, eventID int) ([]models.EventAttendance, error) {
	query := `
        SELECT ea.id, ea.user_id, ea.event_id, ea.status, ea.created_at
        FROM event_attendance ea
        WHERE ea.event_id = $1
        ORDER BY ea.created_at DESC
    `

	rows, err := Pool.Query(ctx, query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendance for event %d: %w", eventID, err)
	}
	defer rows.Close()

	var attendances []models.EventAttendance
	for rows.Next() {
		var attendance models.EventAttendance
		err = rows.Scan(
			&attendance.ID,
			&attendance.UserID,
			&attendance.EventID,
			&attendance.Status,
			&attendance.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan attendance row: %w", err)
		}
		attendances = append(attendances, attendance)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating attendance rows: %w", err)
	}

	return attendances, nil
}

func UpdateEventAttendance(ctx context.Context, data *models.AttendanceData) error {
	// Check if record exists
	checkQuery := `
		SELECT id FROM event_attendance 
		WHERE user_id = $1 AND event_id = $2
	`
	var id int
	err := Pool.QueryRow(ctx, checkQuery, data.UserID, data.EventID).Scan(&id)

	if err == pgx.ErrNoRows {
		// Insert new record if none exists
		insertQuery := `
			INSERT INTO event_attendance (user_id, event_id, status)
			VALUES ($1, $2, $3)
		`
		_, err = Pool.Exec(ctx, insertQuery, data.UserID, data.EventID, data.Status)
		if err != nil {
			return fmt.Errorf("failed to create attendance: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check attendance: %w", err)
	} else {
		// Update existing record
		updateQuery := `
			UPDATE event_attendance 
			SET status = $1
			WHERE user_id = $2 AND event_id = $3
		`
		_, err = Pool.Exec(ctx, updateQuery, data.Status, data.UserID, data.EventID)
		if err != nil {
			return fmt.Errorf("failed to update attendance: %w", err)
		}
	}

	return nil
}
