package handler

import (
	"context"
	"finance-dashboard-backend/db"
	"finance-dashboard-backend/objects"
	"fmt"
)

type GetSummaryRequest struct {
	UserID string `json:"userId"`
}

type CategoryTotal struct {
	Category string  `json:"category"`
	Type     string  `json:"type"`
	Total    float64 `json:"total"`
}

type MonthlyTrend struct {
	Month string  `json:"month"` // e.g. "2024-04"
	Type  string  `json:"type"`
	Total float64 `json:"total"`
}

type GetSummaryResponse struct {
	objects.GenericResponse
	TotalIncome    float64         `json:"totalIncome"`
	TotalExpense   float64         `json:"totalExpense"`
	NetBalance     float64         `json:"netBalance"`
	CategoryTotals []CategoryTotal `json:"categoryTotals"`
	RecentActivity []RecordDetail  `json:"recentActivity"`
	MonthlyTrends  []MonthlyTrend  `json:"monthlyTrends"`
}

func GetSummary(req objects.RequestObject) error {
	reqPay := req.Request.(*GetSummaryRequest)
	resPay := req.Response.(*GetSummaryResponse)

	// 1. Total Income & Expense
	queryTotals := `
		SELECT 
			COALESCE(SUM(CASE WHEN type = 'Income' THEN amount ELSE 0 END), 0) as total_income,
			COALESCE(SUM(CASE WHEN type = 'Expense' THEN amount ELSE 0 END), 0) as total_expense
		FROM financial_records 
		WHERE user_id = $1 AND deleted_at IS NULL
	`
	err := db.Client.QueryRow(context.Background(), queryTotals, reqPay.UserID).Scan(&resPay.TotalIncome, &resPay.TotalExpense)
	if err != nil {
		return fmt.Errorf("failed to fetch totals: %w", err)
	}
	resPay.NetBalance = resPay.TotalIncome - resPay.TotalExpense

	// 2. Category-wise Totals
	categoryQuery := `
		SELECT category, type, SUM(amount) 
		FROM financial_records 
		WHERE user_id = $1 AND deleted_at IS NULL 
		GROUP BY category, type
	`
	rows, err := db.Client.Query(context.Background(), categoryQuery, reqPay.UserID)
	if err != nil {
		return fmt.Errorf("failed to fetch category totals: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ct CategoryTotal
		if err := rows.Scan(&ct.Category, &ct.Type, &ct.Total); err != nil {
			return fmt.Errorf("failed to scan category total: %w", err)
		}
		resPay.CategoryTotals = append(resPay.CategoryTotals, ct)
	}

	// 3. Recent Activity (Last 5 records)
	recentQuery := `
		SELECT id, user_id, amount, type, category, date, notes, created_at 
		FROM financial_records 
		WHERE user_id = $1 AND deleted_at IS NULL 
		ORDER BY date DESC LIMIT 5
	`
	rows, err = db.Client.Query(context.Background(), recentQuery, reqPay.UserID)
	if err != nil {
		return fmt.Errorf("failed to fetch recent activity: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r RecordDetail
		if err := rows.Scan(&r.ID, &r.UserID, &r.Amount, &r.Type, &r.Category, &r.Date, &r.Notes, &r.CreatedAt); err != nil {
			return fmt.Errorf("failed to scan recent record: %w", err)
		}
		resPay.RecentActivity = append(resPay.RecentActivity, r)
	}

	// 4. Monthly Trends (Last 6 months)
	trendQuery := `
		SELECT TO_CHAR(date, 'YYYY-MM') as month, type, SUM(amount) 
		FROM financial_records 
		WHERE user_id = $1 AND deleted_at IS NULL 
		GROUP BY month, type 
		ORDER BY month DESC LIMIT 12
	`
	rows, err = db.Client.Query(context.Background(), trendQuery, reqPay.UserID)
	if err != nil {
		return fmt.Errorf("failed to fetch monthly trends: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var mt MonthlyTrend
		if err := rows.Scan(&mt.Month, &mt.Type, &mt.Total); err != nil {
			return fmt.Errorf("failed to scan trend: %w", err)
		}
		resPay.MonthlyTrends = append(resPay.MonthlyTrends, mt)
	}

	resPay.Success = true
	return nil
}
