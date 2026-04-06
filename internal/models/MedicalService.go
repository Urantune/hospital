package models

type MedicalService struct {
	ID                     string  `db:"id" json:"id"`
	Name                   string  `db:"name" json:"name"`
	Description            string  `db:"description" json:"description"`
	DefaultDurationMinutes int     `db:"default_duration_minutes" json:"default_duration_minutes"`
	Status                 string  `db:"status" json:"status"`
	BasePrice              float64 `db:"base_price" json:"base_price"`
}
