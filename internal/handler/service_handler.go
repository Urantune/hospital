package handler

import (
	"net/http"

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
