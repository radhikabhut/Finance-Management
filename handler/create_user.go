package handler

import (
	"context"
	"finance-dashboard-backend/db"
	"finance-dashboard-backend/objects"
	"finance-dashboard-backend/utils"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`   // Viewer, Analyst, Admin
	Status   string `json:"status"` // Active, Inactive
}

type CreateUserResponse struct {
	objects.GenericResponse
	ID string `json:"id"`
}

func CreateUser(req objects.RequestObject) error {
	reqPay := req.Request.(*CreateUserRequest)
	resPay := req.Response.(*CreateUserResponse)

	// Validation
	if len(reqPay.Username) < 3 {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "username must be at least 3 characters long",
			StatusCode: 400,
		}
	}
	if len(reqPay.Password) < 6 {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "password must be at least 6 characters long",
			StatusCode: 400,
		}
	}
	if reqPay.Role == "" {
		reqPay.Role = "Viewer"
	}
	// Check allowed roles
	allowedRoles := map[string]bool{"Admin": true, "Analyst": true, "Viewer": true}
	if !allowedRoles[reqPay.Role] {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "invalid role. must be Admin, Analyst, or Viewer",
			StatusCode: 400,
		}
	}

	if reqPay.Status == "" {
		reqPay.Status = "Active"
	}

	// Generate UUID
	id := uuid.New().String()

	// Hash Password
	hashedPassword, err := utils.HashPassword(reqPay.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	now := time.Now()
	query := `INSERT INTO users (id, username, password_hash, role, status, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	
	_, err = db.Client.Exec(context.Background(), query, 
		id, reqPay.Username, hashedPassword, reqPay.Role, reqPay.Status, now, now)
	
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	resPay.Success = true
	resPay.ID = id
	return nil
}
