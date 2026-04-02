package models

type AppointmentReminder struct {
	ID            string  `db:"id" json:"id"`
	AppointmentID string  `db:"appointment_id" json:"appointment_id"`
	ReminderType  string  `db:"reminder_type" json:"reminder_type"`
	ScheduledAt   string  `db:"scheduled_at" json:"scheduled_at"`
	SentAt        *string `db:"sent_at" json:"sent_at"`
	Status        string  `db:"status" json:"status"`
	RetryCount    int     `db:"retry_count" json:"retry_count"`
	CreatedAt     string  `db:"created_at" json:"created_at"`
}
