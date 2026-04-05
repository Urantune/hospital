package main

import (
	"hospital/internal/config"
	"hospital/internal/handler"
	"hospital/internal/middleware"
	"hospital/internal/service"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	config.ConnectDB()

	service.StartExpireJob()
	service.StartNotificationWorker()

	appInit, err := service.InitializeServices(1*time.Minute, 5*time.Minute)
	if err != nil {
		log.Fatalf("failed to initialize services: %v", err)
	}
	service.SetAppInitializer(appInit)
	defer appInit.Shutdown()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("[Main] received shutdown signal")
		_ = appInit.Shutdown()
		os.Exit(0)
	}()

	go func() {
		for {
			err := service.RunReconciliation()
			if err != nil {
				log.Println("[Reconciliation Error]:", err)
			}
			time.Sleep(5 * time.Minute)
		}
	}()

	r := gin.Default()

	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", handler.Register)
		authRoutes.GET("/verify", handler.Verify)
		authRoutes.POST("/login", handler.Login)
		authRoutes.POST("/refresh", handler.Refresh)
	}

	r.POST("/payments/callback", handler.HandlePaymentCallback)

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

		// CLINIC + SERVICE
		api.POST("/clinics", handler.CreateClinic)
		api.POST("/services", handler.CreateService)
		api.GET("/services", handler.GetServices)
		api.POST("/pricing/preview", handler.PreviewPrice)

		// DOCTOR SERVICE
		api.POST("/doctor-services", handler.AssignServiceToDoctor)

		// INSURANCE
		api.POST("/insurance/plans", handler.CreateInsurancePlan)
		api.GET("/insurance/plans/:id", handler.GetInsurancePlan)
		api.PUT("/insurance/plans/:id", handler.UpdateInsurancePlan)
		api.GET("/insurance/plans", handler.ListInsurancePlans)
		api.POST("/insurance/plans/:id/service-coverage", handler.SetServiceCoverage)

		// APPOINTMENTS
		api.POST("/appointments", handler.CreateAppointment)
		api.GET("/appointments/me", handler.GetMyAppointments)
		api.GET("/appointments/:id", handler.GetAppointment)
		api.GET("/appointments/:id/history", handler.GetAppointmentHistory)
		api.PATCH("/appointments/:id/status", handler.UpdateAppointmentStatus)

		api.POST("/appointments/:id/payments/initiate", handler.CreateAppointmentPayment)
		api.GET("/appointments/:id/financial-snapshot", handler.GetAppointmentFinancialSnapshot)
		api.GET("/appointments/:id/events", handler.GetAppointmentDomainEvents)

		api.GET("/payments/:id/history", handler.GetPaymentHistory)

		// SCHEDULES
		api.POST("/schedules", handler.CreateSchedule)
		api.POST("/slots/generate", handler.GenerateSlots)
		api.GET("/slots/available", handler.GetAvailableSlots)
		api.POST("/schedules/:id/impact", handler.CheckScheduleImpact)
		api.PUT("/schedules/:id/enforce", handler.UpdateScheduleWithEnforcement)
		api.POST("/schedules/:id/exceptions", handler.CreateExceptionDay)
		api.POST("/slots/:id/lock", handler.LockSlot)

		// MONITORING & HEALTH
		api.GET("/monitoring/dashboard", handler.GetSystemHealthDashboard)
		api.GET("/monitoring/appointments", handler.GetAppointmentMetrics)
		api.GET("/monitoring/payments", handler.GetPaymentMetrics)
		api.GET("/monitoring/reminders", handler.GetReminderMetrics)
		api.GET("/monitoring/clinic", handler.GetClinicMetrics)
		api.GET("/monitoring/users", handler.GetUserMetrics)
		api.GET("/monitoring/metrics-by-range", handler.GetMetricsForDateRange)

		// REMINDERS
		api.GET("/reminders/:id", handler.GetReminder)
		api.GET("/reminders/appointment/:appointment_id", handler.ListReminders)
		api.POST("/reminders/:id/cancel", handler.CancelReminder)

		// AUDIT
		api.GET("/audit-logs", handler.ListAuditLogs)
		api.GET("/activities", handler.GetActivityLogs)

		// CMS
		api.POST("/cms/sync", handler.SyncCMSChange)
	}

	r.Run(":8080")
}
