package db

import (
	"context"
)

// HasPermission checks if a given role has a specific permission in the database.
func HasPermission(role, permission string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM role_permissions WHERE role = $1 AND permission = $2"
	err := Client.QueryRow(context.Background(), query, role, permission).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
