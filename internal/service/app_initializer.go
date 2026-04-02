package service

import (
	"fmt"
	"hospital/internal/config"
	"log"
	"time"
)

type AppInitializer struct {
	isInitialized bool
	scheduler     *JobScheduler
}

func InitializeServices(reminderInterval, retryInterval time.Duration) (*AppInitializer, error) {
	log.Println("[App Initialization] starting application services")

	initializer := &AppInitializer{}

	SetEmailProvider(NewMockEmailProvider())
	log.Println("[App Initialization] email service initialized (mock provider)")

	SetSMSProvider(NewMockSMSProvider())
	log.Println("[App Initialization] SMS service initialized (mock provider)")

	SetPushNotificationProvider(NewMockPushNotificationProvider())
	log.Println("[App Initialization] push notification service initialized (mock provider)")

	if reminderTablesReady() {
		scheduler, err := InitializeJobScheduler(reminderInterval, retryInterval)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize job scheduler: %w", err)
		}
		initializer.scheduler = scheduler
		log.Println("[App Initialization] job scheduler started successfully")
	} else {
		log.Println("[App Initialization] job scheduler skipped because reminder tables are not available yet")
	}

	initializer.isInitialized = true
	return initializer, nil
}

func (a *AppInitializer) ConfigureEmailService(host string, port int, username, password, fromAddr string) error {
	if !a.isInitialized {
		return fmt.Errorf("services not initialized yet")
	}

	SetEmailProvider(NewSMTPEmailProvider(host, port, username, password, fromAddr))
	log.Printf("[App Initialization] email service configured: %s:%d (from: %s)", host, port, fromAddr)
	return nil
}

func (a *AppInitializer) ConfigureSMSService(fromNumber string) error {
	if !a.isInitialized {
		return fmt.Errorf("services not initialized yet")
	}

	SetSMSProvider(NewTwilioSMSProvider(fromNumber))
	log.Printf("[App Initialization] SMS service configured: from %s", fromNumber)
	return nil
}

func (a *AppInitializer) ConfigurePushNotificationService(serverKey string) error {
	if !a.isInitialized {
		return fmt.Errorf("services not initialized yet")
	}

	SetPushNotificationProvider(NewFCMPushNotificationProvider(serverKey))
	log.Println("[App Initialization] push notification service configured")
	return nil
}

func (a *AppInitializer) IsInitialized() bool {
	return a.isInitialized
}

func (a *AppInitializer) GetScheduler() *JobScheduler {
	return a.scheduler
}

func (a *AppInitializer) Shutdown() error {
	if !a.isInitialized {
		return nil
	}

	if a.scheduler != nil && a.scheduler.IsRunning() {
		if err := a.scheduler.Stop(); err != nil {
			return err
		}
	}

	a.isInitialized = false
	log.Println("[App Shutdown] application services shut down successfully")
	return nil
}

var globalAppInitializer *AppInitializer

func GetAppInitializer() *AppInitializer {
	return globalAppInitializer
}

func SetAppInitializer(initializer *AppInitializer) {
	globalAppInitializer = initializer
}

func reminderTablesReady() bool {
	if config.DB == nil {
		return false
	}

	var exists bool
	err := config.DB.QueryRow(`
	SELECT EXISTS (
		SELECT 1
		FROM information_schema.tables
		WHERE table_schema = 'public' AND table_name = 'appointment_reminders'
	)
	`).Scan(&exists)
	if err != nil {
		log.Printf("[App Initialization] could not check reminder table availability: %v", err)
		return false
	}

	return exists
}
