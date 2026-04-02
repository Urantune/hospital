package models

type SystemStats struct {
	TotalNoShows     int `db:"total_no_shows" json:"total_no_shows"`
	PaymentFailures  int `db:"payment_failures" json:"payment_failures"`
	PendingReminders int `db:"pending_reminders" json:"pending_reminders"`
}

type AppointmentMetrics struct {
	TotalAppointments     int     `json:"total_appointments"`
	CompletedAppointments int     `json:"completed_appointments"`
	CancelledAppointments int     `json:"cancelled_appointments"`
	NoShowAppointments    int     `json:"no_show_appointments"`
	PendingAppointments   int     `json:"pending_appointments"`
	RescheduleCount       int     `json:"reschedule_count"`
	CompletionRate        float64 `json:"completion_rate"`
	NoShowRate            float64 `json:"no_show_rate"`
	CancellationRate      float64 `json:"cancellation_rate"`
}

type PaymentMetrics struct {
	TotalTransactions      int     `json:"total_transactions"`
	SuccessfulTransactions int     `json:"successful_transactions"`
	FailedTransactions     int     `json:"failed_transactions"`
	PendingTransactions    int     `json:"pending_transactions"`
	TotalAmount            float64 `json:"total_amount"`
	SuccessfulAmount       float64 `json:"successful_amount"`
	AverageTransactionSize float64 `json:"average_transaction_size"`
	SuccessRate            float64 `json:"success_rate"`
	FailureRate            float64 `json:"failure_rate"`
}

type ReminderMetrics struct {
	TotalReminders   int     `json:"total_reminders"`
	SentReminders    int     `json:"sent_reminders"`
	FailedReminders  int     `json:"failed_reminders"`
	PendingReminders int     `json:"pending_reminders"`
	SuccessRate      float64 `json:"success_rate"`
}

type ClinicMetrics struct {
	TotalClinics      int     `json:"total_clinics"`
	ActiveClinics     int     `json:"active_clinics"`
	InactiveClinics   int     `json:"inactive_clinics"`
	TotalDoctors      int     `json:"total_doctors"`
	AvgBookingsPerDay float64 `json:"avg_bookings_per_day"`
}

type UserMetrics struct {
	TotalUsers        int `json:"total_users"`
	ActiveUsers       int `json:"active_users"`
	SuspendedUsers    int `json:"suspended_users"`
	NewUsersThisMonth int `json:"new_users_this_month"`
	NewUsersThisWeek  int `json:"new_users_this_week"`
}

type TimePeriodMetrics struct {
	Period             string             `json:"period"`
	AppointmentMetrics AppointmentMetrics `json:"appointment_metrics"`
	PaymentMetrics     PaymentMetrics     `json:"payment_metrics"`
	ReminderMetrics    ReminderMetrics    `json:"reminder_metrics"`
	ClinicMetrics      ClinicMetrics      `json:"clinic_metrics"`
}

type SystemHealthDashboard struct {
	Timestamp          string             `json:"timestamp"`
	SystemStats        SystemStats        `json:"system_stats"`
	AppointmentMetrics AppointmentMetrics `json:"appointment_metrics"`
	PaymentMetrics     PaymentMetrics     `json:"payment_metrics"`
	ReminderMetrics    ReminderMetrics    `json:"reminder_metrics"`
	ClinicMetrics      ClinicMetrics      `json:"clinic_metrics"`
	UserMetrics        UserMetrics        `json:"user_metrics"`
	TopErrorsToday     []ErrorOccurrence  `json:"top_errors_today"`
	RecentFailures     []FailureLog       `json:"recent_failures"`
}

type ErrorOccurrence struct {
	ErrorType   string `db:"error_type" json:"error_type"`
	Message     string `db:"message" json:"message"`
	Count       int    `db:"count" json:"count"`
	LastOccured string `db:"last_occurred" json:"last_occurred"`
}

type FailureLog struct {
	ID            string `db:"id" json:"id"`
	ResourceType  string `db:"resource_type" json:"resource_type"`
	ResourceID    string `db:"resource_id" json:"resource_id"`
	FailureReason string `db:"failure_reason" json:"failure_reason"`
	Timestamp     string `db:"timestamp" json:"timestamp"`
	Details       string `db:"details" json:"details"`
}

type SystemMetricsTimeSeries struct {
	Granularity string                `json:"granularity"`
	Periods     []TimePeriodMetrics   `json:"periods"`
	Summary     SystemHealthDashboard `json:"summary"`
}
