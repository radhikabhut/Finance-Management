package handler

import (
	"context"
	"finance-dashboard-backend/db"
	"finance-dashboard-backend/objects"
	"fmt"
	"time"
)

type UpdateRecordRequest struct {
	ID       string    `json:"id"`
	UserID   string    `json:"userId"`
	Amount   float64   `json:"amount"`
	Type     string    `json:"type"` // Income, Expense
	Category string    `json:"category"`
	Date     time.Time `json:"date"`
	Notes    string    `json:"notes"`
}

type UpdateRecordResponse struct {
	objects.GenericResponse
}

func UpdateRecord(req objects.RequestObject) error {
	reqPay := req.Request.(*UpdateRecordRequest)
	resPay := req.Response.(*UpdateRecordResponse)

	if reqPay.ID == "" {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "record id is required",
			StatusCode: 400,
		}
	}

	now := time.Now()
	query := `UPDATE financial_records 
	          SET amount = CASE WHEN $1 > 0 THEN $1 ELSE amount END, 
	              type = COALESCE(NULLIF($2, ''), type), 
	              category = COALESCE(NULLIF($3, ''), category), 
	              date = CASE WHEN $4::timestamp IS NOT NULL AND $4 <> '0001-01-01 00:00:00' THEN $4 ELSE date END, 
	              notes = COALESCE(NULLIF($5, ''), notes), 
	              updated_at = $6 
	          WHERE id = $7 AND user_id = $8 AND deleted_at IS NULL`
	
	res, err := db.Client.Exec(context.Background(), query, 
		reqPay.Amount, reqPay.Type, reqPay.Category, reqPay.Date, reqPay.Notes, now, reqPay.ID, reqPay.UserID)
	
	if err != nil {
		return fmt.Errorf("failed to update record: %w", err)
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
