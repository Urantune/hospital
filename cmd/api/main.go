package main

import (
	"hospital/internal/config"
	"hospital/internal/handler"
	"hospital/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {

	config.ConnectDB()

	r := gin.Default()

	r.POST("/register", handler.Register)
	r.GET("/verify", handler.Verify)
	r.POST("/login", handler.Login)
	r.POST("/refresh", handler.Refresh)
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		r.POST("/logout", handler.Logout)
	}

	r.Run(":8080")

}
