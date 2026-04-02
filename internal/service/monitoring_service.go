package service

import (
	"database/sql"
	"hospital/internal/config"
	"hospital/internal/models"
	"hospital/internal/repository"
	"time"
)

func GetAppointmentMetrics() (*models.AppointmentMetrics, error) {
	metrics := &models.AppointmentMetrics{}

	query := `
	SELECT
		COUNT(*) AS total,
		COALESCE(SUM(CASE WHEN status = 'COMPLETED' THEN 1 ELSE 0 END), 0) AS completed,
		COALESCE(SUM(CASE WHEN status = 'CANCELLED' THEN 1 ELSE 0 END), 0) AS cancelled,
		COALESCE(SUM(CASE WHEN status = 'NO_SHOW' THEN 1 ELSE 0 END), 0) AS no_show,
		COALESCE(SUM(CASE WHEN status IN ('CONFIRMED', 'PENDING_PAYMENT', 'IN_PROGRESS') THEN 1 ELSE 0 END), 0) AS pending
	FROM appointments
	`

	err := config.DB.QueryRow(query).Scan(
		&metrics.TotalAppointments,
		&metrics.CompletedAppointments,
		&metrics.CancelledAppointments,
		&metrics.NoShowAppointments,
		&metrics.PendingAppointments,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	applyAppointmentRates(metrics)
	return metrics, nil
}

func GetPaymentMetrics() (*models.PaymentMetrics, error) {
	metrics := &models.PaymentMetrics{}

	query := `
	SELECT
		COUNT(*) AS total,
		COALESCE(SUM(CASE WHEN status = 'success' THEN 1 ELSE 0 END), 0) AS successful,
		COALESCE(SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END), 0) AS failed,
		COALESCE(SUM(CASE WHEN status = 'initiated' THEN 1 ELSE 0 END), 0) AS pending,
		COALESCE(SUM(amount), 0) AS total_amount,
		COALESCE(SUM(CASE WHEN status = 'success' THEN amount ELSE 0 END), 0) AS successful_amount
	FROM payments
	`

	err := config.DB.QueryRow(query).Scan(
		&metrics.TotalTransactions,
		&metrics.SuccessfulTransactions,
		&metrics.FailedTransactions,
		&metrics.PendingTransactions,
		&metrics.TotalAmount,
		&metrics.SuccessfulAmount,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	applyPaymentRates(metrics)
	return metrics, nil
}

func GetReminderMetrics() (*models.ReminderMetrics, error) {
	metrics := &models.ReminderMetrics{}

	query := `
	SELECT
		COUNT(*) AS total,
		COALESCE(SUM(CASE WHEN status = 'sent' THEN 1 ELSE 0 END), 0) AS sent,
		COALESCE(SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END), 0) AS failed,
		COALESCE(SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END), 0) AS pending
	FROM appointment_reminders
	`

	err := config.DB.QueryRow(query).Scan(
		&metrics.TotalReminders,
		&metrics.SentReminders,
		&metrics.FailedReminders,
		&metrics.PendingReminders,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if metrics.TotalReminders > 0 {
		metrics.SuccessRate = float64(metrics.SentReminders) / float64(metrics.TotalReminders) * 100
	}

	return metrics, nil
}

func GetClinicMetrics() (*models.ClinicMetrics, error) {
	metrics := &models.ClinicMetrics{}

	query := `
	SELECT
		COUNT(*) AS total_clinics,
		COALESCE(SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END), 0) AS active_clinics,
		COALESCE(SUM(CASE WHEN status = 'inactive' THEN 1 ELSE 0 END), 0) AS inactive_clinics
	FROM clinics
	`

	err := config.DB.QueryRow(query).Scan(
		&metrics.TotalClinics,
		&metrics.ActiveClinics,
		&metrics.InactiveClinics,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err := config.DB.QueryRow(`SELECT COALESCE(COUNT(*), 0) FROM doctors WHERE status = 'active'`).Scan(&metrics.TotalDoctors); err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	avgQuery := `
	SELECT COALESCE(AVG(daily_count), 0) FROM (
		SELECT schedule_date, COUNT(*) AS daily_count
		FROM (
			SELECT COALESCE(DATE(ts.start_time), DATE(a.created_at)) AS schedule_date
			FROM appointments a
			LEFT JOIN time_slot ts ON ts.id = a.slot_id
		) appointment_days
		WHERE schedule_date >= CURRENT_DATE - INTERVAL '30 days'
		GROUP BY schedule_date
	) daily_stats
	`

	if err := config.DB.QueryRow(avgQuery).Scan(&metrics.AvgBookingsPerDay); err != nil && err != sql.ErrNoRows {
		metrics.AvgBookingsPerDay = 0
	}

	return metrics, nil
}

func GetUserMetrics() (*models.UserMetrics, error) {
	metrics := &models.UserMetrics{}

	query := `
	SELECT
		COUNT(*) AS total,
		COALESCE(SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END), 0) AS active,
		COALESCE(SUM(CASE WHEN status = 'suspended' THEN 1 ELSE 0 END), 0) AS suspended,
		COALESCE(SUM(CASE WHEN created_at >= NOW() - INTERVAL '1 month' THEN 1 ELSE 0 END), 0) AS new_this_month,
		COALESCE(SUM(CASE WHEN created_at >= NOW() - INTERVAL '1 week' THEN 1 ELSE 0 END), 0) AS new_this_week
	FROM users
	`

	err := config.DB.QueryRow(query).Scan(
		&metrics.TotalUsers,
		&metrics.ActiveUsers,
		&metrics.SuspendedUsers,
		&metrics.NewUsersThisMonth,
		&metrics.NewUsersThisWeek,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return metrics, nil
}

func GetSystemHealthDashboard() (*models.SystemHealthDashboard, error) {
	dashboard := &models.SystemHealthDashboard{
		Timestamp: time.Now().Format(time.RFC3339),
	}

	stats, err := repository.GetSystemGlobalStats()
	if err == nil && stats != nil {
		dashboard.SystemStats = *stats
	}

	if appointmentMetrics, err := GetAppointmentMetrics(); err == nil && appointmentMetrics != nil {
		dashboard.AppointmentMetrics = *appointmentMetrics
	}
	if paymentMetrics, err := GetPaymentMetrics(); err == nil && paymentMetrics != nil {
		dashboard.PaymentMetrics = *paymentMetrics
	}
	if reminderMetrics, err := GetReminderMetrics(); err == nil && reminderMetrics != nil {
		dashboard.ReminderMetrics = *reminderMetrics
	}
	if clinicMetrics, err := GetClinicMetrics(); err == nil && clinicMetrics != nil {
		dashboard.ClinicMetrics = *clinicMetrics
	}
	if userMetrics, err := GetUserMetrics(); err == nil && userMetrics != nil {
		dashboard.UserMetrics = *userMetrics
	}
	if topErrors, err := GetTopErrorsFromAuditLogs(); err == nil {
		dashboard.TopErrorsToday = topErrors
	}
	if recentFailures, err := GetRecentFailuresFromAuditLogs(); err == nil {
		dashboard.RecentFailures = recentFailures
	}

	return dashboard, nil
}

func GetTopErrorsFromAuditLogs() ([]models.ErrorOccurrence, error) {
	var errors []models.ErrorOccurrence

	query := `
	SELECT
		COALESCE(NULLIF(description, ''), action) AS error_type,
		COALESCE(NULLIF(description, ''), action) AS message,
		COUNT(*) AS count,
		MAX(created_at)::text AS last_occurred
	FROM audit_logs
	WHERE (
		LOWER(action) LIKE '%error%' OR
		LOWER(action) LIKE '%failed%' OR
		LOWER(COALESCE(description, '')) LIKE '%error%' OR
		LOWER(COALESCE(description, '')) LIKE '%failed%'
	) AND DATE(created_at) = CURRENT_DATE
	GROUP BY 1, 2
	ORDER BY count DESC
	LIMIT 10
	`

	if err := config.DB.Select(&errors, query); err != nil {
		return []models.ErrorOccurrence{}, err
	}

	return errors, nil
}

func GetRecentFailuresFromAuditLogs() ([]models.FailureLog, error) {
	var failures []models.FailureLog

	query := `
	SELECT
		id,
		resource AS resource_type,
		resource_id,
		COALESCE(description, action) AS failure_reason,
		created_at::text AS timestamp,
		COALESCE(description, '') AS details
	FROM audit_logs
	WHERE
		LOWER(action) LIKE '%failed%' OR
		LOWER(action) LIKE '%error%' OR
		LOWER(COALESCE(description, '')) LIKE '%failed%' OR
		LOWER(COALESCE(description, '')) LIKE '%error%'
	ORDER BY created_at DESC
	LIMIT 20
	`

	if err := config.DB.Select(&failures, query); err != nil {
		return []models.FailureLog{}, err
	}

	return failures, nil
}

func GetMetricsForDateRange(startDate, endDate time.Time) (*models.SystemHealthDashboard, error) {
	dashboard := &models.SystemHealthDashboard{
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if appointmentMetrics, err := GetAppointmentMetricsForPeriod(startDate, endDate); err == nil && appointmentMetrics != nil {
		dashboard.AppointmentMetrics = *appointmentMetrics
	}
	if paymentMetrics, err := GetPaymentMetricsForPeriod(startDate, endDate); err == nil && paymentMetrics != nil {
		dashboard.PaymentMetrics = *paymentMetrics
	}
	if reminderMetrics, err := GetReminderMetricsForPeriod(startDate, endDate); err == nil && reminderMetrics != nil {
		dashboard.ReminderMetrics = *reminderMetrics
	}
	if clinicMetrics, err := GetClinicMetrics(); err == nil && clinicMetrics != nil {
		dashboard.ClinicMetrics = *clinicMetrics
	}
	if userMetrics, err := GetUserMetrics(); err == nil && userMetrics != nil {
		dashboard.UserMetrics = *userMetrics
	}
	if topErrors, err := GetTopErrorsFromAuditLogsForPeriod(startDate, endDate); err == nil {
		dashboard.TopErrorsToday = topErrors
	}
	if recentFailures, err := GetRecentFailuresFromAuditLogsForPeriod(startDate, endDate); err == nil {
		dashboard.RecentFailures = recentFailures
	}

	return dashboard, nil
}

func GetPaymentMetricsForPeriod(startDate, endDate time.Time) (*models.PaymentMetrics, error) {
	metrics := &models.PaymentMetrics{}

	query := `
	SELECT
		COUNT(*) AS total,
		COALESCE(SUM(CASE WHEN status = 'success' THEN 1 ELSE 0 END), 0) AS successful,
		COALESCE(SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END), 0) AS failed,
		COALESCE(SUM(CASE WHEN status = 'initiated' THEN 1 ELSE 0 END), 0) AS pending,
		COALESCE(SUM(amount), 0) AS total_amount,
		COALESCE(SUM(CASE WHEN status = 'success' THEN amount ELSE 0 END), 0) AS successful_amount
	FROM payments
	WHERE created_at >= $1 AND created_at <= $2
	`

	err := config.DB.QueryRow(query, startDate, endDate).Scan(
		&metrics.TotalTransactions,
		&metrics.SuccessfulTransactions,
		&metrics.FailedTransactions,
		&metrics.PendingTransactions,
		&metrics.TotalAmount,
		&metrics.SuccessfulAmount,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	applyPaymentRates(metrics)
	return metrics, nil
}

func GetReminderMetricsForPeriod(startDate, endDate time.Time) (*models.ReminderMetrics, error) {
	metrics := &models.ReminderMetrics{}

	query := `
	SELECT
		COUNT(*) AS total,
		COALESCE(SUM(CASE WHEN status = 'sent' THEN 1 ELSE 0 END), 0) AS sent,
		COALESCE(SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END), 0) AS failed,
		COALESCE(SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END), 0) AS pending
	FROM appointment_reminders
	WHERE created_at >= $1 AND created_at <= $2
	`

	err := config.DB.QueryRow(query, startDate, endDate).Scan(
		&metrics.TotalReminders,
		&metrics.SentReminders,
		&metrics.FailedReminders,
		&metrics.PendingReminders,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if metrics.TotalReminders > 0 {
		metrics.SuccessRate = float64(metrics.SentReminders) / float64(metrics.TotalReminders) * 100
	}

	return metrics, nil
}

func GetTopErrorsFromAuditLogsForPeriod(startDate, endDate time.Time) ([]models.ErrorOccurrence, error) {
	var errors []models.ErrorOccurrence

	query := `
	SELECT
		COALESCE(NULLIF(description, ''), action) AS error_type,
		COALESCE(NULLIF(description, ''), action) AS message,
		COUNT(*) AS count,
		MAX(created_at)::text AS last_occurred
	FROM audit_logs
	WHERE (
		LOWER(action) LIKE '%error%' OR
		LOWER(action) LIKE '%failed%' OR
		LOWER(COALESCE(description, '')) LIKE '%error%' OR
		LOWER(COALESCE(description, '')) LIKE '%failed%'
	) AND created_at >= $1 AND created_at <= $2
	GROUP BY 1, 2
	ORDER BY count DESC
	LIMIT 10
	`

	if err := config.DB.Select(&errors, query, startDate, endDate); err != nil {
		return []models.ErrorOccurrence{}, err
	}

	return errors, nil
}

func GetRecentFailuresFromAuditLogsForPeriod(startDate, endDate time.Time) ([]models.FailureLog, error) {
	var failures []models.FailureLog

	query := `
	SELECT
		id,
		resource AS resource_type,
		resource_id,
		COALESCE(description, action) AS failure_reason,
		created_at::text AS timestamp,
		COALESCE(description, '') AS details
	FROM audit_logs
	WHERE (
		LOWER(action) LIKE '%failed%' OR
		LOWER(action) LIKE '%error%' OR
		LOWER(COALESCE(description, '')) LIKE '%failed%' OR
		LOWER(COALESCE(description, '')) LIKE '%error%'
	) AND created_at >= $1 AND created_at <= $2
	ORDER BY created_at DESC
	LIMIT 20
	`

	if err := config.DB.Select(&failures, query, startDate, endDate); err != nil {
		return []models.FailureLog{}, err
	}

	return failures, nil
}

func GetAppointmentMetricsForPeriod(startDate, endDate time.Time) (*models.AppointmentMetrics, error) {
	metrics := &models.AppointmentMetrics{}

	query := `
	SELECT
		COUNT(*) AS total,
		COALESCE(SUM(CASE WHEN a.status = 'COMPLETED' THEN 1 ELSE 0 END), 0) AS completed,
		COALESCE(SUM(CASE WHEN a.status = 'CANCELLED' THEN 1 ELSE 0 END), 0) AS cancelled,
		COALESCE(SUM(CASE WHEN a.status = 'NO_SHOW' THEN 1 ELSE 0 END), 0) AS no_show,
		COALESCE(SUM(CASE WHEN a.status IN ('CONFIRMED', 'PENDING_PAYMENT', 'IN_PROGRESS') THEN 1 ELSE 0 END), 0) AS pending
	FROM appointments a
	LEFT JOIN time_slot ts ON ts.id = a.slot_id
	WHERE COALESCE(ts.start_time, a.created_at) >= $1
	  AND COALESCE(ts.start_time, a.created_at) <= $2
	`

	err := config.DB.QueryRow(query, startDate, endDate).Scan(
		&metrics.TotalAppointments,
		&metrics.CompletedAppointments,
		&metrics.CancelledAppointments,
		&metrics.NoShowAppointments,
		&metrics.PendingAppointments,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	applyAppointmentRates(metrics)
	return metrics, nil
}

func applyAppointmentRates(metrics *models.AppointmentMetrics) {
	if metrics == nil || metrics.TotalAppointments == 0 {
		return
	}

	metrics.CompletionRate = float64(metrics.CompletedAppointments) / float64(metrics.TotalAppointments) * 100
	metrics.NoShowRate = float64(metrics.NoShowAppointments) / float64(metrics.TotalAppointments) * 100
	metrics.CancellationRate = float64(metrics.CancelledAppointments) / float64(metrics.TotalAppointments) * 100
}

func applyPaymentRates(metrics *models.PaymentMetrics) {
	if metrics == nil || metrics.TotalTransactions == 0 {
		return
	}

	metrics.SuccessRate = float64(metrics.SuccessfulTransactions) / float64(metrics.TotalTransactions) * 100
	metrics.FailureRate = float64(metrics.FailedTransactions) / float64(metrics.TotalTransactions) * 100
	metrics.AverageTransactionSize = metrics.TotalAmount / float64(metrics.TotalTransactions)
}
