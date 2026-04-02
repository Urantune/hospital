package repository

import (
	"hospital/internal/config"
	"hospital/internal/models"
)

func CreateInsurancePlan(plan *models.InsurancePlan) error {
	query := `
	INSERT INTO insurance_plans
	(id, name, provider_name, coverage_percentage, status, created_at)
	VALUES ($1, $2, $3, $4, $5, NOW())
	RETURNING id, created_at
	`

	return config.DB.QueryRow(
		query,
		plan.ID,
		plan.Name,
		plan.ProviderName,
		plan.CoveragePercentage,
		plan.Status,
	).Scan(&plan.ID, &plan.CreatedAt)
}

func GetInsurancePlanByID(id string) (*models.InsurancePlan, error) {
	var plan models.InsurancePlan

	query := `
	SELECT id, name, provider_name, coverage_percentage, status, created_at
	FROM insurance_plans
	WHERE id = $1
	`

	if err := config.DB.Get(&plan, query, id); err != nil {
		return nil, err
	}

	return &plan, nil
}

func ListInsurancePlans() ([]models.InsurancePlan, error) {
	var plans []models.InsurancePlan

	query := `
	SELECT id, name, provider_name, coverage_percentage, status, created_at
	FROM insurance_plans
	WHERE status = 'active'
	ORDER BY created_at DESC
	`

	if err := config.DB.Select(&plans, query); err != nil {
		return nil, err
	}

	return plans, nil
}

func ListAllInsurancePlans() ([]models.InsurancePlan, error) {
	var plans []models.InsurancePlan

	query := `
	SELECT id, name, provider_name, coverage_percentage, status, created_at
	FROM insurance_plans
	ORDER BY created_at DESC
	`

	if err := config.DB.Select(&plans, query); err != nil {
		return nil, err
	}

	return plans, nil
}

func UpdateInsurancePlan(plan *models.InsurancePlan) error {
	query := `
	UPDATE insurance_plans
	SET name = $2,
	    provider_name = $3,
	    coverage_percentage = $4,
	    status = $5
	WHERE id = $1
	`

	_, err := config.DB.Exec(
		query,
		plan.ID,
		plan.Name,
		plan.ProviderName,
		plan.CoveragePercentage,
		plan.Status,
	)

	return err
}

func DeleteInsurancePlan(id string) error {
	_, err := config.DB.Exec(`DELETE FROM insurance_plans WHERE id = $1`, id)
	return err
}

func GetInsuranceServiceCoverage(insurancePlanID, serviceID string) (*float64, error) {
	var coverage *float64

	query := `
	SELECT custom_coverage_percentage
	FROM insurance_service_coverage
	WHERE insurance_plan_id = $1 AND service_id = $2
	`

	if err := config.DB.Get(&coverage, query, insurancePlanID, serviceID); err != nil {
		return nil, err
	}

	return coverage, nil
}

func SetInsuranceServiceCoverage(insurancePlanID, serviceID string, coverage float64) error {
	query := `
	INSERT INTO insurance_service_coverage
	(insurance_plan_id, service_id, custom_coverage_percentage)
	VALUES ($1, $2, $3)
	ON CONFLICT (insurance_plan_id, service_id)
	DO UPDATE SET custom_coverage_percentage = EXCLUDED.custom_coverage_percentage
	`

	_, err := config.DB.Exec(query, insurancePlanID, serviceID, coverage)
	return err
}

func GetAppointmentInsuranceSnapshot(appointmentID string) (*models.AppointmentInsuranceSnapshot, error) {
	return GetAppointmentInsuranceSnapshotByAppointmentID(appointmentID)
}
