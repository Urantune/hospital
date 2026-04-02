package handler

import (
	"net/http"

	"hospital/internal/middleware"
	"hospital/internal/service"

	"github.com/gin-gonic/gin"
)

type CreateAppointmentRequest struct {
	PatientID              string  `json:"patient_id"`
	ClinicID               string  `json:"clinic_id"`
	DoctorID               string  `json:"doctor_id"`
	ServiceID              string  `json:"service_id"`
	SlotID                 string  `json:"slot_id"`
	PaymentWindowExpiresAt *string `json:"payment_window_expires_at"`
	TotalAmount            float64 `json:"total_amount"`
	UserPayAmount          float64 `json:"user_pay_amount"`
	InsuredAmount          float64 `json:"insured_amount"`
}

type UpdateAppointmentStatusRequest struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

func CreateAppointment(c *gin.Context) {
	currentUser := middleware.CurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user context not found"})
		return
	}

	var req CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	patientID := req.PatientID
	if patientID == "" {
		patientID = currentUser.ID
	}

	appointment, err := service.CreateAppointment(service.CreateAppointmentInput{
		PatientID:              patientID,
		ClinicID:               req.ClinicID,
		DoctorID:               req.DoctorID,
		ServiceID:              req.ServiceID,
		SlotID:                 req.SlotID,
		PaymentWindowExpiresAt: req.PaymentWindowExpiresAt,
		TotalAmount:            req.TotalAmount,
		UserPayAmount:          req.UserPayAmount,
		InsuredAmount:          req.InsuredAmount,
	}, currentUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, appointment)
}

func GetAppointment(c *gin.Context) {
	appointment, err := service.GetAppointment(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointment)
}

func GetMyAppointments(c *gin.Context) {
	currentUser := middleware.CurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user context not found"})
		return
	}

	appointments, err := service.ListMyAppointments(currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointments)
}

func GetAppointmentHistory(c *gin.Context) {
	history, err := service.ListAppointmentHistory(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

func UpdateAppointmentStatus(c *gin.Context) {
	currentUser := middleware.CurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user context not found"})
		return
	}

	var req UpdateAppointmentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appointment, err := service.UpdateAppointmentStatus(c.Param("id"), service.UpdateAppointmentStatusInput{
		Status: req.Status,
		Reason: req.Reason,
	}, currentUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointment)

}
