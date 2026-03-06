package models

import (
	"database/sql"
	"time"
)

type RefreshToken struct {
	ID           string       `db:"id"`
	UserID       string       `db:"user_id"`
	RefreshToken string       `db:"refresh_token"`
	DeviceInfo   string       `db:"device_info"`
	IpAddress    string       `db:"ip_address"`
	ExpiresAt    time.Time    `db:"expires_at"`
	RevokedAt    sql.NullTime `db:"revoked_at"`
	CreatedAt    time.Time    `db:"created_at"`
}
