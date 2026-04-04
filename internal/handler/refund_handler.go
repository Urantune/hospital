package handler

import (
	"net/http"
	"time"

	"hospital/internal/repository"

	"github.com/gin-gonic/gin"
)

func PreviewRefund(c *gin.Context) {
	var request struct {
		AppointmentID string `json:"appointment_id" binding:"required"`
		CancelTime    string `json:"cancel_time" binding:"required"`
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

	cancelT, _ := time.Parse(time.RFC3339, request.CancelTime)

	hoursDiff := apt.StartTime.Sub(cancelT).Hours()

	if hoursDiff < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot cancel a past appointment."})
		return
	}

	bookingTime, _ := time.Parse(time.RFC3339, apt.CreatedAt)
	policies, err := repository.GetPoliciesByEffectiveDate(bookingTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cancellation policies."})
		return
	}

	var refundAmount float64 = 0
	appliedPolicy := "No Refund (Last Minute)"

	for _, p := range policies {
		if hoursDiff >= float64(p.HoursBefore) {
			refundAmount = apt.UserPayAmount * (p.RefundPercentage / 100.0)
			appliedPolicy = p.PolicyName
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"appointment_id":     apt.ID,
		"refund_amount":      refundAmount,
		"applied_policy":     appliedPolicy,
		"hours_before_start": hoursDiff,
	})
}
