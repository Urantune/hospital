package handler

import (
	"net/http"

	"hospital/internal/repository"

	"github.com/gin-gonic/gin"
)

type AssignServiceRequest struct {
	DoctorID  string `json:"doctor_id"`
	ServiceID string `json:"service_id"`
	ClinicID  string `json:"clinic_id"`
}

func AssignServiceToDoctor(c *gin.Context) {

	var req AssignServiceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := repository.AssignServiceToDoctor(
		req.DoctorID,
		req.ServiceID,
		req.ClinicID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "service assigned"})
}
