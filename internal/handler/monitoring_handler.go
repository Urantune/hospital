package handler

import (
	"net/http"
	"strconv"
	"time"

	"hospital/internal/repository"
	"hospital/internal/service"

	"github.com/gin-gonic/gin"
)

func GetSystemStats(c *gin.Context) {
	stats, err := repository.GetSystemGlobalStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func GetSystemHealthDashboard(c *gin.Context) {
	dashboard, err := service.GetSystemHealthDashboard()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dashboard)
}

func GetAppointmentMetrics(c *gin.Context) {
	metrics, err := service.GetAppointmentMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

func GetPaymentMetrics(c *gin.Context) {
	metrics, err := service.GetPaymentMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

func GetReminderMetrics(c *gin.Context) {
	metrics, err := service.GetReminderMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

func GetClinicMetrics(c *gin.Context) {
	metrics, err := service.GetClinicMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

func GetUserMetrics(c *gin.Context) {
	metrics, err := service.GetUserMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

func GetMetricsForDateRange(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date query parameters are required"})
		return
	}

	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format (use RFC3339)"})
		return
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format (use RFC3339)"})
		return
	}

	metrics, err := service.GetMetricsForDateRange(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

func GetUnresolvedFailures(c *gin.Context) {
	failures, err := repository.GetUnresolvedFailures()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"unresolved_failures": failures,
		"count":               len(failures),
	})
}

func GetRecentFailures(c *gin.Context) {
	limit := 50
	if rawLimit := c.Query("limit"); rawLimit != "" {
		if parsedLimit, err := strconv.Atoi(rawLimit); err == nil && parsedLimit > 0 && parsedLimit < 1000 {
			limit = parsedLimit
		}
	}

	failures, err := repository.GetRecentFailures(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recent_failures": failures,
		"count":           len(failures),
	})
}

func GetFailuresByResourceType(c *gin.Context) {
	hours := 24
	if rawHours := c.Query("hours"); rawHours != "" {
		if parsedHours, err := strconv.Atoi(rawHours); err == nil && parsedHours > 0 {
			hours = parsedHours
		}
	}

	failureMap, err := repository.GetFailuresByResourceType(hours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"period_hours":     hours,
		"failures_by_type": failureMap,
	})
}

func GetAPIPerformanceMetrics(c *gin.Context) {
	endpoint := c.Query("endpoint")
	if endpoint == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "endpoint query parameter is required"})
		return
	}

	hours := 1
	if rawHours := c.Query("hours"); rawHours != "" {
		if parsedHours, err := strconv.Atoi(rawHours); err == nil && parsedHours > 0 {
			hours = parsedHours
		}
	}

	avgResponseTime, err := repository.GetAverageAPIResponseTime(endpoint, hours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	errorRate, err := repository.GetAPIErrorRate(endpoint, hours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"endpoint":              endpoint,
		"period_hours":          hours,
		"avg_response_time_ms":  avgResponseTime,
		"error_rate_percentage": errorRate,
	})
}
