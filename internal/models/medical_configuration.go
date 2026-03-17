package models

import "database/sql"

type MedicalConfiguration struct {
	ID        string         `db:"id" json:"id"`
	Category  string         `db:"category" json:"category"`
	ConfigKey string         `db:"config_key" json:"config_key"`
	ConfigVal string         `db:"config_val" json:"config_val"`
	Status    string         `db:"status" json:"status"`
	ClinicID  sql.NullString `db:"clinic_id" json:"clinic_id"`
	CreatedBy string         `db:"created_by" json:"created_by"`
	UpdatedBy string         `db:"updated_by" json:"updated_by"`
	CreatedAt string         `db:"created_at" json:"created_at"`
	UpdatedAt string         `db:"updated_at" json:"updated_at"`
}
