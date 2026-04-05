package main

import (
	"finance-dashboard-backend/middleware"
	"finance-dashboard-backend/view"

	"github.com/gin-gonic/gin"
)

func defineRoutes(v1 *gin.RouterGroup) {
	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Public Routes
	v1.POST("/login", view.Login)
	// v1.POST("/users/create", view.CreateUser) // Temporarily public for testing

	// Authenticated Group
	auth := v1.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		// User Management
		auth.POST("/users/create", middleware.RBACMiddleware("user::create"), view.CreateUser)
		auth.POST("/users/list", middleware.RBACMiddleware("user::list"), view.ListUser)
		auth.POST("/users/update", middleware.RBACMiddleware("user::update"), view.UpdateUser)
		auth.POST("/users/delete", middleware.RBACMiddleware("user::delete"), view.DeleteUser)

		// Financial Records
		auth.POST("/records/create", middleware.RBACMiddleware("record::create"), view.CreateRecord)
		auth.POST("/records/list", middleware.RBACMiddleware("record::list"), view.ListRecord)
		auth.POST("/records/update", middleware.RBACMiddleware("record::update"), view.UpdateRecord)
		auth.POST("/records/delete", middleware.RBACMiddleware("record::delete"), view.DeleteRecord)

		// Summary
		auth.POST("/summary", middleware.RBACMiddleware("summary::view"), view.GetSummary)
	}
}
