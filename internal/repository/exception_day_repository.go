package repository

import (
	"hospital/internal/config"
	"hospital/internal/models"
	"time"
)

func CreateExceptionDay(e models.ExceptionDay) error {
	query := `INSERT INTO exception_day (id, doctor_id, date, type, start_time, end_time, reason, created_by) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := config.DB.Exec(query, e.ID, e.DoctorID, e.Date, e.Type, e.StartTime, e.EndTime, e.Reason, e.CreatedBy)
	return err
}

func GetExceptionDaysByDoctorAndDate(doctorID string, date time.Time) ([]models.ExceptionDay, error) {
	var exceptions []models.ExceptionDay
	query := `SELECT id, doctor_id, date, type, start_time, end_time, reason, created_by, created_at 
	          FROM exception_day WHERE doctor_id = $1 AND DATE(date) = $2`

	rows, err := config.DB.Query(query, doctorID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e models.ExceptionDay
		rows.Scan(&e.ID, &e.DoctorID, &e.Date, &e.Type, &e.StartTime, &e.EndTime, &e.Reason, &e.CreatedBy, &e.CreatedAt)
		exceptions = append(exceptions, e)
	}
	return exceptions, nil
}

func IsFullDayOff(doctorID string, date time.Time) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM exception_day 
	          WHERE doctor_id = $1 AND DATE(date) = $2 AND type IN ('OFF', 'HOLIDAY', 'CLOSURE') 
	          AND start_time IS NULL AND end_time IS NULL`

	err := config.DB.QueryRow(query, doctorID, date).Scan(&count)
	return count > 0, err
}

func DeleteExceptionDay(exceptionID string) error {
	query := `DELETE FROM exception_day WHERE id = $1`
	_, err := config.DB.Exec(query, exceptionID)
	return err
}
