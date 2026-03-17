package repository

import (
	"hospital/internal/config"
	"hospital/internal/models"
)

func CreateClinic(clinic *models.Clinic) error {
	query := `
	INSERT INTO clinics
	(id, code, name, owner_user_id, status, effective_from, effective_to, created_at, updated_at)
	VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, GETDATE(), GETDATE())
	`

	_, err := config.DB.Exec(
		query,
		clinic.ID,
		clinic.Code,
		clinic.Name,
		clinic.OwnerUserID,
		clinic.Status,
		clinic.EffectiveFrom,
		clinic.EffectiveTo,
	)

	return err
}

func GetClinicByID(id string) (*models.Clinic, error) {
	var clinic models.Clinic

	query := `
	SELECT id, code, name, owner_user_id, status, effective_from, effective_to, created_at, updated_at
	FROM clinics
	WHERE id = @p1
	`

	err := config.DB.Get(&clinic, query, id)
	if err != nil {
		return nil, err
	}

	return &clinic, nil
}

func ListClinics() ([]models.Clinic, error) {
	var clinics []models.Clinic

	query := `
	SELECT id, code, name, owner_user_id, status, effective_from, effective_to, created_at, updated_at
	FROM clinics
	ORDER BY created_at DESC
	`

	err := config.DB.Select(&clinics, query)
	if err != nil {
		return nil, err
	}

	return clinics, nil
}

func ListClinicsByOwnerOrClinic(ownerUserID, clinicID string) ([]models.Clinic, error) {
	var clinics []models.Clinic

	query := `
	SELECT id, code, name, owner_user_id, status, effective_from, effective_to, created_at, updated_at
	FROM clinics
	WHERE owner_user_id = @p1 OR id = @p2
	ORDER BY created_at DESC
	`

	err := config.DB.Select(&clinics, query, ownerUserID, clinicID)
	if err != nil {
		return nil, err
	}

	return clinics, nil
}

func UpdateClinicStatus(id, status string, effectiveFrom, effectiveTo interface{}) error {
	query := `
	UPDATE clinics
	SET status = @p2,
	    effective_from = @p3,
	    effective_to = @p4,
	    updated_at = GETDATE()
	WHERE id = @p1
	`

	_, err := config.DB.Exec(query, id, status, effectiveFrom, effectiveTo)
	return err
}
