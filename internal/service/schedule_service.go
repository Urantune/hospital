package service

import (
	"errors"
	"fmt"
	"hospital/internal/models"
	"hospital/internal/repository"
	"time"

	"github.com/google/uuid"
)

func SetDoctorSchedule(s models.DoctorSchedule) error {
	return repository.CreateDoctorSchedule(s)
}

func GenerateSlotsForDoctor(doctorID string, clinicID string, dateStr string, durationMin int) error {
	if durationMin <= 0 {
		durationMin = 30
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format, use YYYY-MM-DD: %w", err)
	}

	isOff, err := repository.IsFullDayOff(doctorID, date)
	if err != nil {
		return err
	}
	if isOff {
		return nil
	}

	weekday := int(date.Weekday())
	schedule, err := repository.GetDoctorScheduleByWeekday(doctorID, weekday)
	if err != nil {
		return fmt.Errorf("failed to get schedule: %w", err)
	}
	if schedule == nil {
		return fmt.Errorf("no schedule found for doctor %s on %s", doctorID, date.Weekday().String())
	}

	if date.Before(schedule.EffectiveFrom) || (!schedule.EffectiveTo.IsZero() && date.After(schedule.EffectiveTo)) {
		return nil
	}

	if err := repository.DeleteSlotsForDate(doctorID, date); err != nil {
		return err
	}

	startTime, err := time.Parse("15:04", schedule.StartTime)
	if err != nil {
		return fmt.Errorf("invalid start_time format: %w", err)
	}
	endTime, err := time.Parse("15:04", schedule.EndTime)
	if err != nil {
		return fmt.Errorf("invalid end_time format: %w", err)
	}

	dayStart := time.Date(date.Year(), date.Month(), date.Day(), startTime.Hour(), startTime.Minute(), 0, 0, startTime.Location())
	dayEnd := time.Date(date.Year(), date.Month(), date.Day(), endTime.Hour(), endTime.Minute(), 0, 0, endTime.Location())

	var breakStart, breakEnd *time.Time
	if schedule.BreakStart != "" {
		bt, err := time.Parse("15:04", schedule.BreakStart)
		if err == nil {
			bts := time.Date(date.Year(), date.Month(), date.Day(), bt.Hour(), bt.Minute(), 0, 0, bt.Location())
			breakStart = &bts
		}
	}
	if schedule.BreakEnd != "" {
		bt, err := time.Parse("15:04", schedule.BreakEnd)
		if err == nil {
			bte := time.Date(date.Year(), date.Month(), date.Day(), bt.Hour(), bt.Minute(), 0, 0, bt.Location())
			breakEnd = &bte
		}
	}

	exceptions, err := repository.GetExceptionDaysByDoctorAndDate(doctorID, date)
	if err != nil {
		return err
	}

	currentTime := dayStart
	for currentTime.Add(time.Duration(durationMin)*time.Minute).Before(dayEnd) || currentTime.Add(time.Duration(durationMin)*time.Minute).Equal(dayEnd) {
		slotEnd := currentTime.Add(time.Duration(durationMin) * time.Minute)

		if breakStart != nil && breakEnd != nil {
			if !currentTime.Before(*breakStart) && currentTime.Before(*breakEnd) {
				currentTime = *breakEnd
				continue
			}
		}

		shouldSkip := false
		for _, exc := range exceptions {
			if exc.Type == "OFF" || exc.Type == "CLOSURE" {
				if exc.StartTime != nil && exc.EndTime != nil {
					expStart, _ := time.Parse("15:04", *exc.StartTime)
					expEnd, _ := time.Parse("15:04", *exc.EndTime)
					expStartFull := time.Date(date.Year(), date.Month(), date.Day(), expStart.Hour(), expStart.Minute(), 0, 0, expStart.Location())
					expEndFull := time.Date(date.Year(), date.Month(), date.Day(), expEnd.Hour(), expEnd.Minute(), 0, 0, expEnd.Location())

					if !currentTime.Before(expStartFull) && currentTime.Before(expEndFull) {
						shouldSkip = true
						break
					}
				}
			}
		}
		if shouldSkip {
			currentTime = currentTime.Add(time.Duration(durationMin) * time.Minute)
			continue
		}

		slot := models.TimeSlot{
			ID:        uuid.New(),
			DoctorID:  uuid.MustParse(doctorID),
			ClinicID:  uuid.MustParse(clinicID),
			StartTime: currentTime,
			EndTime:   slotEnd,
			Status:    "available",
		}

		if err := repository.CreateTimeSlot(slot); err != nil {
			return fmt.Errorf("failed to create slot: %w", err)
		}

		currentTime = slotEnd
	}

	return nil
}

