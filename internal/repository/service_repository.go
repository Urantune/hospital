package repository

import (
	"hospital/internal/config"
	"hospital/internal/models"
)

func CreateService(s *models.MedicalService) error {

	query := `
	INSERT INTO medical_services
	(name, description, default_duration_minutes, status)
	VALUES ($1,$2,$3,$4)
	RETURNING id
	`

	return config.DB.QueryRow(
		query,
		s.Name,
		s.Description,
		s.DefaultDurationMinutes,
		s.Status,
	).Scan(&s.ID)
}

func GetAllServices() ([]models.MedicalService, error) {

	var services []models.MedicalService

	query := `SELECT * FROM medical_services`

	err := config.DB.Select(&services, query)
	return services, err
}
