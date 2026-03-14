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

	// AUTH
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", handler.Register)
		authRoutes.GET("/verify", handler.Verify)
		authRoutes.POST("/login", handler.Login)
		authRoutes.POST("/refresh", handler.Refresh)
	}

	// PROTECTED ROUTES
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{

		// USER
		api.POST("/logout", handler.Logout)
		api.GET("/roles", handler.GetRoles)
		api.POST("/roles/assign", handler.AssignRole)

		// DOCTOR
		api.POST("/doctors", handler.CreateDoctor)
		api.PUT("/doctors/:id/activate", handler.ActivateDoctor)
		api.PUT("/doctors/:id/deactivate", handler.DeactivateDoctor)

		// SERVICES
		api.POST("/clinics", handler.CreateClinic)
		api.POST("/services", handler.CreateService)
		api.GET("/services", handler.GetServices)

		// DOCTOR SERVICE
		api.POST("/doctor-services", handler.AssignServiceToDoctor)
	}

	r.Run(":8080")
}
