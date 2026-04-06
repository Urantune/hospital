package models

import "time"

type AuditLog struct {
	ID          string    `db:"id" json:"id"`
	UserID      string    `db:"user_id" json:"user_id"`
	ClinicID    string    `db:"clinic_id" json:"clinic_id"`
	Action      string    `db:"action" json:"action"`
	Resource    string    `db:"resource" json:"resource"`
	ResourceID  string    `db:"resource_id" json:"resource_id"`
	Description string    `db:"description" json:"description"`
	IPAddress   string    `db:"ip_address" json:"ip_address"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
