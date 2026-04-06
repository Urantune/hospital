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

	var clinic struct {
		Code          string `json:"code"`
		Name          string `json:"name"`
		Status        string `json:"status"`
		EffectiveFrom string `json:"effective_from"`
		EffectiveTo   string `json:"effective_to"`
	}

	if err := c.ShouldBindJSON(&clinic); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("userID") // 👈 lấy từ middleware

	query := `
	INSERT INTO clinics (code, name, status, effective_from, effective_to, owner_user_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := config.DB.Exec(
		query,
		clinic.Code,
		clinic.Name,
		clinic.Status,
		clinic.EffectiveFrom,
		clinic.EffectiveTo,
		userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "clinic created",
	})
}
