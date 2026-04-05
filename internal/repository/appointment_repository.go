package repository

import (
	"errors"
	"hospital/internal/config"
	"hospital/internal/models"
)

func CreateAppointment(appointment *models.Appointment) error {
	query := `
	INSERT INTO appointments
	(patient_id, clinic_id, doctor_id, service_id, slot_id, status, payment_window_expires_at, total_amount, user_pay_amount, insured_amount, base_price_at_booking, surcharge_at_booking, total_price_at_booking, applied_policy_snapshot, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, NOW(), NOW())
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
		appointment.BasePriceAtBooking,
		appointment.SurchargeAtBooking,
		appointment.TotalPriceAtBooking,
		appointment.AppliedPolicySnapshot,
	).Scan(&appointment.ID, &appointment.CreatedAt, &appointment.UpdatedAt)
}

func GetAppointmentByID(id string) (*models.Appointment, error) {
	var appointment models.Appointment

	query := `
	SELECT id, patient_id, clinic_id, doctor_id, service_id, slot_id, status, payment_window_expires_at, total_amount, user_pay_amount, insured_amount, start_time, created_at, updated_at
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

func ListAppointmentsByDoctor(doctorID string) ([]models.Appointment, error) {
	var appointments []models.Appointment

	query := `
	SELECT id, patient_id, clinic_id, doctor_id, service_id, slot_id, status, payment_window_expires_at, total_amount, user_pay_amount, insured_amount, created_at, updated_at
	FROM appointments
	WHERE doctor_id = $1
	ORDER BY created_at DESC
	`

	err := config.DB.Select(&appointments, query, doctorID)
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

func ValidateBooking(slotID, doctorID, clinicID, serviceID, userID string) error {

	var slotStatus string
	err := config.DB.Get(&slotStatus,
		`SELECT status FROM time_slot WHERE id = $1`, slotID)
	if err != nil {
		return errors.New("slot not found")
	}
	if slotStatus != "available" {
		return errors.New("slot not available")
	}

	var lockCount int
	err = config.DB.Get(&lockCount,
		`SELECT COUNT(*) FROM slot_lock 
		 WHERE slot_id=$1 AND locked_until > now()`, slotID)

	if lockCount > 0 {
		return errors.New("slot is locked")
	}

	var doctorStatus string
	err = config.DB.Get(&doctorStatus,
		`SELECT status FROM doctors WHERE id=$1`, doctorID)

	if doctorStatus != "active" {
		return errors.New("doctor inactive")
	}

	var mapCount int
	err = config.DB.Get(&mapCount,
		`SELECT COUNT(*) FROM doctor_service_mapping 
		 WHERE doctor_id=$1 AND service_id=$2`,
		doctorID, serviceID)

	if mapCount == 0 {
		return errors.New("service not supported by doctor")
	}

	var userStatus string
	err = config.DB.Get(&userStatus,
		`SELECT status FROM users WHERE id=$1`, userID)

	if userStatus != "active" {
		return errors.New("user inactive")
	}

	var count int
	err = config.DB.Get(&count, `
	SELECT COUNT(*) FROM time_slot
	WHERE id=$1 AND doctor_id=$2 AND clinic_id=$3 AND status='available'
`, slotID, doctorID, clinicID)

	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("invalid slot")
	}

	return nil
}

func MarkSlotBooked(slotID string) error {
	_, err := config.DB.Exec(`
		UPDATE time_slot
		SET status='booked'
		WHERE id=$1
	`, slotID)
	return err
}

func InsertStateHistory(appointmentID, from, to, userID string) {
	config.DB.Exec(`
	INSERT INTO appointment_state_history
	(appointment_id, from_state, to_state, changed_by, created_at)
	VALUES ($1,$2,$3,$4,now())
	`, appointmentID, from, to, userID)
}

func ExpireAppointments() error {

	rows, err := config.DB.Query(`
	SELECT id, slot_id FROM appointments
	WHERE status='PENDING_PAYMENT'
	AND payment_window_expires_at < now()
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id, slotID string
		rows.Scan(&id, &slotID)

		tx, err := config.DB.Begin()
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			UPDATE appointments
			SET status='CANCELLED', updated_at=now()
			WHERE id=$1
		`, id)
		if err != nil {
			tx.Rollback()
			continue
		}

		_, err = tx.Exec(`
			UPDATE time_slot
			SET status='available'
			WHERE id=$1
		`, slotID)
		if err != nil {
			tx.Rollback()
			continue
		}

		_, err = tx.Exec(`
			INSERT INTO appointment_state_history
			(appointment_id, from_state, to_state, changed_by, created_at)
			VALUES ($1,$2,$3,$4,now())
		`, id, "PENDING_PAYMENT", "CANCELLED", "system")
		if err != nil {
			tx.Rollback()
			continue
		}

		tx.Commit()
	}

	return nil
}

func ListAuditLogsByResource(resource string, resourceID string) ([]models.AuditLog, error) {
	var logs []models.AuditLog

	query := `
	SELECT id, user_id, clinic_id, action, resource, resource_id, description, ip_address, created_at
	FROM audit_logs
	WHERE resource = $1 AND resource_id = $2
	ORDER BY created_at DESC
	`

	err := config.DB.Select(&logs, query, resource, resourceID)
	if err != nil {
		return nil, err
	}

	return logs, nil
}
