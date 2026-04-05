package handler

import (
	"context"
	"finance-dashboard-backend/db"
	"finance-dashboard-backend/objects"
	"fmt"
	"time"
)

type DeleteRecordRequest struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
}

type DeleteRecordResponse struct {
	objects.GenericResponse
}

func DeleteRecord(req objects.RequestObject) error {
	reqPay := req.Request.(*DeleteRecordRequest)
	resPay := req.Response.(*DeleteRecordResponse)

	if reqPay.ID == "" {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "record id is required",
			StatusCode: 400,
		}
	}

	now := time.Now()
	query := `UPDATE financial_records SET deleted_at = $1 WHERE id = $2 AND user_id = $3 AND deleted_at IS NULL`
	
	res, err := db.Client.Exec(context.Background(), query, now, reqPay.ID, reqPay.UserID)
	if err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}

	if res.RowsAffected() == 0 {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "record not found or already deleted",
			StatusCode: 404,
		}
	}

	resPay.Success = true
	return nil
}
