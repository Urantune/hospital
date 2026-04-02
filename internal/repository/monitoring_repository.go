package repository

import (
	"encoding/json"
	"hospital/internal/config"
	"hospital/internal/models"
	"time"
)

func RecordSystemFailure(resourceType, resourceID, failureReason, errorMessage string, details map[string]interface{}) error {
	detailsJSON := "{}"
	if details != nil {
		raw, err := json.Marshal(details)
		if err != nil {
			return err
		}
		detailsJSON = string(raw)
	}

	query := `
	INSERT INTO system_failures
	(resource_type, resource_id, failure_reason, error_message, details, severity, created_at)
	VALUES ($1, $2, $3, $4, $5::jsonb, 'ERROR', NOW())
	`

	_, err := config.DB.Exec(query, resourceType, resourceID, failureReason, errorMessage, detailsJSON)
	return err
}

func GetUnresolvedFailures() ([]models.FailureLog, error) {
	var failures []models.FailureLog

	query := `
	SELECT id,
	       resource_type,
	       resource_id,
	       failure_reason,
	       created_at::text AS timestamp,
	       COALESCE(error_message, details::text, '') AS details
	FROM system_failures
	WHERE resolved = FALSE
	ORDER BY created_at DESC
	LIMIT 50
	`

	if err := config.DB.Select(&failures, query); err != nil {
		return nil, err
	}

	return failures, nil
}

func GetRecentFailures(limit int) ([]models.FailureLog, error) {
	var failures []models.FailureLog

	query := `
	SELECT id,
	       resource_type,
	       resource_id,
	       failure_reason,
	       created_at::text AS timestamp,
	       COALESCE(error_message, details::text, '') AS details
	FROM system_failures
	WHERE created_at >= NOW() - INTERVAL '24 hours'
	ORDER BY created_at DESC
	LIMIT $1
	`

	if err := config.DB.Select(&failures, query, limit); err != nil {
		return nil, err
	}

	return failures, nil
}

func RecordSystemMetric(metricName string, value float64, dimensionType, dimensionValue string, dimensionDate *time.Time) error {
	var dimensionHour *int
	if dimensionDate != nil {
		hour := dimensionDate.Hour()
		dimensionHour = &hour
	}

	query := `
	INSERT INTO system_metrics
	(metric_name, metric_value, dimension_type, dimension_value, dimension_date, dimension_hour, recorded_at)
	VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`

	_, err := config.DB.Exec(query, metricName, value, dimensionType, dimensionValue, dimensionDate, dimensionHour)
	return err
}

func RecordAPIPerformance(endpoint, method string, responseTimeMs int, statusCode int, success bool, errorMessage string) error {
	query := `
	INSERT INTO api_performance_metrics
	(endpoint, method, response_time_ms, status_code, success, error_message, recorded_at)
	VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`

	_, err := config.DB.Exec(query, endpoint, method, responseTimeMs, statusCode, success, errorMessage)
	return err
}

func GetAverageAPIResponseTime(endpoint string, lastNHours int) (float64, error) {
	var avgResponseTime float64

	query := `
	SELECT COALESCE(AVG(response_time_ms), 0)
	FROM api_performance_metrics
	WHERE endpoint = $1 AND recorded_at >= NOW() - make_interval(hours => $2)
	`

	if err := config.DB.QueryRow(query, endpoint, lastNHours).Scan(&avgResponseTime); err != nil {
		return 0, err
	}

	return avgResponseTime, nil
}

func GetAPIErrorRate(endpoint string, lastNHours int) (float64, error) {
	var errorRate float64

	query := `
	SELECT COALESCE(
		AVG(CASE WHEN success = FALSE THEN 100.0 ELSE 0 END),
		0
	)
	FROM api_performance_metrics
	WHERE endpoint = $1 AND recorded_at >= NOW() - make_interval(hours => $2)
	`

	if err := config.DB.QueryRow(query, endpoint, lastNHours).Scan(&errorRate); err != nil {
		return 0, err
	}

	return errorRate, nil
}

func GetFailuresByResourceType(lastNHours int) (map[string]int, error) {
	type failureCount struct {
		ResourceType string `db:"resource_type"`
		Count        int    `db:"count"`
	}

	var results []failureCount

	query := `
	SELECT resource_type, COUNT(*) AS count
	FROM system_failures
	WHERE created_at >= NOW() - make_interval(hours => $1)
	GROUP BY resource_type
	ORDER BY count DESC
	`

	if err := config.DB.Select(&results, query, lastNHours); err != nil {
		return nil, err
	}

	failureMap := make(map[string]int, len(results))
	for _, result := range results {
		failureMap[result.ResourceType] = result.Count
	}

	return failureMap, nil
}
