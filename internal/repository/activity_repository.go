package repository

import (
	"hospital/internal/config"
	"hospital/internal/models"
)

func GetAuditLogs(entityType, entityID, userID string, limit, offset int) ([]models.AuditLog, error) {

	query := `
	SELECT id, entity_type, entity_id, action, performed_by, old_value, new_value, created_at
	FROM audit_logs
	WHERE ($1 = '' OR entity_type = $1)
	  AND ($2 = '' OR entity_id = $2)
	  AND ($3 = '' OR performed_by = $3)
	ORDER BY created_at DESC
	LIMIT $4 OFFSET $5
	`

	var logs []models.AuditLog
	err := config.DB.Select(&logs, query, entityType, entityID, userID, limit, offset)
	return logs, err
}
