package repository

import (
	"hospital/internal/config"
	"hospital/internal/models"
)

func CreateMedicalConfiguration(cfg *models.MedicalConfiguration) error {
	query := `
	INSERT INTO medical_configurations
	(id, category, config_key, config_val, status, clinic_id, created_by, updated_by, created_at, updated_at)
	VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, GETDATE(), GETDATE())
	`

	_, err := config.DB.Exec(
		query,
		cfg.ID,
		cfg.Category,
		cfg.ConfigKey,
		cfg.ConfigVal,
		cfg.Status,
		cfg.ClinicID,
		cfg.CreatedBy,
		cfg.UpdatedBy,
	)

	return err
}

func ListMedicalConfigurations(clinicID string, isAdmin bool) ([]models.MedicalConfiguration, error) {
	var configs []models.MedicalConfiguration

	query := `
	SELECT id, category, config_key, config_val, status, clinic_id, created_by, updated_by, created_at, updated_at
	FROM medical_configurations
	`

	if isAdmin {
		err := config.DB.Select(&configs, query+" ORDER BY created_at DESC")
		if err != nil {
			return nil, err
		}
		return configs, nil
	}

	err := config.DB.Select(&configs, query+" WHERE clinic_id = @p1 ORDER BY created_at DESC", clinicID)
	if err != nil {
		return nil, err
	}

	return configs, nil
}

func UpdateMedicalConfiguration(cfg *models.MedicalConfiguration) error {
	query := `
	UPDATE medical_configurations
	SET category = @p2,
	    config_key = @p3,
	    config_val = @p4,
	    status = @p5,
	    clinic_id = @p6,
	    updated_by = @p7,
	    updated_at = GETDATE()
	WHERE id = @p1
	`

	_, err := config.DB.Exec(
		query,
		cfg.ID,
		cfg.Category,
		cfg.ConfigKey,
		cfg.ConfigVal,
		cfg.Status,
		cfg.ClinicID,
		cfg.UpdatedBy,
	)

	return err
}
