package repository

import "hospital/internal/config"

func UpsertClinicFromCMS(id, name, address, phone, status string) error {
	query := `
	INSERT INTO clinics (id, name, address, phone, status)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (id)
	DO UPDATE SET
		name = EXCLUDED.name,
		address = EXCLUDED.address,
		phone = EXCLUDED.phone,
		status = EXCLUDED.status
	`

	_, err := config.DB.Exec(query, id, name, address, phone, status)
	return err
}

func MarkClinicInactive(id string) error {
	query := `
	UPDATE clinics
	SET status = 'inactive'
	WHERE id = $1
	`

	_, err := config.DB.Exec(query, id)
	return err
}

func UpsertMedicalServiceFromCMS(id, name, description string, duration int, status string) error {
	query := `
	INSERT INTO medical_services (id, name, description, default_duration_minutes, status)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (id)
	DO UPDATE SET
		name = EXCLUDED.name,
		description = EXCLUDED.description,
		default_duration_minutes = EXCLUDED.default_duration_minutes,
		status = EXCLUDED.status
	`

	_, err := config.DB.Exec(query, id, name, description, duration, status)
	return err
}

func MarkMedicalServiceInactive(id string) error {
	query := `
	UPDATE medical_services
	SET status = 'inactive'
	WHERE id = $1
	`

	_, err := config.DB.Exec(query, id)
	return err
}

func UpsertDoctorServiceMappingFromCMS(doctorID, serviceID, clinicID string) error {
	query := `
	INSERT INTO doctor_service_mapping (doctor_id, service_id, clinic_id)
	VALUES ($1, $2, $3)
	ON CONFLICT (doctor_id, service_id, clinic_id) DO NOTHING
	`

	_, err := config.DB.Exec(query, doctorID, serviceID, clinicID)
	return err
}

func DeleteDoctorServiceMappingFromCMS(doctorID, serviceID, clinicID string) error {
	query := `
	DELETE FROM doctor_service_mapping
	WHERE doctor_id = $1 AND service_id = $2 AND clinic_id = $3
	`

	_, err := config.DB.Exec(query, doctorID, serviceID, clinicID)
	return err
}
