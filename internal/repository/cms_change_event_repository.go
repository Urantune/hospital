package repository

import (
	"database/sql"
	"errors"
	"hospital/internal/config"
	"hospital/internal/models"

	"github.com/google/uuid"
)

func GetCMSChangeEventByEventID(eventID string) (*models.CMSChangeEvent, error) {
	var event models.CMSChangeEvent

	query := `
	SELECT id, event_id, source, entity_type, entity_id, action, payload, status, error_message, processed_by, processed_at, created_at
	FROM cms_change_events
	WHERE event_id = $1
	`

	err := config.DB.Get(&event, query, eventID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &event, nil
}

func CreateCMSChangeEvent(event *models.CMSChangeEvent) error {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}

	query := `
	INSERT INTO cms_change_events
	(id, event_id, source, entity_type, entity_id, action, payload, status, error_message, processed_by, processed_at, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
	`

	_, err := config.DB.Exec(
		query,
		event.ID,
		event.EventID,
		event.Source,
		event.EntityType,
		event.EntityID,
		event.Action,
		event.Payload,
		event.Status,
		event.ErrorMessage,
		event.ProcessedBy,
	)

	return err
}

func UpdateCMSChangeEventStatus(eventID, status, errorMessage, processedBy string) error {
	query := `
	UPDATE cms_change_events
	SET status = $2,
	    error_message = $3,
	    processed_by = $4,
	    processed_at = NOW()
	WHERE event_id = $1
	`

	_, err := config.DB.Exec(query, eventID, status, errorMessage, processedBy)
	return err
}
