package service

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"time"
)

type EmailProvider interface {
	SendEmail(to, subject, body string) error
}

type SMTPEmailProvider struct {
	Host     string
	Port     int
	Username string
	Password string
	FromAddr string
}

func (p *SMTPEmailProvider) SendEmail(to, subject, body string) error {
	if !strings.Contains(to, "@") {
		return fmt.Errorf("invalid email address: %s", to)
	}

	auth := smtp.PlainAuth("", p.Username, p.Password, p.Host)
	addr := fmt.Sprintf("%s:%d", p.Host, p.Port)

	headers := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/plain; charset=\"UTF-8\";\r\n\r\n",
		p.FromAddr,
		to,
		subject,
	)

	if err := smtp.SendMail(addr, auth, p.FromAddr, []string{to}, []byte(headers+body)); err != nil {
		log.Printf("failed to send email to %s: %v", to, err)
		return err
	}

	log.Printf("email sent successfully to %s", to)
	return nil
}

type MockEmailProvider struct {
	SentEmails []EmailRecord
}

type EmailRecord struct {
	To      string
	Subject string
	Body    string
}

func (p *MockEmailProvider) SendEmail(to, subject, body string) error {
	if !strings.Contains(to, "@") {
		return fmt.Errorf("invalid email address: %s", to)
	}

	p.SentEmails = append(p.SentEmails, EmailRecord{
		To:      to,
		Subject: subject,
		Body:    body,
	})

	log.Printf("[MOCK] email would be sent to %s with subject: %s", to, subject)
	return nil
}

func NewSMTPEmailProvider(host string, port int, username, password, fromAddr string) *SMTPEmailProvider {
	return &SMTPEmailProvider{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		FromAddr: fromAddr,
	}
}

func NewMockEmailProvider() *MockEmailProvider {
	return &MockEmailProvider{SentEmails: make([]EmailRecord, 0)}
}

var emailProvider EmailProvider = NewMockEmailProvider()

func SetEmailProvider(provider EmailProvider) {
	if provider != nil {
		emailProvider = provider
	}
}

func GetEmailProvider() EmailProvider {
	return emailProvider
}

func SendAppointmentReminderEmail(patientEmail, patientName, doctorName, appointmentTime string) error {
	if patientEmail == "" {
		return fmt.Errorf("patient email is required")
	}

	subject := "Appointment Reminder - Your Upcoming Visit"
	body := fmt.Sprintf(`Dear %s,

This is a friendly reminder about your upcoming appointment:

Doctor: %s
Date/Time: %s

Please make sure to arrive 10-15 minutes early to check in.

If you need to reschedule or cancel, please contact us as soon as possible.

Best regards,
Healthcare Provider Team`,
		patientName,
		doctorName,
		appointmentTime,
	)

	return emailProvider.SendEmail(patientEmail, subject, body)
}

func SendPaymentConfirmationEmail(patientEmail, patientName, transactionCode string, amount float64) error {
	if patientEmail == "" {
		return fmt.Errorf("patient email is required")
	}

	subject := "Payment Confirmation"
	body := fmt.Sprintf(`Dear %s,

Your payment has been successfully processed.

Transaction Code: %s
Amount: %.2f
Time: %s

Please keep this confirmation for your records.

Best regards,
Healthcare Provider Team`,
		patientName,
		transactionCode,
		amount,
		time.Now().Format(time.RFC3339),
	)

	return emailProvider.SendEmail(patientEmail, subject, body)
}

func SendAppointmentConfirmationEmail(patientEmail, patientName, doctorName, appointmentTime string) error {
	if patientEmail == "" {
		return fmt.Errorf("patient email is required")
	}

	subject := "Appointment Confirmed"
	body := fmt.Sprintf(`Dear %s,

Your appointment has been confirmed.

Doctor: %s
Date/Time: %s

Please arrive 10 minutes early.

Best regards,
Healthcare Provider Team`,
		patientName,
		doctorName,
		appointmentTime,
	)

	return emailProvider.SendEmail(patientEmail, subject, body)
}
