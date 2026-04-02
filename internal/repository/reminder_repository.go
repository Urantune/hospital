package repository

import (
	"hospital/internal/config"
	"hospital/internal/models"
	"time"
)

func CreateReminder(reminder *models.AppointmentReminder) error {
	query := `
	INSERT INTO appointment_reminders
	(appointment_id, reminder_type, scheduled_at, status, retry_count, created_at)
	VALUES ($1, $2, $3, $4, 0, NOW())
	RETURNING id, created_at
	`

	return config.DB.QueryRow(
		query,
		reminder.AppointmentID,
		reminder.ReminderType,
		reminder.ScheduledAt,
		reminder.Status,
	).Scan(&reminder.ID, &reminder.CreatedAt)
}

func GetReminderByID(id string) (*models.AppointmentReminder, error) {
	var reminder models.AppointmentReminder

	query := `
	SELECT id, appointment_id, reminder_type, scheduled_at, sent_at, status, retry_count, created_at
	FROM appointment_reminders
	WHERE id = $1
	`

	if err := config.DB.Get(&reminder, query, id); err != nil {
		return nil, err
	}

	return &reminder, nil
}

func ListRemindersByAppointment(appointmentID string) ([]models.AppointmentReminder, error) {
	var reminders []models.AppointmentReminder

	query := `
	SELECT id, appointment_id, reminder_type, scheduled_at, sent_at, status, retry_count, created_at
	FROM appointment_reminders
	WHERE appointment_id = $1
	ORDER BY scheduled_at DESC
	`

	if err := config.DB.Select(&reminders, query, appointmentID); err != nil {
		return nil, err
	}

	return reminders, nil
}

func ListPendingReminders() ([]models.AppointmentReminder, error) {
	var reminders []models.AppointmentReminder

	query := `
	SELECT id, appointment_id, reminder_type, scheduled_at, sent_at, status, retry_count, created_at
	FROM appointment_reminders
	WHERE status = 'pending' AND scheduled_at <= NOW()
	ORDER BY scheduled_at ASC
	LIMIT 100
	`

	if err := config.DB.Select(&reminders, query); err != nil {
		return nil, err
	}

	return reminders, nil
}

func ListFailedReminders() ([]models.AppointmentReminder, error) {
	var reminders []models.AppointmentReminder

	query := `
	SELECT id, appointment_id, reminder_type, scheduled_at, sent_at, status, retry_count, created_at
	FROM appointment_reminders
	WHERE status = 'failed' AND retry_count < 3
	ORDER BY scheduled_at ASC
	LIMIT 50
	`

	if err := config.DB.Select(&reminders, query); err != nil {
		return nil, err
	}

	return reminders, nil
}

func MarkReminderAsSent(id string) error {
	_, err := config.DB.Exec(`
	UPDATE appointment_reminders
	SET status = 'sent', sent_at = NOW()
	WHERE id = $1
	`, id)
	return err
}

func MarkReminderAsFailed(id string) error {
	_, err := config.DB.Exec(`
	UPDATE appointment_reminders
	SET status = 'failed', retry_count = retry_count + 1
	WHERE id = $1
	`, id)
	return err
}

func MarkReminderAsCancelled(id string) error {
	_, err := config.DB.Exec(`
	UPDATE appointment_reminders
	SET status = 'cancelled'
	WHERE id = $1
	`, id)
	return err
}

func RetryFailedReminder(id string, newScheduledAt *time.Time) error {
	if newScheduledAt == nil {
		t := time.Now().Add(1 * time.Hour)
		newScheduledAt = &t
	}

	_, err := config.DB.Exec(`
	UPDATE appointment_reminders
	SET status = 'pending', scheduled_at = $2
	WHERE id = $1 AND status = 'failed' AND retry_count < 3
	`, id, newScheduledAt)
	return err
}

func CancelRemindersByAppointment(appointmentID string) error {
	_, err := config.DB.Exec(`
	UPDATE appointment_reminders
	SET status = 'cancelled'
	WHERE appointment_id = $1 AND status != 'sent'
	`, appointmentID)
	return err
}

func GetSystemGlobalStats() (*models.SystemStats, error) {
	var stats models.SystemStats

	query := `
	SELECT
		(SELECT COUNT(*) FROM appointments WHERE status = 'NO_SHOW') AS total_no_shows,
		(SELECT COUNT(*) FROM payments WHERE status = 'failed') AS payment_failures,
		(SELECT COUNT(*) FROM appointment_reminders WHERE status = 'pending') AS pending_reminders
	`

	if err := config.DB.Get(&stats, query); err != nil {
		return nil, err
	}

	return &stats, nil
}
