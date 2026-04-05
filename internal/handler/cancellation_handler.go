package handler

import (
	"net/http"

	"hospital/internal/repository"

	"github.com/gin-gonic/gin"
)

func CancelAppointment(c *gin.Context) {
	var request struct {
		AppointmentID string `json:"appointment_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request. Appointment ID is required."})
		return
	}

	apt, err := repository.GetAppointmentByID(request.AppointmentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found."})
		return
	}

	if apt.Status == "CANCELLED" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This appointment is already cancelled."})
		return
	}

	err = repository.UpdateAppointmentStatus(request.AppointmentID, "CANCELLED")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel the appointment."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Appointment cancelled successfully.",
		"appointment_id": apt.ID,
		"new_status":     "CANCELLED",
	})
}
