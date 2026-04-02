package service

import (
	"fmt"
	"log"
)

type SMSProvider interface {
	SendAppointmentReminderSMS(phoneNumber, patientName, doctorName, appointmentTime string) error
	SendPaymentConfirmationSMS(phoneNumber, transactionRef string, amount float64) error
	SendAppointmentConfirmationSMS(phoneNumber, appointmentRef, appointmentTime string) error
}

type TwilioSMSProvider struct {
	fromNumber string
}

func NewTwilioSMSProvider(fromNumber string) *TwilioSMSProvider {
	return &TwilioSMSProvider{fromNumber: fromNumber}
}

func (p *TwilioSMSProvider) SendAppointmentReminderSMS(phoneNumber, patientName, doctorName, appointmentTime string) error {
	message := fmt.Sprintf(
		"Hello %s, reminder: Your appointment with Dr. %s is scheduled for %s. Reply to confirm.",
		patientName,
		doctorName,
		appointmentTime,
	)

	log.Printf("[TWILIO SMS] to=%s message=%s", phoneNumber, message)
	return nil
}

func (p *TwilioSMSProvider) SendPaymentConfirmationSMS(phoneNumber, transactionRef string, amount float64) error {
	message := fmt.Sprintf(
		"Payment confirmed! Transaction Ref: %s | Amount: %.2f | Your appointment is confirmed. Thank you!",
		transactionRef,
		amount,
	)

	log.Printf("[TWILIO SMS] to=%s message=%s", phoneNumber, message)
	return nil
}

func (p *TwilioSMSProvider) SendAppointmentConfirmationSMS(phoneNumber, appointmentRef, appointmentTime string) error {
	message := fmt.Sprintf(
		"Appointment confirmed! Ref: %s | Scheduled for: %s | Please arrive 15 minutes early.",
		appointmentRef,
		appointmentTime,
	)

	log.Printf("[TWILIO SMS] to=%s message=%s", phoneNumber, message)
	return nil
}

type MockSMSProvider struct {
	sentMessages []string
}

func NewMockSMSProvider() *MockSMSProvider {
	return &MockSMSProvider{sentMessages: []string{}}
}

func (p *MockSMSProvider) SendAppointmentReminderSMS(phoneNumber, patientName, doctorName, appointmentTime string) error {
	message := fmt.Sprintf(
		"[MOCK SMS] Reminder to %s: Appointment with Dr. %s at %s",
		patientName,
		doctorName,
		appointmentTime,
	)
	p.sentMessages = append(p.sentMessages, message)
	log.Printf("mock sms sent: %s", message)
	return nil
}

func (p *MockSMSProvider) SendPaymentConfirmationSMS(phoneNumber, transactionRef string, amount float64) error {
	message := fmt.Sprintf("[MOCK SMS] Payment confirmed: %s for %.2f", transactionRef, amount)
	p.sentMessages = append(p.sentMessages, message)
	log.Printf("mock sms sent: %s", message)
	return nil
}

func (p *MockSMSProvider) SendAppointmentConfirmationSMS(phoneNumber, appointmentRef, appointmentTime string) error {
	message := fmt.Sprintf("[MOCK SMS] Appointment confirmed: %s at %s", appointmentRef, appointmentTime)
	p.sentMessages = append(p.sentMessages, message)
	log.Printf("mock sms sent: %s", message)
	return nil
}

var smsProvider SMSProvider = NewMockSMSProvider()

func SetSMSProvider(provider SMSProvider) {
	if provider != nil {
		smsProvider = provider
	}
}

func GetSMSProvider() SMSProvider {
	return smsProvider
}
