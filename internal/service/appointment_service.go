package service

import (
	"database/sql"
	"errors"
	"hospital/internal/models"
	"hospital/internal/repository"
)

type CreateAppointmentInput struct {
	PatientID              string
	ClinicID               string
	DoctorID               string
	ServiceID              string
	SlotID                 string
	PaymentWindowExpiresAt *string
	TotalAmount            float64
	UserPayAmount          float64
	InsuredAmount          float64
}

type UpdateAppointmentStatusInput struct {
	Status string
	Reason string
}

var appointmentTransitions = map[string]map[string]bool{
	"CREATED": {
		"PENDING_PAYMENT": true,
		"CANCELLED":       true,
	},
	"PENDING_PAYMENT": {
		"CONFIRMED": true,
		"CANCELLED": true,
	},
	"CONFIRMED": {
		"IN_PROGRESS": true,
		"CANCELLED":   true,
		"NO_SHOW":     true,
	},
	"IN_PROGRESS": {
		"COMPLETED": true,
		"CANCELLED": true,
	},
	"COMPLETED": {},
	"CANCELLED": {},
	"NO_SHOW":   {},
}

func CreateAppointment(input CreateAppointmentInput, changedBy string) (*models.Appointment, error) {
	if input.PatientID == "" || input.ClinicID == "" || input.DoctorID == "" || input.ServiceID == "" || input.SlotID == "" {
		return nil, errors.New("patient_id, clinic_id, doctor_id, service_id and slot_id are required")
	}

	appointment := &models.Appointment{
		PatientID:              input.PatientID,
		ClinicID:               input.ClinicID,
		DoctorID:               input.DoctorID,
		ServiceID:              input.ServiceID,
		SlotID:                 input.SlotID,
		Status:                 "CREATED",
		PaymentWindowExpiresAt: input.PaymentWindowExpiresAt,
		TotalAmount:            input.TotalAmount,
		UserPayAmount:          input.UserPayAmount,
		InsuredAmount:          input.InsuredAmount,
	}

	if err := repository.CreateAppointment(appointment); err != nil {
		return nil, err
	}

	if err := repository.CreateAppointmentStateHistory(&models.AppointmentStateHistory{
		AppointmentID: appointment.ID,
		FromState:     nil,
		ToState:       appointment.Status,
		ChangedBy:     changedBy,
		Reason:        nil,
	}); err != nil {
		return nil, err
	}

	return appointment, nil
}

func GetAppointment(id string) (*models.Appointment, error) {
	return repository.GetAppointmentByID(id)
}

func ListMyAppointments(patientID string) ([]models.Appointment, error) {
	return repository.ListAppointmentsByPatient(patientID)
}

func ListAppointmentHistory(appointmentID string) ([]models.AppointmentStateHistory, error) {
	return repository.ListAppointmentStateHistory(appointmentID)
}

func UpdateAppointmentStatus(appointmentID string, input UpdateAppointmentStatusInput, changedBy string) (*models.Appointment, error) {
	if input.Status == "" {
		return nil, errors.New("status is required")
	}

	appointment, err := repository.GetAppointmentByID(appointmentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("appointment not found")
		}
		return nil, err
	}

	if !canTransitionAppointment(appointment.Status, input.Status) {
		return nil, errors.New("invalid appointment status transition")
	}

	if err := repository.UpdateAppointmentStatus(appointmentID, input.Status); err != nil {
		return nil, err
	}

	var fromState *string
	if appointment.Status != "" {
		state := appointment.Status
		fromState = &state
	}

	var reason *string
	if input.Reason != "" {
		r := input.Reason
		reason = &r
	}

	if err := repository.CreateAppointmentStateHistory(&models.AppointmentStateHistory{
		AppointmentID: appointmentID,
		FromState:     fromState,
		ToState:       input.Status,
		ChangedBy:     changedBy,
		Reason:        reason,
	}); err != nil {
		return nil, err
	}

	return repository.GetAppointmentByID(appointmentID)
}

func canTransitionAppointment(fromState, toState string) bool {
	if fromState == toState {
		return true
	}

	allowedNext, ok := appointmentTransitions[fromState]
	if !ok {
		return false
	}

	return allowedNext[toState]
}
