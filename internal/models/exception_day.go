package models

import "time"

type ExceptionDay struct {
	ID        string    `db:"id" json:"id"`
	DoctorID  string    `db:"doctor_id" json:"doctor_id"`
	Date      time.Time `db:"date" json:"date"`
	Type      string    `db:"type" json:"type"`
	StartTime *string   `db:"start_time" json:"start_time"`
	EndTime   *string   `db:"end_time" json:"end_time"`
	Reason    *string   `db:"reason" json:"reason"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	CreatedBy string    `db:"created_by" json:"created_by"`
}
