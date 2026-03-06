package models

import "database/sql"

type User struct {
	ID                string         `db:"id"`
	Email             string         `db:"email"`
	Phone             sql.NullString `db:"phone"`
	PasswordHash      string         `db:"password_hash"`
	IsVerified        bool           `db:"is_verified"`
	Status            string         `db:"status"`
	VerificationToken sql.NullString `db:"verification_token"`
	CreatedAt         string         `db:"created_at"`
	UpdatedAt         string         `db:"updated_at"`
}
