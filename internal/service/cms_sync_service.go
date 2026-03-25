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
	"strings"
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

type CMSClinicPayload struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Status  string `json:"status"`
}

type CMSMedicalServicePayload struct {
	ID                     string `json:"id"`
	Name                   string `json:"name"`
	Description            string `json:"description"`
	DefaultDurationMinutes int    `json:"default_duration_minutes"`
	Status                 string `json:"status"`
}

type CMSDoctorServiceMappingPayload struct {
	DoctorID  string `json:"doctor_id"`
	ServiceID string `json:"service_id"`
	ClinicID  string `json:"clinic_id"`
}

type CMSProcessResult struct {
	EventID string `json:"event_id"`
	Status  string `json:"status"`
}

func ValidateCMSSignature(body []byte, signature string) bool {
	secret := os.Getenv("CMS_SYNC_SECRET")
	if secret == "" {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(normalizeCMSSignature(signature)))
}

func normalizeCMSSignature(signature string) string {
	normalized := strings.TrimSpace(signature)
	normalized = strings.TrimPrefix(normalized, "sha256=")
	return strings.ToLower(normalized)
}

func ProcessCMSChange(envelope CMSChangeEnvelope, body []byte, processedBy, ipAddress, signature string) (*CMSProcessResult, error) {
	if !ValidateCMSSignature(body, signature) {
		return nil, errors.New("invalid cms signature")
	}

	if envelope.EventID == "" || envelope.EntityType == "" || envelope.Action == "" {
		return nil, errors.New("event_id, entity_type and action are required")
	}

	if !isAllowedCMSSource(envelope.Source) {
		return nil, errors.New("unsupported cms source")
	}

	if !isAllowedCMSAction(envelope.Action) {
		return nil, errors.New("unsupported cms action")
	}

	existingEvent, err := repository.GetCMSChangeEventByEventID(envelope.EventID)
	if err == nil && existingEvent != nil {
		return &CMSProcessResult{EventID: envelope.EventID, Status: "duplicate"}, nil
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
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
		return nil, err
	}

	applyErr := applyCMSChange(envelope, processedBy, ipAddress)
	if applyErr != nil {
		_ = repository.UpdateCMSChangeEventStatus(envelope.EventID, "failed", applyErr.Error(), processedBy)
		return nil, applyErr
	}

	_ = repository.UpdateCMSChangeEventStatus(envelope.EventID, "applied", "", processedBy)
	return &CMSProcessResult{EventID: envelope.EventID, Status: "applied"}, nil
}

func isAllowedCMSSource(source string) bool {
	switch strings.ToLower(strings.TrimSpace(source)) {
	case "cms", "directus", "admin_portal":
		return true
	default:
		return false
	}
}

func isAllowedCMSAction(action string) bool {
	switch strings.ToLower(strings.TrimSpace(action)) {
	case "created", "updated", "deleted", "activated", "deactivated":
		return true
	default:
		return false
	}
}

func applyCMSChange(envelope CMSChangeEnvelope, processedBy, ipAddress string) error {
	switch envelope.EntityType {
	case "clinic":
		return applyClinicChange(envelope)
	case "medical_service":
		return applyMedicalServiceChange(envelope)
	case "doctor_service_mapping":
		return applyDoctorServiceMappingChange(envelope)
	case "medical_configuration":
		return applyMedicalConfigurationChange(envelope, processedBy, ipAddress)
	default:
		return errors.New("unsupported cms entity_type")
	}
}

func applyClinicChange(envelope CMSChangeEnvelope) error {
	var payload CMSClinicPayload
	if err := json.Unmarshal(envelope.Payload, &payload); err != nil {
		return errors.New("invalid clinic payload")
	}

	clinicID := payload.ID
	if clinicID == "" {
		clinicID = envelope.EntityID
	}
	if clinicID == "" || payload.Name == "" {
		return errors.New("clinic id and name are required")
	}

	status := payload.Status
	if status == "" {
		status = "active"
	}

	switch strings.ToLower(envelope.Action) {
	case "created", "updated", "activated", "deactivated":
		if envelope.Action == "deactivated" {
			status = "inactive"
		}
		return repository.UpsertClinicFromCMS(clinicID, payload.Name, payload.Address, payload.Phone, status)
	case "deleted":
		return repository.MarkClinicInactive(clinicID)
	default:
		return errors.New("unsupported clinic action")
	}
}

func applyMedicalServiceChange(envelope CMSChangeEnvelope) error {
	var payload CMSMedicalServicePayload
	if err := json.Unmarshal(envelope.Payload, &payload); err != nil {
		return errors.New("invalid medical_service payload")
	}

	serviceID := payload.ID
	if serviceID == "" {
		serviceID = envelope.EntityID
	}
	if serviceID == "" || payload.Name == "" || payload.DefaultDurationMinutes <= 0 {
		return errors.New("medical_service id, name and default_duration_minutes are required")
	}

	status := payload.Status
	if status == "" {
		status = "active"
	}

	switch strings.ToLower(envelope.Action) {
	case "created", "updated", "activated", "deactivated":
		if envelope.Action == "deactivated" {
			status = "inactive"
		}
		return repository.UpsertMedicalServiceFromCMS(serviceID, payload.Name, payload.Description, payload.DefaultDurationMinutes, status)
	case "deleted":
		return repository.MarkMedicalServiceInactive(serviceID)
	default:
		return errors.New("unsupported medical_service action")
	}
}

func applyDoctorServiceMappingChange(envelope CMSChangeEnvelope) error {
	var payload CMSDoctorServiceMappingPayload
	if err := json.Unmarshal(envelope.Payload, &payload); err != nil {
		return errors.New("invalid doctor_service_mapping payload")
	}

	if payload.DoctorID == "" || payload.ServiceID == "" || payload.ClinicID == "" {
		return errors.New("doctor_id, service_id and clinic_id are required")
	}

	switch strings.ToLower(envelope.Action) {
	case "created", "updated", "activated":
		return repository.UpsertDoctorServiceMappingFromCMS(payload.DoctorID, payload.ServiceID, payload.ClinicID)
	case "deleted", "deactivated":
		return repository.DeleteDoctorServiceMappingFromCMS(payload.DoctorID, payload.ServiceID, payload.ClinicID)
	default:
		return errors.New("unsupported doctor_service_mapping action")
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
