package main

import (
	"finance-dashboard-backend/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	v1 := r.Group("/v1")
	defineRoutes(v1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/ping", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	
	// Protected route
	auth := r.Group("/auth")
	auth.Use(middleware.AuthMiddleware())
	auth.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// 1. Test Without Token
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/auth/test", nil)
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusUnauthorized, w1.Code)

	// 2. Test Invalid Token
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/auth/test", nil)
	req2.Header.Set("Authorization", "Bearer invalid_token")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusUnauthorized, w2.Code)
}
