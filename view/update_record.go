package view

import (
	"finance-dashboard-backend/handler"
	"finance-dashboard-backend/objects"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateRecord(ctx *gin.Context) {
	var req objects.RequestObject
	var reqPay handler.UpdateRecordRequest
	var resPay handler.UpdateRecordResponse

	if err := ctx.ShouldBindJSON(&reqPay); err != nil {
		slog.Error("error while binding json", "error", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, objects.MalformedError)
		return
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

	if err := handler.UpdateRecord(req); err != nil {
		slog.Error("error while updating record", "error", err.Error())
		
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
