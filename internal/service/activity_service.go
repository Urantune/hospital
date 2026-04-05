package service

import (
	"hospital/internal/models"
	"hospital/internal/repository"
)

type ActivityFilter struct {
	EntityType string
	EntityID   string
	UserID     string
	Limit      int
	Offset     int
}

func GetActivityLogs(filter ActivityFilter) ([]models.AuditLog, error) {

	// default pagination
	if filter.Limit == 0 {
		filter.Limit = 20
	}

	return repository.GetAuditLogs(
		filter.EntityType,
		filter.EntityID,
		filter.UserID,
		filter.Limit,
		filter.Offset,
	)
}
