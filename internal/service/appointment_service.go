package service

import (
	"database/sql"
	"errors"
	"hospital/internal/models"
	"hospital/internal/repository"
	"time"

	"github.com/google/uuid"
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
	if err := repository.ValidateBooking(
		input.SlotID,
		input.DoctorID,
		input.ClinicID,
		input.ServiceID,
		input.PatientID,
	); err != nil {
		return nil, err
	}

	expireTime := time.Now().Add(15 * time.Minute)

	appointment := &models.Appointment{
		ID:                     uuid.New().String(),
		PatientID:              input.PatientID,
		ClinicID:               input.ClinicID,
		DoctorID:               input.DoctorID,
		ServiceID:              input.ServiceID,
		SlotID:                 input.SlotID,
		Status:                 "PENDING_PAYMENT",
		PaymentWindowExpiresAt: &expireTime,
		TotalAmount:            input.TotalAmount,
		UserPayAmount:          input.UserPayAmount,
		InsuredAmount:          input.InsuredAmount,
	}

	if err := repository.CreateAppointment(appointment); err != nil {
		return nil, err
	}

	if err := repository.MarkSlotBooked(input.SlotID); err != nil {
		return nil, err
	}

	from := ""
	to := "PENDING_PAYMENT"

	repository.InsertStateHistory(appointment.ID, from, to, changedBy)

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
	return transitionAppointmentStatus(appointmentID, input.Status, input.Reason, changedBy)
}

func TransitionAppointmentStatus(appointmentID, status, reason, changedBy string) (*models.Appointment, error) {
	return transitionAppointmentStatus(appointmentID, status, reason, changedBy)
}

func transitionAppointmentStatus(appointmentID, status, reasonText, changedBy string) (*models.Appointment, error) {
	if status == "" {
		return nil, errors.New("status is required")
	}

	appointment, err := repository.GetAppointmentByID(appointmentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("appointment not found")
		}
		return nil, err
	}

	if !canTransitionAppointment(appointment.Status, status) {
		return nil, errors.New("invalid appointment status transition")
	}

	if err := repository.UpdateAppointmentStatus(appointmentID, status); err != nil {
		return nil, err
	}

	var fromState *string
	if appointment.Status != "" {
		state := appointment.Status
		fromState = &state
	}

	var reason *string
	if reasonText != "" {
		r := reasonText
		reason = &r
	}

	if err := repository.CreateAppointmentStateHistory(&models.AppointmentStateHistory{
		AppointmentID: appointmentID,
		FromState:     fromState,
		ToState:       status,
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

func StartExpireJob() {

	go func() {
		for {
			repository.ExpireAppointments()
			time.Sleep(30 * time.Second)
		}
	}()
}
