package db

import (
	"context"
	"errors"
	"fmt"
	"nest/models"

	"github.com/jackc/pgx/v4"
)

func CreateGroup(ctx context.Context, group *models.Group) (*models.Group, error) {
	query := `
		INSERT INTO groups (created_by, group_name, code)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := Pool.QueryRow(
		ctx,
		query,
		group.CreatedByID,
		group.Name,
		group.Code,
	).Scan(&group.ID, &group.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	return group, nil
}

func GetGroupByID(ctx context.Context, groupID int) (*models.Group, error) {
	var group models.Group
	query := `
        SELECT id, created_by, created_at, group_name, code, do_send_emails
        FROM groups 
        WHERE id = $1
    `
	err := Pool.QueryRow(ctx, query, groupID).Scan(
		&group.ID,
		&group.CreatedByID,
		&group.CreatedAt,
		&group.Name,
		&group.Code,
		&group.DoSendEmails,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("group not found")
		}
		return nil, fmt.Errorf("query error: %w", err)
	}
	return &group, nil
}

func GetGroupByCode(ctx context.Context, groupCode string) (*models.Group, error) {
	var group models.Group
	query := `
        SELECT id, created_by, created_at, group_name, code, do_send_emails
        FROM groups 
        WHERE code = $1
    `
	err := Pool.QueryRow(ctx, query, groupCode).Scan(
		&group.ID,
		&group.CreatedByID,
		&group.CreatedAt,
		&group.Name,
		&group.Code,
		&group.DoSendEmails,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
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
		SELECT g.id, g.group_name, g.created_by, g.created_at, g.code, g.do_send_emails
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
			&group.Code,
			&group.DoSendEmails,
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

func GetAllGroups(ctx context.Context) ([]models.Group, error) {
	query := `
		SELECT id, group_name, created_by, created_at, code, do_send_emails
		FROM groups
	`

	rows, err := Pool.Query(ctx, query)
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
			&group.Code,
			&group.DoSendEmails,
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
		SELECT u.id, u.username, u.email, u.first_name, u.last_name, u.role, u.created_at
		FROM users u
		JOIN group_memberships gm ON u.id = gm.user_id
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
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt)
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

func GetAllNonMembersForGroup(ctx context.Context, groupID int) ([]models.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.first_name, u.last_name, u.role, u.created_at
		FROM users u
		LEFT JOIN group_memberships gm ON u.id = gm.user_id AND gm.group_id = $1
		WHERE gm.user_id IS NULL;
	`

	rows, err := Pool.Query(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt)
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

func GetAllNonAdminMembersForGroup(ctx context.Context, groupID int) ([]models.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.first_name, u.last_name, u.role, u.created_at
		FROM users u
		JOIN group_memberships gm ON u.id = gm.user_id
		WHERE gm.group_id = $1 and gm.role_in_group = 'member';
	`

	rows, err := Pool.Query(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt)
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

func GetAllAdminMembersForGroup(ctx context.Context, groupID int) ([]models.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.first_name, u.last_name, u.role, u.created_at
		FROM users u
		JOIN group_memberships gm ON u.id = gm.user_id
		WHERE gm.group_id = $1 and gm.role_in_group != 'member';
	`

	rows, err := Pool.Query(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt)
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
		if err == pgx.ErrNoRows {
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
		if err == pgx.ErrNoRows {
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

func UpdateGroupDoSendEmails(ctx context.Context, groupID int, doSendEmails bool) error {
	query := `
		UPDATE groups
		SET do_send_emails = $1
		WHERE id = $2;
	`
	_, err := Pool.Exec(
		ctx,
		query,
		doSendEmails,
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
