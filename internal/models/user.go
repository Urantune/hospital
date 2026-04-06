package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID                string         `db:"id"`
	Email             string         `db:"email"`
	FullName          sql.NullString `db:"full_name"`
	Phone             sql.NullString `db:"phone"`
	Address           sql.NullString `db:"address"`
	DateOfBirth       sql.NullString `db:"date_of_birth"`
	PasswordHash      string         `db:"password_hash"`
	IsVerified        bool           `db:"is_verified"`
	Status            string         `db:"status"`
	Role              string         `db:"role"`
	ClinicID          sql.NullString `db:"clinic_id"`
	VerificationToken sql.NullString `db:"verification_token"`
	RoleID            sql.NullInt64  `db:"role_id"`
	CreatedAt         time.Time      `db:"created_at"`
	UpdatedAt         time.Time      `db:"updated_at"`
}

type Role struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}
