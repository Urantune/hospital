package service

import (
	"database/sql"
	"errors"
	"hospital/internal/models"
	"hospital/internal/repository"

	"github.com/google/uuid"
)

type UpsertMedicalConfigurationInput struct {
	Category  string
	ConfigKey string
	ConfigVal string
	Status    string
	ClinicID  string
}

func ListMedicalConfigurations(requestingUser *models.User) ([]models.MedicalConfiguration, error) {
	return repository.ListMedicalConfigurations(nullableStringValue(requestingUser.ClinicID), requestingUser.Role == RoleSystemAdmin)
}

func CreateMedicalConfiguration(input UpsertMedicalConfigurationInput, actingUser *models.User, ipAddress string) (*models.MedicalConfiguration, error) {
	if input.Category == "" || input.ConfigKey == "" || input.ConfigVal == "" {
		return nil, errors.New("category, config_key and config_val are required")
	}

	cfg := &models.MedicalConfiguration{
		ID:        uuid.New().String(),
		Category:  input.Category,
		ConfigKey: input.ConfigKey,
		ConfigVal: input.ConfigVal,
		Status:    input.Status,
		CreatedBy: actingUser.ID,
		UpdatedBy: actingUser.ID,
	}

	if cfg.Status == "" {
		cfg.Status = "active"
	}

	clinicID := input.ClinicID
	if clinicID == "" && actingUser.ClinicID.Valid {
		clinicID = actingUser.ClinicID.String
	}
	if clinicID != "" {
		cfg.ClinicID = sql.NullString{String: clinicID, Valid: true}
	}

	if err := repository.CreateMedicalConfiguration(cfg); err != nil {
		return nil, err
	}

	_ = repository.CreateAuditLog(&models.AuditLog{
		UserID:      actingUser.ID,
		ClinicID:    nullableStringValue(cfg.ClinicID),
		Action:      "config.create",
		Resource:    "medical_configurations",
		ResourceID:  cfg.ID,
		Description: "Medical configuration created",
		IPAddress:   ipAddress,
	})

	return cfg, nil
}

func UpdateMedicalConfiguration(id string, input UpsertMedicalConfigurationInput, actingUser *models.User, ipAddress string) (*models.MedicalConfiguration, error) {
	if input.Category == "" || input.ConfigKey == "" || input.ConfigVal == "" {
		return nil, errors.New("category, config_key and config_val are required")
	}

	cfg := &models.MedicalConfiguration{
		ID:        id,
		Category:  input.Category,
		ConfigKey: input.ConfigKey,
		ConfigVal: input.ConfigVal,
		Status:    input.Status,
		UpdatedBy: actingUser.ID,
	}

	if cfg.Status == "" {
		cfg.Status = "active"
	}
	if input.ClinicID != "" {
		cfg.ClinicID = sql.NullString{String: input.ClinicID, Valid: true}
	} else if actingUser.ClinicID.Valid {
		cfg.ClinicID = sql.NullString{String: actingUser.ClinicID.String, Valid: true}
	}

	if err := repository.UpdateMedicalConfiguration(cfg); err != nil {
		return nil, err
	}

	_ = repository.CreateAuditLog(&models.AuditLog{
		UserID:      actingUser.ID,
		ClinicID:    nullableStringValue(cfg.ClinicID),
		Action:      "config.update",
		Resource:    "medical_configurations",
		ResourceID:  cfg.ID,
		Description: "Medical configuration updated",
		IPAddress:   ipAddress,
	})

	return cfg, nil
}
