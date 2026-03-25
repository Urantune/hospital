package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"hospital/internal/middleware"
	"hospital/internal/service"

	"github.com/gin-gonic/gin"
)

func SyncCMSChange(c *gin.Context) {
	user := middleware.CurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user context not found"})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}

	var envelope service.CMSChangeEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON payload"})
		return
	}

	signature := c.GetHeader("X-CMS-Signature")
	result, err := service.ProcessCMSChange(envelope, body, user.ID, c.ClientIP(), signature)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "CMS change processed safely",
		"event_id": result.EventID,
		"status":   result.Status,
	})
}
