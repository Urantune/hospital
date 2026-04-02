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

func PreviewPrice(c *gin.Context) {

	var request struct {
		ServiceID string `json:"service_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Thiếu mã dịch vụ!"})
		return
	}

	service, err := repository.GetServiceByID(request.ServiceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy dịch vụ này!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"service_id":   service.ID,
		"service_name": service.Name,
		"base_price":   service.BasePrice,
		"total_price":  service.BasePrice,
	})
}
