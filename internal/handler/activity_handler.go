package handler

import (
	"net/http"
	"strconv"

	"hospital/internal/service"

	"github.com/gin-gonic/gin"
)

func GetActivityLogs(c *gin.Context) {

	entityType := c.Query("entity_type")
	entityID := c.Query("entity_id")
	userID := c.Query("user_id")

	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	logs, err := service.GetActivityLogs(service.ActivityFilter{
		EntityType: entityType,
		EntityID:   entityID,
		UserID:     userID,
		Limit:      limit,
		Offset:     offset,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}
