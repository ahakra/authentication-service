package data

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"slices"
)

type PermissionsRepository struct {
	DB *sql.DB
}

func NewPermissionsRepository(db *sql.DB) *PermissionsRepository {
	return &PermissionsRepository{DB: db}
}

func (p Permissions) HasPermission(permission string) bool {
	return slices.Contains(p, permission)
}

func (m *PermissionsRepository) InsertPermissions(permission string) error {
	stmt := `INSERT INTO permissions (permission) VALUES ($1)`
	_, err := m.DB.Exec(stmt, permission)
	if err != nil {
		log.Printf("Error inserting permission: %v", err)
		return fmt.Errorf("could not insert permission: %w", err)
	}
	return nil
}

func (m *PermissionsRepository) InsertUserPermissions(userID, permissionID int64) error {
	stmt := `INSERT INTO users_permissions (user_id, permission_id) VALUES ($1, $2)`
	_, err := m.DB.Exec(stmt, userID, permissionID)
	if err != nil {
		log.Printf("Error inserting user permission: %v", err)
		return fmt.Errorf("could not insert user permission: %w", err)
	}
	return nil
}

func (m *PermissionsRepository) DeleteUserPermissions(userID, permissionID int64) error {
	stmt := `DELETE FROM users_permissions WHERE user_id = $1 AND permission_id = $2`
	_, err := m.DB.Exec(stmt, userID, permissionID)
	if err != nil {
		log.Printf("Error deleting user permission: %v", err)
		return fmt.Errorf("could not delete user permission: %w", err)
	}
	return nil
}
func (m *PermissionsRepository) GetPermissionIDByName(permission string) (int64, error) {
	query := `SELECT id FROM permissions WHERE permission = $1 LIMIT 1`

	var permissionID int64
	err := m.DB.QueryRow(query, permission).Scan(&permissionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("permission not found: %w", err)
		}
		return 0, fmt.Errorf("could not retrieve permission ID: %w", err)
	}

	return permissionID, nil
}
func (m *PermissionsRepository) GetAllForUser(userID int64) (Permissions, error) {
	query := `
        SELECT permissions.permission
        FROM permissions
        INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
        INNER JOIN users ON users_permissions.user_id = users.id
        WHERE users.id = $1`

	rows, err := m.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions Permissions

	for rows.Next() {
		var permission string

		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}
