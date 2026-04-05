package service

import (
	"database/sql"
	"errors"
	"hospital/internal/repository"
)

func CalculateCoverage(serviceID, planID string, totalAmount float64) (float64, float64, error) {
	plan, err := repository.GetInsurancePlanByID(planID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, errors.New("insurance plan not found")
		}
		return 0, 0, err
	}

	coveragePercent := plan.CoveragePercentage

	customCov, err := repository.GetInsuranceServiceCoverage(planID, serviceID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, 0, err
	}

	if customCov != nil {
		coveragePercent = *customCov
	}

	insuredAmount := totalAmount * (coveragePercent / 100.0)
	userPayAmount := totalAmount - insuredAmount

	return insuredAmount, userPayAmount, nil
}
