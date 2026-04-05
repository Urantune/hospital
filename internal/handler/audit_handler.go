package handler

import (
	"hospital/internal/repository"
	"hospital/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListAuditLogs(c *gin.Context) {
	limit := 50
	if rawLimit := c.Query("limit"); rawLimit != "" {
		if parsedLimit, err := strconv.Atoi(rawLimit); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	logs, err := repository.ListAuditLogs(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": logs})
}

func GetAppointmentAuditLogs(c *gin.Context) {
	appointmentID := c.Param("id")

	if appointmentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Appointment ID is required"})
		return
	}

	logs, err := repository.ListAuditLogsByResource("appointments", appointmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve appointment audit logs"})
		return
	}

	if len(logs) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": []interface{}{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": logs})
}

func ProcessRefund(c *gin.Context) {
	var request struct {
		AppointmentID string `json:"appointment_id" binding:"required"`
		Reason        string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data."})
		return
	}

	apt, err := repository.GetAppointmentByID(request.AppointmentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found."})
		return
	}

	bookingTime, _ := time.Parse(time.RFC3339, apt.CreatedAt)
	policies, _ := repository.GetPoliciesByEffectiveDate(bookingTime)
	hoursDiff := apt.StartTime.Sub(time.Now()).Hours()
	var refundAmount float64 = 0

	for _, p := range policies {
		if hoursDiff >= float64(p.HoursBefore) {
			refundAmount = apt.UserPayAmount * (p.RefundPercentage / 100.0)
			break
		}
	}

	if refundAmount <= 0 {
		repository.UpdateAppointmentStatus(request.AppointmentID, "CANCELLED")
		c.JSON(http.StatusOK, gin.H{
			"message":       "Appointment cancelled. No refund applicable due to policy.",
			"refund_amount": 0,
		})
		return
	}

	payment, err := repository.GetLatestPaymentByAppointmentID(request.AppointmentID)
	if err != nil || payment == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No successful payment found to refund."})
		return
	}

	refundStatus := "refunded"
	if refundAmount < payment.Amount {
		refundStatus = "partial_refunded"
	}

	callbackInput := service.HandlePaymentCallbackInput{
		CallbackID:      uuid.New().String(),
		PaymentID:       payment.ID,
		Gateway:         payment.Gateway,
		TransactionCode: payment.TransactionCode,
		Status:          refundStatus,
		Payload:         []byte(`{"source": "system_refund", "reason": "` + request.Reason + `"}`),
	}

	_, _, err = service.HandlePaymentCallback(callbackInput, "system_admin", c.ClientIP())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process refund", "details": err.Error()})
		return
	}

	repository.UpdateAppointmentStatus(request.AppointmentID, "CANCELLED")

	c.JSON(http.StatusOK, gin.H{
		"message":        "Refund processed successfully and appointment cancelled.",
		"appointment_id": apt.ID,
		"refund_amount":  refundAmount,
		"refund_status":  refundStatus,
	})
}
