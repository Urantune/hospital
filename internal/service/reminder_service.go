package service

import (
	"errors"
	"hospital/internal/models"
	"hospital/internal/repository"
	"strings"
	"time"
)

type ReminderService struct{}

func SendEmailReminder(_ *models.AppointmentReminder, appointmentDetails map[string]interface{}) error {
	patientEmail, ok := appointmentDetails["patient_email"].(string)
	if !ok || patientEmail == "" {
		return errors.New("patient email not found in appointment details")
	}

	patientName, _ := appointmentDetails["patient_name"].(string)
	if patientName == "" {
		patientName = "Patient"
	}

	doctorName, _ := appointmentDetails["doctor_name"].(string)
	if doctorName == "" {
		doctorName = "your doctor"
	}

	appointmentTime, _ := appointmentDetails["appointment_time"].(string)
	if appointmentTime == "" {
		appointmentTime = "your scheduled time"
	}

	return SendAppointmentReminderEmail(patientEmail, patientName, doctorName, appointmentTime)
}

func SendSMSReminder(_ *models.AppointmentReminder, appointmentDetails map[string]interface{}) error {
	patientPhone, ok := appointmentDetails["patient_phone"].(string)
	if !ok || patientPhone == "" {
		return errors.New("patient phone not found in appointment details")
	}

	patientName, _ := appointmentDetails["patient_name"].(string)
	if patientName == "" {
		patientName = "Valued Patient"
	}

	doctorName, _ := appointmentDetails["doctor_name"].(string)
	appointmentTime, _ := appointmentDetails["appointment_time"].(string)
	if appointmentTime == "" {
		appointmentTime = "your scheduled time"
	}

	return GetSMSProvider().SendAppointmentReminderSMS(patientPhone, patientName, doctorName, appointmentTime)
}

func SendPushReminder(_ *models.AppointmentReminder, appointmentDetails map[string]interface{}) error {
	deviceToken, ok := appointmentDetails["device_token"].(string)
	if !ok || deviceToken == "" {
		return errors.New("device token not found in appointment details")
	}

	patientName, _ := appointmentDetails["patient_name"].(string)
	if patientName == "" {
		patientName = "Valued Patient"
	}

	doctorName, _ := appointmentDetails["doctor_name"].(string)
	appointmentTime, _ := appointmentDetails["appointment_time"].(string)
	if appointmentTime == "" {
		appointmentTime = "your scheduled time"
	}

	return GetPushNotificationProvider().SendAppointmentReminderNotification(deviceToken, patientName, doctorName, appointmentTime)
}

func ProcessReminder(reminder *models.AppointmentReminder, appointmentDetails map[string]interface{}) error {
	var err error

	switch strings.ToUpper(reminder.ReminderType) {
	case "EMAIL":
		err = SendEmailReminder(reminder, appointmentDetails)
	case "SMS":
		err = SendSMSReminder(reminder, appointmentDetails)
	case "PUSH":
		err = SendPushReminder(reminder, appointmentDetails)
	default:
		return errors.New("unknown reminder type: " + reminder.ReminderType)
	}

	if err != nil {
		_ = repository.MarkReminderAsFailed(reminder.ID)
		return err
	}

	return repository.MarkReminderAsSent(reminder.ID)
}

func ProcessPendingReminders() (int, int, error) {
	reminders, err := repository.ListPendingReminders()
	if err != nil {
		return 0, 0, err
	}

	processedCount := 0
	failedCount := 0

	for idx := range reminders {
		reminder := reminders[idx]

		appointmentDetails, detailErr := buildReminderDetails(reminder.AppointmentID)
		if detailErr != nil {
			_ = repository.MarkReminderAsFailed(reminder.ID)
			failedCount++
			continue
		}

		if err := ProcessReminder(&reminder, appointmentDetails); err != nil {
			failedCount++
			continue
		}

		processedCount++
	}

	return processedCount, failedCount, nil
}

func RetryFailedReminders() (int, int, error) {
	reminders, err := repository.ListFailedReminders()
	if err != nil {
		return 0, 0, err
	}

	processedCount := 0
	failedCount := 0

	for idx := range reminders {
		reminder := reminders[idx]

		appointment, err := repository.GetAppointmentByID(reminder.AppointmentID)
		if err != nil || appointment.Status == "CANCELLED" {
			_ = repository.MarkReminderAsCancelled(reminder.ID)
			continue
		}

		appointmentDetails, detailErr := buildReminderDetails(reminder.AppointmentID)
		if detailErr != nil {
			failedCount++
			continue
		}

		if err := ProcessReminder(&reminder, appointmentDetails); err != nil {
			failedCount++
			continue
		}

		processedCount++
	}

	return processedCount, failedCount, nil
}

func GetAppointmentReminderDetails(appointmentID string) (map[string]interface{}, error) {
	return buildReminderDetails(appointmentID)
}

func buildReminderDetails(appointmentID string) (map[string]interface{}, error) {
	appointment, err := repository.GetAppointmentByID(appointmentID)
	if err != nil {
		return nil, err
	}

	appointmentTime := appointment.CreatedAt
	if slot, slotErr := repository.GetSlotByID(appointment.SlotID); slotErr == nil {
		appointmentTime = slot.StartTime.Format(time.RFC3339)
	}

	details := map[string]interface{}{
		"appointment_id":   appointment.ID,
		"appointment_time": appointmentTime,
		"patient_id":       appointment.PatientID,
		"doctor_id":        appointment.DoctorID,
		"clinic_id":        appointment.ClinicID,
		"service_id":       appointment.ServiceID,
		"status":           appointment.Status,
		"total_amount":     appointment.TotalAmount,
	}

	patient, err := repository.GetUserByID(appointment.PatientID)
	if err == nil {
		details["patient_email"] = patient.Email
		details["patient_phone"] = nullableStringValue(patient.Phone)
		details["patient_name"] = displayReminderUserName(patient, "Patient")
	}

	if doctor, err := repository.GetDoctorByID(appointment.DoctorID); err == nil {
		if doctorUser, doctorErr := repository.GetUserByID(doctor.UserID); doctorErr == nil {
			details["doctor_name"] = displayReminderUserName(doctorUser, "Doctor")
		}
	}

	return details, nil
}

func displayReminderUserName(user *models.User, fallback string) string {
	if user == nil {
		return fallback
	}

	if user.FullName.Valid && strings.TrimSpace(user.FullName.String) != "" {
		return user.FullName.String
	}

	if strings.TrimSpace(user.Email) != "" {
		return user.Email
	}

	return fallback
}
