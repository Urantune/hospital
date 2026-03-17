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
	WHERE event_id = @p1
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
	VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, GETDATE(), GETDATE())
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
	SET status = @p2,
	    error_message = @p3,
	    processed_by = @p4,
	    processed_at = GETDATE()
	WHERE event_id = @p1
	`

	_, err := config.DB.Exec(query, eventID, status, errorMessage, processedBy)
	return err
}
