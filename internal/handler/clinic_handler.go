package handler

import (
	"net/http"

	"hospital/internal/config"

	"github.com/gin-gonic/gin"
)

type Clinic struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Status  string `json:"status"`
}

func CreateClinic(c *gin.Context) {

	var clinic Clinic

	if err := c.ShouldBindJSON(&clinic); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
	INSERT INTO clinics (name,address,phone,status)
	VALUES ($1,$2,$3,$4)
	`

	_, err := config.DB.Exec(
		query,
		clinic.Name,
		clinic.Address,
		clinic.Phone,
		clinic.Status,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "clinic created",
	})
}
