package service

import (
	"fmt"
	"log"
)

type PushNotificationProvider interface {
	SendAppointmentReminderNotification(deviceToken, patientName, doctorName, appointmentTime string) error
	SendPaymentConfirmationNotification(deviceToken, transactionRef string, amount float64) error
	SendAppointmentConfirmationNotification(deviceToken, appointmentRef, appointmentTime string) error
	SendAppointmentCancelledNotification(deviceToken, appointmentRef, reason string) error
}

type FCMPushNotificationProvider struct {
	serverKey string
}

func NewFCMPushNotificationProvider(serverKey string) *FCMPushNotificationProvider {
	return &FCMPushNotificationProvider{serverKey: serverKey}
}

func (p *FCMPushNotificationProvider) SendAppointmentReminderNotification(deviceToken, patientName, doctorName, appointmentTime string) error {
	title := fmt.Sprintf("Appointment Reminder - Dr. %s", doctorName)
	body := fmt.Sprintf("Hello %s, your appointment is scheduled for %s. Please confirm your attendance.", patientName, appointmentTime)
	log.Printf("[FCM Push] to=%s title=%s body=%s", deviceToken, title, body)
	return nil
}

func (p *FCMPushNotificationProvider) SendPaymentConfirmationNotification(deviceToken, transactionRef string, amount float64) error {
	title := "Payment Confirmed"
	body := fmt.Sprintf("Your payment of %.2f has been successfully processed. Transaction ID: %s", amount, transactionRef)
	log.Printf("[FCM Push] to=%s title=%s body=%s", deviceToken, title, body)
	return nil
}

func (p *FCMPushNotificationProvider) SendAppointmentConfirmationNotification(deviceToken, appointmentRef, appointmentTime string) error {
	title := "Appointment Confirmed"
	body := fmt.Sprintf("Your appointment (%s) is confirmed for %s. Please arrive 15 minutes early.", appointmentRef, appointmentTime)
	log.Printf("[FCM Push] to=%s title=%s body=%s", deviceToken, title, body)
	return nil
}

func (p *FCMPushNotificationProvider) SendAppointmentCancelledNotification(deviceToken, appointmentRef, reason string) error {
	title := "Appointment Cancelled"
	body := fmt.Sprintf("Your appointment (%s) has been cancelled. Reason: %s", appointmentRef, reason)
	log.Printf("[FCM Push] to=%s title=%s body=%s", deviceToken, title, body)
	return nil
}

type MockPushNotificationProvider struct {
	sentNotifications []string
}

func NewMockPushNotificationProvider() *MockPushNotificationProvider {
	return &MockPushNotificationProvider{sentNotifications: []string{}}
}

func (p *MockPushNotificationProvider) SendAppointmentReminderNotification(deviceToken, patientName, doctorName, appointmentTime string) error {
	message := fmt.Sprintf("[MOCK Push] Reminder to %s: Dr. %s at %s", patientName, doctorName, appointmentTime)
	p.sentNotifications = append(p.sentNotifications, message)
	log.Printf("mock push sent: %s", message)
	return nil
}

func (p *MockPushNotificationProvider) SendPaymentConfirmationNotification(deviceToken, transactionRef string, amount float64) error {
	message := fmt.Sprintf("[MOCK Push] Payment confirmed: %s for %.2f", transactionRef, amount)
	p.sentNotifications = append(p.sentNotifications, message)
	log.Printf("mock push sent: %s", message)
	return nil
}

func (p *MockPushNotificationProvider) SendAppointmentConfirmationNotification(deviceToken, appointmentRef, appointmentTime string) error {
	message := fmt.Sprintf("[MOCK Push] Appointment confirmed: %s at %s", appointmentRef, appointmentTime)
	p.sentNotifications = append(p.sentNotifications, message)
	log.Printf("mock push sent: %s", message)
	return nil
}

func (p *MockPushNotificationProvider) SendAppointmentCancelledNotification(deviceToken, appointmentRef, reason string) error {
	message := fmt.Sprintf("[MOCK Push] Appointment cancelled: %s - %s", appointmentRef, reason)
	p.sentNotifications = append(p.sentNotifications, message)
	log.Printf("mock push sent: %s", message)
	return nil
}

var pushProvider PushNotificationProvider = NewMockPushNotificationProvider()

func SetPushNotificationProvider(provider PushNotificationProvider) {
	if provider != nil {
		pushProvider = provider
	}
}

func GetPushNotificationProvider() PushNotificationProvider {
	return pushProvider
}
