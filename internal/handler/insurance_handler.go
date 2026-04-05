package handler

import (
	"net/http"

	"hospital/internal/models"
	"hospital/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateInsurancePlanRequest struct {
	Name               string  `json:"name" binding:"required"`
	ProviderName       string  `json:"provider_name"`
	CoveragePercentage float64 `json:"coverage_percentage" binding:"required"`
}

type UpdateInsurancePlanRequest struct {
	Name               string  `json:"name"`
	ProviderName       string  `json:"provider_name"`
	CoveragePercentage float64 `json:"coverage_percentage"`
	Status             string  `json:"status"`
}

type SetServiceCoverageRequest struct {
	ServiceID                string  `json:"service_id" binding:"required"`
	CustomCoveragePercentage float64 `json:"custom_coverage_percentage" binding:"required"`
}

func CreateInsurancePlan(c *gin.Context) {
	var req CreateInsurancePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan := models.InsurancePlan{
		ID:                 uuid.New().String(),
		Name:               req.Name,
		ProviderName:       req.ProviderName,
		CoveragePercentage: req.CoveragePercentage,
		Status:             "active",
	}

	if err := repository.CreateInsurancePlan(&plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, plan)
}

func GetInsurancePlan(c *gin.Context) {
	plan, err := repository.GetInsurancePlanByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "insurance plan not found"})
		return
	}

	c.JSON(http.StatusOK, plan)
}

func ListInsurancePlans(c *gin.Context) {
	plans, err := repository.ListInsurancePlans()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plans)
}

func ListAllInsurancePlans(c *gin.Context) {
	plans, err := repository.ListAllInsurancePlans()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plans)
}

func UpdateInsurancePlan(c *gin.Context) {
	id := c.Param("id")

	var req UpdateInsurancePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan, err := repository.GetInsurancePlanByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "insurance plan not found"})
		return
	}

	if req.Name != "" {
		plan.Name = req.Name
	}
	if req.ProviderName != "" {
		plan.ProviderName = req.ProviderName
	}
	if req.CoveragePercentage > 0 {
		plan.CoveragePercentage = req.CoveragePercentage
	}
	if req.Status != "" {
		plan.Status = req.Status
	}

	if err := repository.UpdateInsurancePlan(plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plan)
}

func DeleteInsurancePlan(c *gin.Context) {
	if err := repository.DeleteInsurancePlan(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "insurance plan deleted"})
}

func SetServiceCoverage(c *gin.Context) {
	planID := c.Param("id")

	var req SetServiceCoverageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := repository.SetInsuranceServiceCoverage(planID, req.ServiceID, req.CustomCoveragePercentage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "service coverage updated"})
}

func GetAppointmentInsuranceSnapshot(c *gin.Context) {
	snapshot, err := repository.GetAppointmentInsuranceSnapshot(c.Param("appointment_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "insurance snapshot not found"})
		return
	}

	c.JSON(http.StatusOK, snapshot)
}

type CalculateCoverageRequest struct {
	ServiceID       string  `json:"service_id" binding:"required"`
	InsurancePlanID string  `json:"insurance_plan_id" binding:"required"`
	TotalAmount     float64 `json:"total_amount" binding:"required"`
}

func CalculateInsuranceCoverage(c *gin.Context) {
	var req CalculateCoverageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload."})
		return
	}

	plan, err := repository.GetInsurancePlanByID(req.InsurancePlanID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Insurance plan not found"})
		return
	}

	coveragePercent := plan.CoveragePercentage

	customCov, err := repository.GetInsuranceServiceCoverage(req.InsurancePlanID, req.ServiceID)

	if err == nil && customCov != nil {
		coveragePercent = *customCov
	}

	insuredAmount := req.TotalAmount * (coveragePercent / 100.0)
	userPayAmount := req.TotalAmount - insuredAmount

	c.JSON(http.StatusOK, gin.H{
		"service_id":        req.ServiceID,
		"insurance_plan_id": req.InsurancePlanID,
		"total_amount":      req.TotalAmount,
		"coverage_percent":  coveragePercent,
		"insured_amount":    insuredAmount,
		"user_pay_amount":   userPayAmount,
	})
}
