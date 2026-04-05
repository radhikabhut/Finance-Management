package view

import (
	"finance-dashboard-backend/handler"
	"finance-dashboard-backend/objects"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func ListRecord(ctx *gin.Context) {
	var req objects.RequestObject
	var reqPay handler.ListRecordRequest
	var resPay handler.ListRecordResponse

	// Bind from JSON if present
	_ = ctx.ShouldBindJSON(&reqPay)

	// Override with Query Params if present
	if page, _ := strconv.Atoi(ctx.Query("page")); page > 0 {
		reqPay.Page = page
	}
	if limit, _ := strconv.Atoi(ctx.Query("limit")); limit > 0 {
		reqPay.Limit = limit
	}
	if t := ctx.Query("type"); t != "" {
		reqPay.Type = t
	}
	if cat := ctx.Query("category"); cat != "" {
		reqPay.Category = cat
	}
	if sd := ctx.Query("startDate"); sd != "" {
		reqPay.StartDate, _ = time.Parse(time.RFC3339, sd)
	}
	if ed := ctx.Query("endDate"); ed != "" {
		reqPay.EndDate, _ = time.Parse(time.RFC3339, ed)
	}

	// Extract UserID from context
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, objects.GenericResponse{
			Success: false,
			ErrMsg:  "user id not found in context",
		})
		return
	}
	reqPay.UserID = userID.(string)

	req.Request = &reqPay
	req.Response = &resPay

	if err := handler.ListRecord(req); err != nil {
		slog.Error("error while listing records", "error", err.Error())
		
		genericErr, ok := err.(objects.GenericError)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, objects.InternalError)
			return
		}
		
		ctx.AbortWithStatusJSON(genericErr.StatusCode, genericErr)
		return
	}

	ctx.JSON(http.StatusOK, resPay)
}
