package db

import (
	"context"
	"fmt"
	"nest/models"
)

func ReactToEvent(ctx context.Context, userID int, eventID int, reaction *models.Reaction) error {
	query := `
		INSERT INTO event_reactions (user_id, event_id, reaction)
		VALUES ($1, $2, $3)
	`

	_, err := Pool.Exec(
		ctx,
		query,
		userID,
		eventID,
		reaction,
	)

	if err != nil {
		return fmt.Errorf("failed to react **%s** to event: %w", string(*reaction), err)
	}

	return nil
}

func UnreactToEvent(ctx context.Context, userID int, eventID int, reaction *models.Reaction) error {
	query := `
		DELETE FROM event_reactions
		WHERE user_id = $1 AND event_id = $2 AND reaction = $3
	`

	_, err := Pool.Exec(
		ctx,
		query,
		userID,
		eventID,
		reaction,
	)

	if err != nil {
		return fmt.Errorf("failed to unreact **%d** to event: %w", eventID, err)
	}

	return nil
}

func GetReactionsByEvent(ctx context.Context, eventID int) ([]models.UserReaction, error) {
	query := `
		SELECT u.id AS user_id, er.reaction AS reaction, er.event_id as event_id
		FROM event_reactions er
		JOIN users u ON er.user_id = u.id
		WHERE er.event_id = $1;
	`

	rows, err := Pool.Query(ctx, query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reactions for event %d: %w", eventID, err)
	}
	defer rows.Close()

	var userReactions []models.UserReaction

	for rows.Next() {
		var userReaction models.UserReaction
		err = rows.Scan(
			&userReaction.UserID,
			&userReaction.Reaction,
			&userReaction.EventID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event reaction row: %w", err)
		}
		userReactions = append(userReactions, userReaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating event reaction rows: %w", err)
	}

	return userReactions, nil
}
