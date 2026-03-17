package handler

import (
	"database/sql"
	"net/http"

	"hospital/internal/middleware"
	"hospital/internal/service"

	"github.com/gin-gonic/gin"
)

type UpdateProfileRequest struct {
	FullName    string `json:"full_name"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	DateOfBirth string `json:"date_of_birth"`
}

func GetProfile(c *gin.Context) {
	user := middleware.CurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user context not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":            user.ID,
		"email":         user.Email,
		"full_name":     nullStringJSON(user.FullName),
		"phone":         nullStringJSON(user.Phone),
		"address":       nullStringJSON(user.Address),
		"date_of_birth": nullStringJSON(user.DateOfBirth),
		"role":          user.Role,
		"clinic_id":     nullStringJSON(user.ClinicID),
		"status":        user.Status,
		"is_verified":   user.IsVerified,
	})
}

func UpdateProfile(c *gin.Context) {
	user := middleware.CurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user context not found"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := service.UpdateProfile(user.ID, service.UpdateProfileInput{
		FullName:    req.FullName,
		Phone:       req.Phone,
		Address:     req.Address,
		DateOfBirth: req.DateOfBirth,
	}, c.ClientIP())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated",
		"user": gin.H{
			"id":            updatedUser.ID,
			"email":         updatedUser.Email,
			"full_name":     nullStringJSON(updatedUser.FullName),
			"phone":         nullStringJSON(updatedUser.Phone),
			"address":       nullStringJSON(updatedUser.Address),
			"date_of_birth": nullStringJSON(updatedUser.DateOfBirth),
			"role":          updatedUser.Role,
			"clinic_id":     nullStringJSON(updatedUser.ClinicID),
			"status":        updatedUser.Status,
			"is_verified":   updatedUser.IsVerified,
		},
	})
}

func nullStringJSON(value sql.NullString) interface{} {
	if !value.Valid {
		return nil
	}

	return value.String
}
