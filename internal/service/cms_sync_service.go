package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"hospital/internal/models"
	"hospital/internal/repository"
	"os"
)

type CMSChangeEnvelope struct {
	EventID    string          `json:"event_id"`
	Source     string          `json:"source"`
	EntityType string          `json:"entity_type"`
	EntityID   string          `json:"entity_id"`
	Action     string          `json:"action"`
	Payload    json.RawMessage `json:"payload"`
}

type CMSMedicalConfigurationPayload struct {
	Category  string `json:"category"`
	ConfigKey string `json:"config_key"`
	ConfigVal string `json:"config_val"`
	Status    string `json:"status"`
	ClinicID  string `json:"clinic_id"`
}

func ValidateCMSSignature(body []byte, signature string) bool {
	secret := os.Getenv("CMS_SYNC_SECRET")
	if secret == "" {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(signature))
}

func ProcessCMSChange(envelope CMSChangeEnvelope, body []byte, processedBy, ipAddress, signature string) error {
	if !ValidateCMSSignature(body, signature) {
		return errors.New("invalid cms signature")
	}

	if envelope.EventID == "" || envelope.EntityType == "" || envelope.Action == "" {
		return errors.New("event_id, entity_type and action are required")
	}

	existingEvent, err := repository.GetCMSChangeEventByEventID(envelope.EventID)
	if err == nil && existingEvent != nil {
		return nil
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	event := &models.CMSChangeEvent{
		EventID:      envelope.EventID,
		Source:       envelope.Source,
		EntityType:   envelope.EntityType,
		EntityID:     envelope.EntityID,
		Action:       envelope.Action,
		Payload:      string(body),
		Status:       "received",
		ErrorMessage: "",
		ProcessedBy:  processedBy,
	}

	if err := repository.CreateCMSChangeEvent(event); err != nil {
		return err
	}

	applyErr := applyCMSChange(envelope, processedBy, ipAddress)
	if applyErr != nil {
		_ = repository.UpdateCMSChangeEventStatus(envelope.EventID, "failed", applyErr.Error(), processedBy)
		return applyErr
	}

	_ = repository.UpdateCMSChangeEventStatus(envelope.EventID, "applied", "", processedBy)
	return nil
}

func applyCMSChange(envelope CMSChangeEnvelope, processedBy, ipAddress string) error {
	switch envelope.EntityType {
	case "medical_configuration":
		return applyMedicalConfigurationChange(envelope, processedBy, ipAddress)
	default:
		return errors.New("unsupported cms entity_type")
	}
}

func applyMedicalConfigurationChange(envelope CMSChangeEnvelope, processedBy, ipAddress string) error {
	var payload CMSMedicalConfigurationPayload
	if err := json.Unmarshal(envelope.Payload, &payload); err != nil {
		return errors.New("invalid medical_configuration payload")
	}

	actingUser := &models.User{
		ID:   processedBy,
		Role: RoleSystemAdmin,
	}

	switch envelope.Action {
	case "created":
		_, err := CreateMedicalConfiguration(UpsertMedicalConfigurationInput{
			Category:  payload.Category,
			ConfigKey: payload.ConfigKey,
			ConfigVal: payload.ConfigVal,
			Status:    payload.Status,
			ClinicID:  payload.ClinicID,
		}, actingUser, ipAddress)
		return err
	case "updated":
		_, err := UpdateMedicalConfiguration(envelope.EntityID, UpsertMedicalConfigurationInput{
			Category:  payload.Category,
			ConfigKey: payload.ConfigKey,
			ConfigVal: payload.ConfigVal,
			Status:    payload.Status,
			ClinicID:  payload.ClinicID,
		}, actingUser, ipAddress)
		return err
	default:
		return errors.New("unsupported cms action")
	}
}
