package models

type Doctor struct {
	ID             string `db:"id" json:"id"`
	UserID         string `db:"user_id" json:"user_id"`
	ClinicID       string `db:"clinic_id" json:"clinic_id"`
	Specialization string `db:"specialization" json:"specialization"`
	Status         string `db:"status" json:"status"`
	CreatedAt      string `db:"created_at" json:"created_at"`
}
