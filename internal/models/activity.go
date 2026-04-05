package models

import "time"

type auditLog struct {
	ID          string    `db:"id" json:"id"`
	EntityType  string    `db:"entity_type" json:"entity_type"`
	EntityID    string    `db:"entity_id" json:"entity_id"`
	Action      string    `db:"action" json:"action"`
	PerformedBy *string   `db:"performed_by" json:"performed_by"`
	OldValue    *string   `db:"old_value" json:"old_value"`
	NewValue    *string   `db:"new_value" json:"new_value"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
