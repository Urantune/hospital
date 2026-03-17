package handler

import (
	"net/http"

	"hospital/internal/middleware"
	"hospital/internal/service"

	"github.com/gin-gonic/gin"
)

type UpsertMedicalConfigurationRequest struct {
	Category  string `json:"category"`
	ConfigKey string `json:"config_key"`
	ConfigVal string `json:"config_val"`
	Status    string `json:"status"`
	ClinicID  string `json:"clinic_id"`
}

func ListMedicalConfigurations(c *gin.Context) {
	user := middleware.CurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user context not found"})
		return
	}

	configs, err := service.ListMedicalConfigurations(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": configs})
}

func CreateMedicalConfiguration(c *gin.Context) {
	user := middleware.CurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user context not found"})
		return
	}

	var req UpsertMedicalConfigurationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cfg, err := service.CreateMedicalConfiguration(service.UpsertMedicalConfigurationInput{
		Category:  req.Category,
		ConfigKey: req.ConfigKey,
		ConfigVal: req.ConfigVal,
		Status:    req.Status,
		ClinicID:  req.ClinicID,
	}, user, c.ClientIP())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Medical configuration created", "data": cfg})
}

func UpdateMedicalConfiguration(c *gin.Context) {
	user := middleware.CurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user context not found"})
		return
	}

	var req UpsertMedicalConfigurationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cfg, err := service.UpdateMedicalConfiguration(c.Param("id"), service.UpsertMedicalConfigurationInput{
		Category:  req.Category,
		ConfigKey: req.ConfigKey,
		ConfigVal: req.ConfigVal,
		Status:    req.Status,
		ClinicID:  req.ClinicID,
	}, user, c.ClientIP())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Medical configuration updated", "data": cfg})
}
