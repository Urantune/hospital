package handler

import (
	"encoding/json"
	"net/http"

	"hospital/internal/middleware"
	"hospital/internal/service"

	"github.com/gin-gonic/gin"
)

type CreateAppointmentPaymentRequest struct {
	Amount                     float64         `json:"amount"`
	Currency                   string          `json:"currency"`
	Gateway                    string          `json:"gateway"`
	TransactionCode            string          `json:"transaction_code"`
	IdempotencyKey             string          `json:"idempotency_key"`
	BasePrice                  float64         `json:"base_price"`
	ExtraFee                   float64         `json:"extra_fee"`
	CancellationPolicySnapshot json.RawMessage `json:"cancellation_policy_snapshot"`
	PlanID                     *string         `json:"plan_id"`
	CoverageRate               float64         `json:"coverage_rate"`
	CoverageCap                float64         `json:"coverage_cap"`
}

func CreateAppointmentPayment(c *gin.Context) {
	currentUser := middleware.CurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user context not found"})
		return
	}

	var req CreateAppointmentPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := service.CreateAppointmentPayment(service.CreateAppointmentPaymentInput{
		AppointmentID:              c.Param("id"),
		Amount:                     req.Amount,
		Currency:                   req.Currency,
		Gateway:                    req.Gateway,
		TransactionCode:            req.TransactionCode,
		IdempotencyKey:             req.IdempotencyKey,
		BasePrice:                  req.BasePrice,
		ExtraFee:                   req.ExtraFee,
		CancellationPolicySnapshot: req.CancellationPolicySnapshot,
		PlanID:                     req.PlanID,
		CoverageRate:               req.CoverageRate,
		CoverageCap:                req.CoverageCap,
	}, currentUser.ID, c.ClientIP())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, payment)
}

func HandlePaymentCallback(c *gin.Context) {
	currentUser := middleware.CurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user context not found"})
		return
	}

	var req service.HandlePaymentCallbackInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, payment, err := service.HandlePaymentCallback(req, currentUser.ID, c.ClientIP())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  result,
		"payment": payment,
	})
}

func GetAppointmentFinancialSnapshot(c *gin.Context) {
	snapshot, err := service.GetAppointmentFinancialSnapshot(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, snapshot)
}

func GetPaymentHistory(c *gin.Context) {
	items, err := service.ListPaymentHistory(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func GetAppointmentDomainEvents(c *gin.Context) {
	items, err := service.ListAppointmentDomainEvents(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}
