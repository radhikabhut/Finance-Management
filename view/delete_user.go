package view

import (
	"finance-dashboard-backend/handler"
	"finance-dashboard-backend/objects"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteUser(ctx *gin.Context) {
	var req objects.RequestObject
	var reqPay handler.DeleteUserRequest
	var resPay handler.DeleteUserResponse

	// Try to get ID from path or body
	reqPay.ID = ctx.Param("id")
	if reqPay.ID == "" {
		if err := ctx.ShouldBindJSON(&reqPay); err != nil {
			slog.Error("error while binding json", "error", err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, objects.MalformedError)
			return
		}
	}

	req.Request = &reqPay
	req.Response = &resPay

	if err := handler.DeleteUser(req); err != nil {
		slog.Error("error while deleting user", "error", err.Error())
		
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
