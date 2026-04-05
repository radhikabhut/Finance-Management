package handler

import (
	"context"
	"finance-dashboard-backend/db"
	"finance-dashboard-backend/objects"
	"fmt"
	"time"
)

type DeleteUserRequest struct {
	ID string `json:"id"`
}

type DeleteUserResponse struct {
	objects.GenericResponse
}

func DeleteUser(req objects.RequestObject) error {
	reqPay := req.Request.(*DeleteUserRequest)
	resPay := req.Response.(*DeleteUserResponse)

	if reqPay.ID == "" {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "user id is required",
			StatusCode: 400,
		}
	}

	now := time.Now()
	query := `UPDATE users SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	
	res, err := db.Client.Exec(context.Background(), query, now, reqPay.ID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
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
