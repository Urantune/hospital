package repository

import "hospital/internal/config"

func AssignServiceToDoctor(doctorID, serviceID, clinicID string) error {

	query := `
	INSERT INTO doctor_service_mapping
	(doctor_id, service_id, clinic_id)
	VALUES ($1,$2,$3)
	`

	_, err := config.DB.Exec(query, doctorID, serviceID, clinicID)
	return err
}
