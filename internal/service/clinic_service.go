package service

import (
	"database/sql"
	"errors"
	"hospital/internal/models"
	"hospital/internal/repository"

	"github.com/google/uuid"
)

type CreateClinicInput struct {
	Code          string
	Name          string
	OwnerUserID   string
	Status        string
	EffectiveFrom string
	EffectiveTo   string
}

type UpdateClinicStatusInput struct {
	Status        string
	EffectiveFrom string
	EffectiveTo   string
}

func ListClinics(requestingUser *models.User) ([]models.Clinic, error) {
	if requestingUser.Role == RoleSystemAdmin {
		return repository.ListClinics()
	}

	return repository.ListClinicsByOwnerOrClinic(requestingUser.ID, nullableStringValue(requestingUser.ClinicID))
}

func GetClinic(clinicID string) (*models.Clinic, error) {
	return repository.GetClinicByID(clinicID)
}

func CreateClinic(input CreateClinicInput, actingUser *models.User, ipAddress string) (*models.Clinic, error) {
	if input.Code == "" || input.Name == "" {
		return nil, errors.New("code and name are required")
	}

	clinic := &models.Clinic{
		ID:     uuid.New().String(),
		Code:   input.Code,
		Name:   input.Name,
		Status: input.Status,
	}

	if clinic.Status == "" {
		clinic.Status = "active"
	}
	if input.OwnerUserID != "" {
		clinic.OwnerUserID = sql.NullString{String: input.OwnerUserID, Valid: true}
	}
	if input.EffectiveFrom != "" {
		clinic.EffectiveFrom = sql.NullString{String: input.EffectiveFrom, Valid: true}
	}
	if input.EffectiveTo != "" {
		clinic.EffectiveTo = sql.NullString{String: input.EffectiveTo, Valid: true}
	}

	if err := repository.CreateClinic(clinic); err != nil {
		return nil, err
	}

	_ = repository.CreateAuditLog(&models.AuditLog{
		UserID:      actingUser.ID,
		ClinicID:    clinic.ID,
		Action:      "clinic.create",
		Resource:    "clinics",
		ResourceID:  clinic.ID,
		Description: "Clinic created",
		IPAddress:   ipAddress,
	})

	return clinic, nil
}

func UpdateClinicStatus(clinicID string, input UpdateClinicStatusInput, actingUser *models.User, ipAddress string) error {
	if input.Status == "" {
		return errors.New("status is required")
	}

	var effectiveFrom interface{}
	var effectiveTo interface{}

	if input.EffectiveFrom != "" {
		effectiveFrom = input.EffectiveFrom
	}
	if input.EffectiveTo != "" {
		effectiveTo = input.EffectiveTo
	}

	if err := repository.UpdateClinicStatus(clinicID, input.Status, effectiveFrom, effectiveTo); err != nil {
		return err
	}

	_ = repository.CreateAuditLog(&models.AuditLog{
		UserID:      actingUser.ID,
		ClinicID:    clinicID,
		Action:      "clinic.status.update",
		Resource:    "clinics",
		ResourceID:  clinicID,
		Description: "Clinic status/effectivity updated",
		IPAddress:   ipAddress,
	})

	return nil
}
