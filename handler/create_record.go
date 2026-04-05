package handler

import (
	"context"
	"finance-dashboard-backend/db"
	"finance-dashboard-backend/objects"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CreateRecordRequest struct {
	UserID   string    `json:"userId"`
	Amount   float64   `json:"amount"`
	Type     string    `json:"type"` // Income, Expense
	Category string    `json:"category"`
	Date     time.Time `json:"date"`
	Notes    string    `json:"notes"`
}

type CreateRecordResponse struct {
	objects.GenericResponse
	ID string `json:"id"`
}

func CreateRecord(req objects.RequestObject) error {
	reqPay := req.Request.(*CreateRecordRequest)
	resPay := req.Response.(*CreateRecordResponse)

	// Validation
	if reqPay.Amount <= 0 {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "amount must be greater than 0",
			StatusCode: 400,
		}
	}

	// Check allowed types
	allowedTypes := map[string]bool{"Income": true, "Expense": true}
	if !allowedTypes[reqPay.Type] {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "invalid type. must be Income or Expense",
			StatusCode: 400,
		}
	}

	if reqPay.Category == "" {
		return objects.GenericError{
			Success:    false,
			ErrMsg:     "category is required",
			StatusCode: 400,
		}
	}
	if reqPay.Date.IsZero() {
		reqPay.Date = time.Now()
	}

	// Generate UUID
	id := uuid.New().String()

	now := time.Now()
	query := `INSERT INTO financial_records (id, user_id, amount, type, category, date, notes, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	
	_, err := db.Client.Exec(context.Background(), query, 
		id, reqPay.UserID, reqPay.Amount, reqPay.Type, reqPay.Category, reqPay.Date, reqPay.Notes, now, now)
	
	if err != nil {
		return fmt.Errorf("failed to insert financial record: %w", err)
	}

	resPay.Success = true
	resPay.ID = id
	return nil
}
