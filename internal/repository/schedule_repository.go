package repository

import (
	"context"
	"errors"
	"hospital/internal/config"
	"hospital/internal/models"
	"time"
)

func CreateDoctorSchedule(s models.DoctorSchedule) error {
	query := `INSERT INTO doctor_schedule (doctor_id, weekday, start_time, end_time, break_start, break_end, effective_from) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := config.DB.Exec(query, s.DoctorID, s.Weekday, s.StartTime, s.EndTime, s.BreakStart, s.BreakEnd, s.EffectiveFrom)
	return err
}

func GetDoctorScheduleByWeekday(doctorID string, weekday int) (*models.DoctorSchedule, error) {
	var s models.DoctorSchedule
	query := `SELECT id, doctor_id, weekday, start_time, end_time, break_start, break_end, effective_from, effective_to 
	          FROM doctor_schedule WHERE doctor_id = $1 AND weekday = $2 LIMIT 1`

	err := config.DB.QueryRow(query, doctorID, weekday).Scan(
		&s.ID, &s.DoctorID, &s.Weekday, &s.StartTime, &s.EndTime, &s.BreakStart, &s.BreakEnd, &s.EffectiveFrom, &s.EffectiveTo)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func GetAvailableSlots(doctorID string, date string) ([]models.TimeSlot, error) {
	var slots []models.TimeSlot
	query := `SELECT id, doctor_id, clinic_id, start_time, end_time, status FROM time_slot 
	          WHERE doctor_id = $1 AND DATE(start_time) = $2 AND status = 'available'
	          AND id NOT IN (SELECT slot_id FROM slot_lock WHERE locked_until > NOW())`

	rows, err := config.DB.Query(query, doctorID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s models.TimeSlot
		rows.Scan(&s.ID, &s.DoctorID, &s.ClinicID, &s.StartTime, &s.EndTime, &s.Status)
		slots = append(slots, s)
	}
	return slots, nil
}

func CreateTimeSlot(slot models.TimeSlot) error {
	query := `INSERT INTO time_slot (id, doctor_id, clinic_id, start_time, end_time, status) 
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := config.DB.Exec(query, slot.ID, slot.DoctorID, slot.ClinicID, slot.StartTime, slot.EndTime, slot.Status)
	return err
}

func DeleteSlotsForDate(doctorID string, date time.Time) error {
	query := `DELETE FROM time_slot WHERE doctor_id = $1 AND DATE(start_time) = $2`
	_, err := config.DB.Exec(query, doctorID, date)
	return err
}

func GetSlotByID(slotID string) (*models.TimeSlot, error) {
	var s models.TimeSlot
	query := `SELECT id, doctor_id, clinic_id, start_time, end_time, status FROM time_slot WHERE id = $1`

	err := config.DB.QueryRow(query, slotID).Scan(&s.ID, &s.DoctorID, &s.ClinicID, &s.StartTime, &s.EndTime, &s.Status)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func UpdateSlotStatus(slotID string, status string) error {
	query := `UPDATE time_slot SET status = $1 WHERE id = $2`
	_, err := config.DB.Exec(query, status, slotID)
	return err
}

func GetActiveLocksBySlotID(slotID string) ([]models.SlotLock, error) {
	var locks []models.SlotLock
	query := `SELECT id, slot_id, user_id, locked_until FROM slot_lock 
	          WHERE slot_id = $1 AND locked_until > NOW()`

	rows, err := config.DB.Query(query, slotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var l models.SlotLock
		rows.Scan(&l.ID, &l.SlotID, &l.UserID, &l.LockedUntil)
		locks = append(locks, l)
	}
	return locks, nil
}

func LockSlotTransaction(slotID, userID string, duration time.Duration) error {
	tx, err := config.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var status string
	err = tx.QueryRow(`SELECT status FROM time_slot WHERE id = $1 FOR UPDATE`, slotID).Scan(&status)
	if err != nil {
		return err
	}

	if status != "available" {
		return errors.New("slot is no longer available")
	}

	lockedUntil := time.Now().Add(duration)
	_, err = tx.Exec(`INSERT INTO slot_lock (slot_id, user_id, locked_until) VALUES ($1, $2, $3)`,
		slotID, userID, lockedUntil)

	if err != nil {
		return err
	}

	return tx.Commit()
}
