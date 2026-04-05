package view

import (
	"finance-dashboard-backend/handler"
	"finance-dashboard-backend/objects"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListUser(ctx *gin.Context) {
	var req objects.RequestObject
	var reqPay handler.ListUserRequest
	var resPay handler.ListUserResponse

	// Get pagination from query params or body
	page, _ := strconv.Atoi(ctx.Query("page"))
	limit, _ := strconv.Atoi(ctx.Query("limit"))

	if page > 0 {
		reqPay.Page = page
	}
	if limit > 0 {
		reqPay.Limit = limit
	}

	// Also try to bind from JSON if present
	_ = ctx.ShouldBindJSON(&reqPay)

	req.Request = &reqPay
	req.Response = &resPay

	if err := handler.ListUser(req); err != nil {
		slog.Error("error while listing users", "error", err.Error())
		
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
