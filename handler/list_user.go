package handler

import (
	"context"
	"finance-dashboard-backend/db"
	"finance-dashboard-backend/objects"
	"fmt"
	"time"
)

type ListUserRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type UserDetail struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

type ListUserResponse struct {
	objects.GenericResponse
	Users []UserDetail `json:"users"`
	Total int          `json:"total"`
}

func ListUser(req objects.RequestObject) error {
	reqPay := req.Request.(*ListUserRequest)
	resPay := req.Response.(*ListUserResponse)

	if reqPay.Limit <= 0 {
		reqPay.Limit = 10
	}
	if reqPay.Page <= 0 {
		reqPay.Page = 1
	}
	offset := (reqPay.Page - 1) * reqPay.Limit

	// Count total
	var total int
	err := db.Client.QueryRow(context.Background(), "SELECT COUNT(*) FROM users WHERE deleted_at IS NULL").Scan(&total)
	if err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}

	// Fetch users
	query := `SELECT id, username, role, status, created_at 
	          FROM users 
	          WHERE deleted_at IS NULL 
	          ORDER BY created_at DESC 
	          LIMIT $1 OFFSET $2`
	
	rows, err := db.Client.Query(context.Background(), query, reqPay.Limit, offset)
	if err != nil {
		return fmt.Errorf("failed to fetch users: %w", err)
	}
	defer rows.Close()

	var users []UserDetail
	for rows.Next() {
		var u UserDetail
		if err := rows.Scan(&u.ID, &u.Username, &u.Role, &u.Status, &u.CreatedAt); err != nil {
			return fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, u)
	}

	resPay.Success = true
	resPay.Users = users
	resPay.Total = total
	return nil
}
