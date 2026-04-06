package repository

import (
	"hospital/internal/config"
	"hospital/internal/models"
)

func GetAuditLogs(entityType, entityID, userID string, limit, offset int) ([]models.AuditLog, error) {

	query := `
SELECT id, user_id, clinic_id, action, resource, resource_id, description, ip_address, created_at
FROM audit_logs
WHERE ($1 = '' OR resource = $1)
  AND ($2 = '' OR resource_id = $2)
  AND ($3 = '' OR user_id = $3)
ORDER BY created_at DESC
LIMIT $4 OFFSET $5
`

	var logs []models.AuditLog
	err := config.DB.Select(&logs, query, entityType, entityID, userID, limit, offset)
	return logs, err
}
