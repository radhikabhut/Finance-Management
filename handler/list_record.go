package handler

import (
	"context"
	"finance-dashboard-backend/db"
	"finance-dashboard-backend/objects"
	"fmt"
	"time"
)

type ListRecordRequest struct {
	UserID   string    `json:"userId"`
	Type     string    `json:"type"`     // Income, Expense
	Category string    `json:"category"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Page      int       `json:"page"`
	Limit     int       `json:"limit"`
}

type RecordDetail struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Amount    float64   `json:"amount"`
	Type      string    `json:"type"`
	Category  string    `json:"category"`
	Date      time.Time `json:"date"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"createdAt"`
}

type ListRecordResponse struct {
	objects.GenericResponse
	Records []RecordDetail `json:"records"`
	Total   int            `json:"total"`
}

func ListRecord(req objects.RequestObject) error {
	reqPay := req.Request.(*ListRecordRequest)
	resPay := req.Response.(*ListRecordResponse)

	if reqPay.Limit <= 0 {
		reqPay.Limit = 10
	}
	if reqPay.Page <= 0 {
		reqPay.Page = 1
	}
	offset := (reqPay.Page - 1) * reqPay.Limit

	// Build query with filters
	query := `SELECT id, user_id, amount, type, category, date, notes, created_at 
	          FROM financial_records 
	          WHERE deleted_at IS NULL AND user_id = $1`
	countQuery := `SELECT COUNT(*) FROM financial_records WHERE deleted_at IS NULL AND user_id = $1`
	
	args := []interface{}{reqPay.UserID}
	argIdx := 2

	if reqPay.Type != "" {
		query += fmt.Sprintf(" AND type = $%d", argIdx)
		countQuery += fmt.Sprintf(" AND type = $%d", argIdx)
		args = append(args, reqPay.Type)
		argIdx++
	}
	if reqPay.Category != "" {
		query += fmt.Sprintf(" AND category = $%d", argIdx)
		countQuery += fmt.Sprintf(" AND category = $%d", argIdx)
		args = append(args, reqPay.Category)
		argIdx++
	}
	if !reqPay.StartDate.IsZero() {
		query += fmt.Sprintf(" AND date >= $%d", argIdx)
		countQuery += fmt.Sprintf(" AND date >= $%d", argIdx)
		args = append(args, reqPay.StartDate)
		argIdx++
	}
	if !reqPay.EndDate.IsZero() {
		query += fmt.Sprintf(" AND date <= $%d", argIdx)
		countQuery += fmt.Sprintf(" AND date <= $%d", argIdx)
		args = append(args, reqPay.EndDate)
		argIdx++
	}

	// Count total
	var total int
	err := db.Client.QueryRow(context.Background(), countQuery, args...).Scan(&total)
	if err != nil {
		return fmt.Errorf("failed to count records: %w", err)
	}

	// Add ordering and pagination
	query += fmt.Sprintf(" ORDER BY date DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, reqPay.Limit, offset)

	rows, err := db.Client.Query(context.Background(), query, args...)
	if err != nil {
		return fmt.Errorf("failed to fetch records: %w", err)
	}
	defer rows.Close()

	var records []RecordDetail
	for rows.Next() {
		var r RecordDetail
		if err := rows.Scan(&r.ID, &r.UserID, &r.Amount, &r.Type, &r.Category, &r.Date, &r.Notes, &r.CreatedAt); err != nil {
			return fmt.Errorf("failed to scan record: %w", err)
		}
		records = append(records, r)
	}

	resPay.Success = true
	resPay.Records = records
	resPay.Total = total
	return nil
}
