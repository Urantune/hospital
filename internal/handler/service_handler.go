package handler

import (
	"net/http"
	"time"

	"hospital/internal/models"
	"hospital/internal/repository"

	"github.com/gin-gonic/gin"
)

func CreateService(c *gin.Context) {

	var service models.MedicalService

	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service.Status = "active"

	err := repository.CreateService(&service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, service)
}

func GetServices(c *gin.Context) {

	services, err := repository.GetAllServices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, services)
}

func PreviewPrice(c *gin.Context) {
	var request struct {
		ServiceID string `json:"service_id" binding:"required"`
		StartTime string `json:"start_time" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: Service ID and Start Time are required."})
		return
	}

	service, err := repository.GetServiceByID(request.ServiceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Medical service not found in system."})
		return
	}

	t, err := time.Parse(time.RFC3339, request.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time format. Please use ISO 8601 (RFC3339)."})
		return
	}

	var surcharge float64 = 0
	policyName := "Standard Pricing"

	if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
		surcharge = service.BasePrice * 0.20
		policyName = "Weekend Surcharge (20%)"
	} else if t.Hour() >= 17 || t.Hour() < 8 {
		surcharge = service.BasePrice * 0.10
		policyName = "After-Hours Surcharge (10%)"
	}

	c.JSON(http.StatusOK, gin.H{
		"service_id":     service.ID,
		"service_name":   service.Name,
		"base_price":     service.BasePrice,
		"surcharge":      surcharge,
		"total_price":    service.BasePrice + surcharge,
		"applied_policy": policyName,
	})
}