func ReserveSlot(slotID, userID string) error {
	return repository.LockSlotTransaction(slotID, userID, 10*time.Minute)
}

func GetAvailableSlots(doctorID string, date string) ([]models.TimeSlot, error) {
	slots, err := repository.GetAvailableSlots(doctorID, date)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch slots: %w", err)
	}

	if slots == nil {
		return []models.TimeSlot{}, nil
	}
	return slots, nil
}

func CheckScheduleImpact(doctorID string, newSchedule models.DoctorSchedule) (int, error) {
	appointments, err := repository.ListAppointmentsByDoctor(doctorID)
	if err != nil {
		return 0, err
	}

	impactedCount := 0

	for _, appt := range appointments {
		if appt.Status != "CONFIRMED" && appt.Status != "IN_PROGRESS" {
			continue
		}

		slot, err := repository.GetSlotByID(appt.SlotID)
		if err != nil {
			continue
		}

		weekday := int(slot.StartTime.Weekday())
		if weekday != newSchedule.Weekday {
			continue
		}

		if slot.StartTime.Before(newSchedule.EffectiveFrom) || (!newSchedule.EffectiveTo.IsZero() && slot.StartTime.After(newSchedule.EffectiveTo)) {
			continue
		}

		newStart, err := time.Parse("15:04", newSchedule.StartTime)
		if err != nil {
			continue
		}
		newEnd, err := time.Parse("15:04", newSchedule.EndTime)
		if err != nil {
			continue
		}

		slotTime := slot.StartTime.Hour()*60 + slot.StartTime.Minute()
		newStartMin := newStart.Hour()*60 + newStart.Minute()
		newEndMin := newEnd.Hour()*60 + newEnd.Minute()

		if slotTime < newStartMin || slotTime >= newEndMin {
			impactedCount++
			continue
		}

		if newSchedule.BreakStart != "" && newSchedule.BreakEnd != "" {
			breakStart, err := time.Parse("15:04", newSchedule.BreakStart)
			if err == nil {
				breakEnd, err := time.Parse("15:04", newSchedule.BreakEnd)
				if err == nil {
					breakStartMin := breakStart.Hour()*60 + breakStart.Minute()
					breakEndMin := breakEnd.Hour()*60 + breakEnd.Minute()

					if slotTime >= breakStartMin && slotTime < breakEndMin {
						impactedCount++
					}
				}
			}
		}
	}

	return impactedCount, nil
}

func CreateExceptionDay(doctorID string, dateStr string, exceptionType string, startTime *string, endTime *string, reason *string, createdBy string) error {
	validTypes := map[string]bool{
		"OFF":           true,
		"HOLIDAY":       true,
		"CLOSURE":       true,
		"SPECIAL_HOURS": true,
	}
	if !validTypes[exceptionType] {
		return errors.New("invalid exception type")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format, use YYYY-MM-DD: %w", err)
	}

	exception := models.ExceptionDay{
		ID:        uuid.New().String(),
		DoctorID:  doctorID,
		Date:      date,
		Type:      exceptionType,
		StartTime: startTime,
		EndTime:   endTime,
		Reason:    reason,
		CreatedBy: createdBy,
	}

	return repository.CreateExceptionDay(exception)
}
