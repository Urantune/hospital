package models

import (
	"time"

	"github.com/google/uuid"
)

type DoctorSchedule struct {
	ID            uuid.UUID `json:"id"`
	DoctorID      uuid.UUID `json:"doctor_id"`
	Weekday       int       `json:"weekday"`
	StartTime     string    `json:"start_time"`
	EndTime       string    `json:"end_time"`
	BreakStart    string    `json:"break_start"`
	BreakEnd      string    `json:"break_end"`
	EffectiveFrom time.Time `json:"effective_from"`
	EffectiveTo   time.Time `json:"effective_to"`
}

type TimeSlot struct {
	ID        uuid.UUID `json:"id"`
	DoctorID  uuid.UUID `json:"doctor_id"`
	ClinicID  uuid.UUID `json:"clinic_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
}

type SlotLock struct {
	ID          uuid.UUID `json:"id"`
	SlotID      uuid.UUID `json:"slot_id"`
	UserID      uuid.UUID `json:"user_id"`
	LockedUntil time.Time `json:"locked_until"`
}
