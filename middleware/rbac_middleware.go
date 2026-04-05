package middleware

import (
	"finance-dashboard-backend/db"
	"finance-dashboard-backend/objects"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RBACMiddleware(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, objects.GenericResponse{
				Success: false,
				ErrMsg:  "user role not found in context",
			})
			return
		}

		hasPerm, err := db.HasPermission(role.(string), permission)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, objects.GenericResponse{
				Success: false,
				ErrMsg:  "failed to check permissions",
			})
			return
		}

		if !hasPerm {
			c.AbortWithStatusJSON(http.StatusForbidden, objects.GenericResponse{
				Success: false,
				ErrMsg:  "permission denied",
			})
			return
		}

		c.Next()
	}
}
