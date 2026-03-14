package handler

import (
	"net/http"

	"hospital/internal/models"
	"hospital/internal/repository"

	"github.com/gin-gonic/gin"
)

func CreateDoctor(c *gin.Context) {

	var doctor models.Doctor

	if err := c.ShouldBindJSON(&doctor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doctor.Status = "active"

	err := repository.CreateDoctor(&doctor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, doctor)
}

func DeactivateDoctor(c *gin.Context) {

	id := c.Param("id")

	err := repository.UpdateDoctorStatus(id, "inactive")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "doctor deactivated"})
}

func ActivateDoctor(c *gin.Context) {

	id := c.Param("id")

	err := repository.UpdateDoctorStatus(id, "active")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "doctor activated"})
}
