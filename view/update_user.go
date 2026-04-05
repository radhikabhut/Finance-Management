package view

import (
	"finance-dashboard-backend/handler"
	"finance-dashboard-backend/objects"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateUser(ctx *gin.Context) {
	var req objects.RequestObject
	var reqPay handler.UpdateUserRequest
	var resPay handler.UpdateUserResponse

	if err := ctx.ShouldBindJSON(&reqPay); err != nil {
		slog.Error("error while binding json", "error", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, objects.MalformedError)
		return
	}

	// ID can also be in the URL path if needed, but here we expect it in body
	if reqPay.ID == "" {
		reqPay.ID = ctx.Param("id")
	}

	req.Request = &reqPay
	req.Response = &resPay

	if err := handler.UpdateUser(req); err != nil {
		slog.Error("error while updating user", "error", err.Error())
		
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
