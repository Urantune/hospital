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
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
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
	SELECT id, user_id, clinic_id, action, resource, resource_id, description, ip_address, created_at
	FROM audit_logs
	ORDER BY created_at DESC
	LIMIT $1
	`

	err := config.DB.Select(&logs, query, limit)
	if err != nil {
		return nil, err
	}

	return logs, nil
}
