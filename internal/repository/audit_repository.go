package repository

import (
	"hospital/internal/config"
	"hospital/internal/models"

	"github.com/google/uuid"
)

func CreateAuditLog(log *models.AuditLog) error {
	if log.ID == "" {
		log.ID = uuid.New().String()
	}

	query := `
	INSERT INTO audit_logs
	(id, user_id, clinic_id, action, resource, resource_id, description, ip_address, created_at)
	VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, GETDATE())
	`

	_, err := config.DB.Exec(
		query,
		log.ID,
		log.UserID,
		log.ClinicID,
		log.Action,
		log.Resource,
		log.ResourceID,
		log.Description,
		log.IPAddress,
	)

	return err
}

func ListAuditLogs(limit int) ([]models.AuditLog, error) {
	if limit <= 0 {
		limit = 100
	}

	var logs []models.AuditLog
	query := `
	SELECT TOP (@p1) id, user_id, clinic_id, action, resource, resource_id, description, ip_address, created_at
	FROM audit_logs
	ORDER BY created_at DESC
	`

	err := config.DB.Select(&logs, query, limit)
	if err != nil {
		return nil, err
	}

	return logs, nil
}
