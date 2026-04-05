package handler

import (
	"context"
	"finance-dashboard-backend/db"
	"finance-dashboard-backend/objects"
	"fmt"
	"time"
)

type UpdateUserRequest struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

type UpdateUserResponse struct {
	objects.GenericResponse
}

func UpdateUser(req objects.RequestObject) error {
	reqPay := req.Request.(*UpdateUserRequest)
	resPay := req.Response.(*UpdateUserResponse)

	if reqPay.ID == "" {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "user id is required",
			StatusCode: 400,
		}
	}

	now := time.Now()
	query := `UPDATE users 
	          SET username = COALESCE(NULLIF($1, ''), username), 
	              role = COALESCE(NULLIF($2, ''), role), 
	              status = COALESCE(NULLIF($3, ''), status), 
	              updated_at = $4 
	          WHERE id = $5 AND deleted_at IS NULL`
	
	res, err := db.Client.Exec(context.Background(), query, 
		reqPay.Username, reqPay.Role, reqPay.Status, now, reqPay.ID)
	
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if res.RowsAffected() == 0 {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "user not found or already deleted",
			StatusCode: 404,
		}
	}

	resPay.Success = true
	return nil
}
