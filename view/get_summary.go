package view

import (
	"finance-dashboard-backend/handler"
	"finance-dashboard-backend/objects"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSummary(ctx *gin.Context) {
	var req objects.RequestObject
	var reqPay handler.GetSummaryRequest
	var resPay handler.GetSummaryResponse

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

	if err := handler.GetSummary(req); err != nil {
		slog.Error("error while getting summary", "error", err.Error())
		
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
