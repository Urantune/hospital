package repository

import (
	"hospital/internal/config"
	"hospital/internal/models"
)

func CreateDoctor(d *models.Doctor) error {

	query := `
	INSERT INTO doctors (user_id, clinic_id, specialization, status)
	VALUES ($1,$2,$3,$4)
	RETURNING id
	`

	return config.DB.QueryRow(
		query,
		d.UserID,
		d.ClinicID,
		d.Specialization,
		d.Status,
	).Scan(&d.ID)
}

func GetDoctorByID(id string) (*models.Doctor, error) {

	var doctor models.Doctor

	query := `
	SELECT id,user_id,clinic_id,specialization,status,created_at
	FROM doctors
	WHERE id=$1
	`

	err := config.DB.Get(&doctor, query, id)
	if err != nil {
		return nil, err
	}

	return &doctor, nil
}

func UpdateDoctorStatus(id string, status string) error {

	query := `
	UPDATE doctors
	SET status=$1
	WHERE id=$2
	`

	_, err := config.DB.Exec(query, status, id)
	return err
}
