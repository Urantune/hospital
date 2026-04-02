package repository

import (
	"hospital/internal/config"
	"hospital/internal/models"
)

func CreateService(s *models.MedicalService) error {

	query := `
	INSERT INTO medical_services
	(name, description, default_duration_minutes, base_price, status)
	VALUES ($1,$2,$3,$4,$5)
	RETURNING id
	`

	return config.DB.QueryRow(
		query,
		s.Name,
		s.Description,
		s.DefaultDurationMinutes,
		s.BasePrice,
		s.Status,
	).Scan(&s.ID)
}

func GetAllServices() ([]models.MedicalService, error) {

	var services []models.MedicalService

	query := `SELECT * FROM medical_services`

	err := config.DB.Select(&services, query)
	return services, err
}
func ListAppointmentsByDoctor(doctorID string) ([]models.Appointment, error) {
	// Trả về một danh sách rỗng để nó không thắc mắc nữa
	return []models.Appointment{}, nil
}
