package models

import "database/sql"

type Clinic struct {
	ID            string         `db:"id" json:"id"`
	Code          string         `db:"code" json:"code"`
	Name          string         `db:"name" json:"name"`
	OwnerUserID   sql.NullString `db:"owner_user_id" json:"owner_user_id"`
	Status        string         `db:"status" json:"status"`
	EffectiveFrom sql.NullString `db:"effective_from" json:"effective_from"`
	EffectiveTo   sql.NullString `db:"effective_to" json:"effective_to"`
	CreatedAt     string         `db:"created_at" json:"created_at"`
	UpdatedAt     string         `db:"updated_at" json:"updated_at"`
}
