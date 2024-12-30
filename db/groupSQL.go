package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"nest/models"
)

func CreateGroup(ctx context.Context, group *models.Group) (*models.Group, error) {
	query := `
		INSERT INTO groups (created_by, group_name)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	err := Pool.QueryRow(
		ctx,
		query,
		group.CreatedByID,
		group.Name,
	).Scan(&group.ID, &group.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	return group, nil
}

func GetGroupByID(ctx context.Context, groupID int) (*models.Group, error) {
	var group models.Group
	query := `
        SELECT id, created_by, created_at, group_name
        FROM groups 
        WHERE id = $1
    `
	err := Pool.QueryRow(ctx, query, groupID).Scan(
		&group.ID,
		&group.CreatedByID,
		&group.Name,
		&group.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("group not found")
		}
		return nil, fmt.Errorf("query error: %w", err)
	}
	return &group, nil
}

func AddGroupMember(ctx context.Context, userID, groupID int, roleInGroup models.Role) error {
	// validate user exists
	_, err := GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user with ID %d does not exist", userID)
	}

	// validate group exists
	_, err = GetGroupByID(ctx, groupID)
	if err != nil {
		return fmt.Errorf("group with ID %d does not exist", groupID)
	}

	query := `
		INSERT INTO group_memberships (group_id, user_id, role_in_group)
		VALUES ($1, $2, $3)
	`

	_, err = Pool.Exec(
		ctx,
		query,
		groupID,
		userID,
		roleInGroup,
	)

	if err != nil {
		return fmt.Errorf("failed to add user to group: %w", err)
	}

	return nil
}

func RemoveGroupMember(ctx context.Context, userID, groupID int) error {
	query := `
		DELETE FROM group_memberships
		WHERE group_id = $2 AND user_id = $1;
	`
	_, err := Pool.Exec(
		ctx,
		query,
		userID,
		groupID,
	)

	if err != nil {
		return fmt.Errorf("failed to remove user from group: %w", err)
	}

	return nil
}

func GetAllGroupsForUser(ctx context.Context, userID int) ([]models.Group, error) {
	query := `
		SELECT g.id, g.group_name, g.created_by, g.created_at
		FROM groups g
		INNER JOIN group_memberships gm ON g.id = gm.group_id
		WHERE gm.user_id = $1
	`

	rows, err := Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var groups []models.Group

	for rows.Next() {
		var group models.Group
		err := rows.Scan(
			&group.ID,
			&group.Name,
			&group.CreatedByID,
			&group.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan group row: %w", err)
		}

		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating group rows: %w", err)
	}

	return groups, nil
}

func GetAllMembersForGroup(ctx context.Context, groupID int) ([]models.User, error) {
	query := `
		SELECT u.user_id, u.username, gm.role_in_group, gm.joined_at
		FROM users u
		JOIN group_memberships gm ON u.user_id = gm.user_id
		WHERE gm.group_id = $1;
	`

	rows, err := Pool.Query(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Role, &user.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return users, nil
}

func IsUserGroupAdmin(ctx context.Context, userID, groupID int) (bool, error) {
	query := `
		SELECT 1
		FROM group_memberships
		WHERE group_id = $2 AND user_id = $1 AND role_in_group = 'group_admin';
	`
	var result int

	err := Pool.QueryRow(ctx, query, userID, groupID).Scan(&result)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("query error: %w", err)
	}
	return true, nil
}

func IsUserGroupMember(ctx context.Context, userID, groupID int) (bool, error) {
	query := `
		SELECT 1
		FROM group_memberships
		WHERE group_id = $2 AND user_id = $1;
	`
	var result int

	err := Pool.QueryRow(ctx, query, userID, groupID).Scan(&result)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("query error: %w", err)
	}
	return true, nil
}

func DeleteGroup(ctx context.Context, groupID int) error {
	query := `
		DELETE FROM groups
		WHERE id = $1;
	`
	_, err := Pool.Exec(
		ctx,
		query,
		groupID,
	)

	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}

	return nil
}

func UpdateGroupName(ctx context.Context, groupID int, groupName string) error {
	query := `
		UPDATE groups
		SET group_name = $1
		WHERE id = $2;
	`
	_, err := Pool.Exec(
		ctx,
		query,
		groupName,
		groupID,
	)

	if err != nil {
		return fmt.Errorf("failed to update group: %w", err)
	}

	return nil
}

func AddGroupAdmin(ctx context.Context, groupID, userID int) error {
	query := `
		UPDATE group_memberships
		SET role_in_group = 'group_admin'
		WHERE group_id = $1 AND user_id = $2;
	`
	_, err := Pool.Exec(
		ctx,
		query,
		groupID,
		userID,
	)

	if err != nil {
		return fmt.Errorf("failed to add group admin: %w", err)
	}

	return nil
}

func RemoveGroupAdmin(ctx context.Context, groupID, userID int) error {
	query := `
		UPDATE group_memberships
		SET role_in_group = 'member'
		WHERE group_id = $1 AND user_id = $2;
	`
	_, err := Pool.Exec(
		ctx,
		query,
		groupID,
		userID,
	)

	if err != nil {
		return fmt.Errorf("failed to remove group admin: %w", err)
	}

	return nil
}
