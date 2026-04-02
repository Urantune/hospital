package handler

import (
	"net/http"

	"hospital/internal/repository"
	"hospital/internal/service"

	"github.com/gin-gonic/gin"
)

type CreateReminderRequest struct {
	AppointmentID string `json:"appointment_id" binding:"required"`
	ReminderType  string `json:"reminder_type" binding:"required"`
	ScheduledAt   string `json:"scheduled_at" binding:"required"`
}

func ListReminders(c *gin.Context) {
	reminders, err := repository.ListRemindersByAppointment(c.Param("appointment_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reminders)
}

func GetReminder(c *gin.Context) {
	reminder, err := repository.GetReminderByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "reminder not found"})
		return
	}

	c.JSON(http.StatusOK, reminder)
}

func ManuallyTriggerReminder(c *gin.Context) {
	id := c.Param("id")

	reminder, err := repository.GetReminderByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "reminder not found"})
		return
	}

	details, err := service.GetAppointmentReminderDetails(reminder.AppointmentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "appointment not found"})
		return
	}

	if err := service.ProcessReminder(reminder, details); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send reminder: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reminder sent successfully"})
}

func CancelReminder(c *gin.Context) {
	if err := repository.MarkReminderAsCancelled(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reminder cancelled"})
}

func RetryReminder(c *gin.Context) {
	id := c.Param("id")

	reminder, err := repository.GetReminderByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "reminder not found"})
		return
	}

	if reminder.Status != "failed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only failed reminders can be retried"})
		return
	}

	if err := repository.RetryFailedReminder(id, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reminder queued for retry"})
}

func ProcessPendingReminders(c *gin.Context) {
	processedCount, failedCount, err := service.ProcessPendingReminders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"processed": processedCount,
		"failed":    failedCount,
		"message":   "pending reminders processed",
	})
}

func ProcessFailedReminders(c *gin.Context) {
	processedCount, failedCount, err := service.RetryFailedReminders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"retried": processedCount,
		"failed":  failedCount,
		"message": "failed reminders processed",
	})
}

func CancelAppointmentReminders(c *gin.Context) {
	if err := repository.CancelRemindersByAppointment(c.Param("appointment_id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "all reminders for appointment cancelled"})
}

func GetReminderStats(c *gin.Context) {
	pendingReminders, err := repository.ListPendingReminders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	failedReminders, err := repository.ListFailedReminders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pending_count": len(pendingReminders),
		"failed_count":  len(failedReminders),
	})
}
