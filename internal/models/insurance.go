package models

type InsurancePlan struct {
	ID                 string  `db:"id" json:"id"`
	Name               string  `db:"name" json:"name"`
	ProviderName       string  `db:"provider_name" json:"provider_name"`
	CoveragePercentage float64 `db:"coverage_percentage" json:"coverage_percentage"`
	Status             string  `db:"status" json:"status"`
	CreatedAt          string  `db:"created_at" json:"created_at"`
}
