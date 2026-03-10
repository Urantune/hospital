package handler

import (
	"hospital/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetRoles(c *gin.Context) {
	roles, err := repository.GetSystemRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
}

func AssignRole(c *gin.Context) {
	var input struct {
		UserID string `json:"user_id" binding:"required"`
		RoleID int    `json:"role_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := repository.AssignUserRole(input.UserID, input.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_ = repository.RevokeAllUserTokens(input.UserID)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
