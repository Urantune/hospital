package repository

import (
	"database/sql"
	"errors"
	"hospital/internal/config"
	"hospital/internal/models"
)

func CreateUser(user *models.User) error {

	query := `
	INSERT INTO users 
	(id, email, password_hash, verification_token, is_verified, status, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	`

	_, err := config.DB.Exec(
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.VerificationToken,
		user.IsVerified,
		user.Status,
	)

	return err
}

func GetUserByEmail(email string) (*models.User, error) {

	var user models.User

	query := `
	SELECT id, email, phone, password_hash, is_verified, status, verification_token, created_at, updated_at
	FROM users
	WHERE email = $1
	`

	err := config.DB.Get(&user, query, email)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}

		return nil, err
	}

	return &user, nil
}

func VerifyUser(token string) error {

	query := `
	UPDATE users
	SET is_verified = TRUE,
	    verification_token = NULL,
	    status = 'active',
	    updated_at = NOW()
	WHERE verification_token = $1
	`

	result, err := config.DB.Exec(query, token)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("invalid token")
	}

	return nil
}

func SaveRefreshToken(userID, token string) error {

	query := `
	INSERT INTO refresh_tokens 
	(user_id, refresh_token, expires_at, created_at)
	VALUES ($1, $2, NOW() + INTERVAL '7 days', NOW())
	`

	_, err := config.DB.Exec(query, userID, token)

	return err
}

func GetRefreshToken(token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	query := `
		SELECT id, user_id, refresh_token, expires_at, revoked_at, created_at 
		FROM refresh_tokens 
		WHERE refresh_token = $1
	`
	err := config.DB.Get(&rt, query, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("token not found")
		}
		return nil, err
	}
	return &rt, nil
}

func RevokeRefreshToken(token string) error {
	query := `
		UPDATE refresh_tokens 
		SET revoked_at = NOW() 
		WHERE refresh_token = $1
	`
	_, err := config.DB.Exec(query, token)
	return err
}

func RevokeAllUserTokens(userID string) error {
	query := `
		UPDATE refresh_tokens 
		SET revoked_at = NOW() 
		WHERE user_id = $1 AND revoked_at IS NULL
	`
	_, err := config.DB.Exec(query, userID)
	return err
}
func GetSystemRoles() ([]models.Role, error) {
	var roles []models.Role
	query := `SELECT id, name, description FROM roles`
	err := config.DB.Select(&roles, query)
	return roles, err
}

func AssignUserRole(userID string, roleID int) error {
	query := `
	UPDATE users 
	SET role_id = $1, updated_at = NOW() 
	WHERE id = $2
	`
	_, err := config.DB.Exec(query, roleID, userID)
	return err
}
