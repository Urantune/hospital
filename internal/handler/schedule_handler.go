package handler

import (
	"hospital/internal/models"
	"hospital/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GenerateSlotsRequest struct {
	DoctorID    string `json:"doctor_id" binding:"required"`
	ClinicID    string `json:"clinic_id" binding:"required"`
	Date        string `json:"date" binding:"required"`
	DurationMin int    `json:"duration_min"`
}

type CheckScheduleImpactRequest struct {
	Weekday    int    `json:"weekday" binding:"required"`
	StartTime  string `json:"start_time" binding:"required"`
	EndTime    string `json:"end_time" binding:"required"`
	BreakStart string `json:"break_start"`
	BreakEnd   string `json:"break_end"`
}

type CreateExceptionDayRequest struct {
	Date      string  `json:"date" binding:"required"`
	Type      string  `json:"type" binding:"required"`
	StartTime *string `json:"start_time"`
	EndTime   *string `json:"end_time"`
	Reason    *string `json:"reason"`
}

func CreateSchedule(c *gin.Context) {
	var req models.DoctorSchedule
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := service.SetDoctorSchedule(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Schedule created"})
}

func GenerateSlots(c *gin.Context) {
	var req GenerateSlotsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.DurationMin == 0 {
		req.DurationMin = 30
	}

	err := service.GenerateSlotsForDoctor(req.DoctorID, req.ClinicID, req.Date, req.DurationMin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Slots generated successfully"})
}

func GetAvailableSlots(c *gin.Context) {
	doctorID := c.Query("doctor_id")
	date := c.Query("date")

	if doctorID == "" || date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "doctor_id and date are required"})
		return
	}

	slots, err := service.GetAvailableSlots(doctorID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, slots)
}

func CheckScheduleImpact(c *gin.Context) {
	doctorID := c.Param("id")
	var req CheckScheduleImpactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newSchedule := models.DoctorSchedule{
		Weekday:    req.Weekday,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		BreakStart: req.BreakStart,
		BreakEnd:   req.BreakEnd,
	}

	impactedCount, err := service.CheckScheduleImpact(doctorID, newSchedule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"affected_appointments": impactedCount,
	})
}

func UpdateScheduleWithEnforcement(c *gin.Context) {
	doctorID := c.Param("id")

	var req models.DoctorSchedule
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.DoctorID == uuid.Nil {
		parsedDoctorID, err := uuid.Parse(doctorID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doctor id"})
			return
		}
		req.DoctorID = parsedDoctorID
	}

	canChange, reason, err := service.EnforceScheduleChangeConstraints(doctorID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !canChange {
		c.JSON(http.StatusConflict, gin.H{
			"error":      reason,
			"error_code": "SCHEDULE_CHANGE_BLOCKED",
		})
		return
	}

	if err := service.SetDoctorSchedule(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "schedule updated successfully",
		"doctor_id": doctorID,
	})
}

func LockSlot(c *gin.Context) {
	slotID := c.Param("id")
	userID := c.GetString("user_id")

	err := service.ReserveSlot(slotID, userID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Could not lock slot: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Slot locked for 10 minutes"})
}

func CreateExceptionDay(c *gin.Context) {
	doctorID := c.Param("id")
	var req CreateExceptionDayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUser := c.GetString("user_id")
	err := service.CreateExceptionDay(doctorID, req.Date, req.Type, req.StartTime, req.EndTime, req.Reason, currentUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exception day created"})
}
