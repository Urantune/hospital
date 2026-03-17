package service

import (
	"database/sql"
	"errors"
	"hospital/internal/models"
	"hospital/internal/repository"
)

type UpdateProfileInput struct {
	FullName    string
	Phone       string
	Address     string
	DateOfBirth string
}

func GetProfile(userID string) (*models.User, error) {
	return repository.GetUserByID(userID)
}

func UpdateProfile(userID string, input UpdateProfileInput, ipAddress string) (*models.User, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	if input.FullName != "" {
		user.FullName = sql.NullString{String: input.FullName, Valid: true}
	}
	if input.Phone != "" {
		user.Phone = sql.NullString{String: input.Phone, Valid: true}
	}
	if input.Address != "" {
		user.Address = sql.NullString{String: input.Address, Valid: true}
	}
	if input.DateOfBirth != "" {
		user.DateOfBirth = sql.NullString{String: input.DateOfBirth, Valid: true}
	}

	if !user.IsVerified {
		return nil, errors.New("account not verified")
	}

	if err := repository.UpdateUserProfile(user); err != nil {
		return nil, err
	}

	_ = repository.CreateAuditLog(&models.AuditLog{
		UserID:      user.ID,
		ClinicID:    nullableStringValue(user.ClinicID),
		Action:      "profile.update",
		Resource:    "users",
		ResourceID:  user.ID,
		Description: "User profile updated",
		IPAddress:   ipAddress,
	})

	return user, nil
}
