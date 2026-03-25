package repository

import (
	"hospital/internal/config"
	"hospital/internal/models"
)

func CreateAppointment(appointment *models.Appointment) error {
	query := `
	INSERT INTO appointments
	(patient_id, clinic_id, doctor_id, service_id, slot_id, status, payment_window_expires_at, total_amount, user_pay_amount, insured_amount, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
	RETURNING id, created_at, updated_at
	`

	return config.DB.QueryRow(
		query,
		appointment.PatientID,
		appointment.ClinicID,
		appointment.DoctorID,
		appointment.ServiceID,
		appointment.SlotID,
		appointment.Status,
		appointment.PaymentWindowExpiresAt,
		appointment.TotalAmount,
		appointment.UserPayAmount,
		appointment.InsuredAmount,
	).Scan(&appointment.ID, &appointment.CreatedAt, &appointment.UpdatedAt)
}

func GetAppointmentByID(id string) (*models.Appointment, error) {
	var appointment models.Appointment

	query := `
	SELECT id, patient_id, clinic_id, doctor_id, service_id, slot_id, status, payment_window_expires_at, total_amount, user_pay_amount, insured_amount, created_at, updated_at
	FROM appointments
	WHERE id = $1
	`

	err := config.DB.Get(&appointment, query, id)
	if err != nil {
		return nil, err
	}

	return &appointment, nil
}

func ListAppointmentsByPatient(patientID string) ([]models.Appointment, error) {
	var appointments []models.Appointment

	query := `
	SELECT id, patient_id, clinic_id, doctor_id, service_id, slot_id, status, payment_window_expires_at, total_amount, user_pay_amount, insured_amount, created_at, updated_at
	FROM appointments
	WHERE patient_id = $1
	ORDER BY created_at DESC
	`

	err := config.DB.Select(&appointments, query, patientID)
	if err != nil {
		return nil, err
	}

	return appointments, nil
}

func UpdateAppointmentStatus(id, status string) error {
	query := `
	UPDATE appointments
	SET status = $1,
	    updated_at = NOW()
	WHERE id = $2
	`

	_, err := config.DB.Exec(query, status, id)
	return err
}

func CreateAppointmentStateHistory(history *models.AppointmentStateHistory) error {
	query := `
	INSERT INTO appointment_state_history
	(appointment_id, from_state, to_state, changed_by, reason, created_at)
	VALUES ($1, $2, $3, $4, $5, NOW())
	RETURNING id, created_at
	`

	return config.DB.QueryRow(
		query,
		history.AppointmentID,
		history.FromState,
		history.ToState,
		history.ChangedBy,
		history.Reason,
	).Scan(&history.ID, &history.CreatedAt)
}

func ListAppointmentStateHistory(appointmentID string) ([]models.AppointmentStateHistory, error) {
	var items []models.AppointmentStateHistory

	query := `
	SELECT id, appointment_id, from_state, to_state, changed_by, reason, created_at
	FROM appointment_state_history
	WHERE appointment_id = $1
	ORDER BY created_at ASC
	`

	err := config.DB.Select(&items, query, appointmentID)
	if err != nil {
		return nil, err
	}

	return items, nil
}
