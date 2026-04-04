package repository

import (
	"time"

	"hospital/internal/config"
	"hospital/internal/models"
)

func GetPoliciesByEffectiveDate(targetTime time.Time) ([]models.CancellationPolicy, error) {
	var policies []models.CancellationPolicy

	query := `
	SELECT id, policy_name, hours_before, refund_percentage, description, 
	       is_active, version, effective_from, effective_to, created_at
	FROM cancellation_policies
	WHERE is_active = true 
	  AND effective_from <= $1 
	  AND (effective_to IS NULL OR effective_to >= $1)
	ORDER BY hours_before DESC
	`

	err := config.DB.Select(&policies, query, targetTime)
	return policies, err
}
