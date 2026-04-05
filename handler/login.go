package handler

import (
	"context"
	"finance-dashboard-backend/db"
	"finance-dashboard-backend/objects"
	"finance-dashboard-backend/utils"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	objects.GenericResponse
	Token string `json:"token"`
	Role  string `json:"role"`
}

func Login(req objects.RequestObject) error {
	reqPay := req.Request.(*LoginRequest)
	resPay := req.Response.(*LoginResponse)

	// Fetch user from DB
	var id, role, passwordHash string
	query := "SELECT id, role, password_hash FROM users WHERE username = $1 AND deleted_at IS NULL AND status = 'Active'"
	err := db.Client.QueryRow(context.Background(), query, reqPay.Username).Scan(&id, &role, &passwordHash)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return objects.GenericError{
				Success:    false,
				ErrMsg:     "invalid username or password",
				StatusCode: 401,
			}
		}
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	// Verify password
	if !utils.CheckPassword(reqPay.Password, passwordHash) {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "invalid username or password",
			StatusCode: 401,
		}
	}

	// Generate JWT
	token, err := utils.GenerateToken(id, role)
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	resPay.Success = true
	resPay.Token = token
	resPay.Role = role
	return nil
}
