package repository

import (
	"hospital/internal/config"
	"hospital/internal/models"
)

func CreateDomainEvent(event *models.DomainEvent) error {
	query := `
	INSERT INTO domain_events
	(event_type, aggregate_type, aggregate_id, payload, status, created_at)
	VALUES ($1, $2, $3, $4::jsonb, $5, NOW())
	RETURNING id, created_at
	`

	return config.DB.QueryRow(
		query,
		event.EventType,
		event.AggregateType,
		event.AggregateID,
		event.Payload,
		event.Status,
	).Scan(&event.ID, &event.CreatedAt)
}

func ListDomainEventsByAggregateID(aggregateID string) ([]models.DomainEvent, error) {
	var items []models.DomainEvent

	query := `
	SELECT id, event_type, aggregate_type, aggregate_id, payload::text AS payload, status, created_at
	FROM domain_events
	WHERE aggregate_id = $1
	ORDER BY created_at ASC
	`

	err := config.DB.Select(&items, query, aggregateID)
	if err != nil {
		return nil, err
	}

	return items, nil
}
