package view

import (
	"finance-dashboard-backend/handler"
	"finance-dashboard-backend/objects"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {
	var req objects.RequestObject
	var reqPay handler.LoginRequest
	var resPay handler.LoginResponse

	if err := ctx.ShouldBindJSON(&reqPay); err != nil {
		slog.Error("error while binding json", "error", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, objects.MalformedError)
		return
	}

	req.Request = &reqPay
	req.Response = &resPay

	if err := handler.Login(req); err != nil {
		slog.Error("error while logging in", "error", err.Error())
		
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
