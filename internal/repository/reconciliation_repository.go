package repository

import (
	"hospital/internal/config"
	"hospital/internal/models"
)

func GetAllTodayPayments() ([]models.Payment, error) {
	var payments []models.Payment

	err := config.DB.Select(&payments, `
		SELECT * FROM payments
		WHERE DATE(created_at) = CURRENT_DATE
	`)

	return payments, err
}

func InsertReconciliationReport(gateway string, mismatchCount int) error {
	_, err := config.DB.Exec(`
		INSERT INTO reconciliation_reports
		(id, gateway, report_date, status, mismatch_count, created_at)
		VALUES (gen_random_uuid(), $1, CURRENT_DATE, $2, $3, now())
	`,
		gateway,
		"done",
		mismatchCount,
	)

	return err
}
