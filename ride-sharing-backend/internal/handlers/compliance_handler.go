package handlers

import (
	"net/http"
	"strconv"

	"kenyan-ride-share-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ComplianceHandler struct {
	db                *gorm.DB
	complianceService *services.ComplianceService
}

func NewComplianceHandler(db *gorm.DB) *ComplianceHandler {
	return &ComplianceHandler{
		db:                db,
		complianceService: services.NewComplianceService(db),
	}
}

func (h *ComplianceHandler) CheckDriverCompliance(c *gin.Context) {
	driverID := c.Param("id")
	currentUserID := c.GetString("user_id")
	currentUserType := c.GetString("user_type")

	// Only drivers can check their own compliance or admins can check any driver
	if currentUserType != "admin" && driverID != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	driverUUID, err := uuid.Parse(driverID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid driver ID"})
		return
	}

	status, err := h.complianceService.ValidateDriverCompliance(driverUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check compliance"})
		return
	}

	c.JSON(http.StatusOK, status)
}

func (h *ComplianceHandler) CalculateCommission(c *gin.Context) {
	fareStr := c.Query("fare")
	if fareStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fare amount required"})
		return
	}

	fare, err := strconv.ParseFloat(fareStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fare amount"})
		return
	}

	breakdown := h.complianceService.CalculateCommission(fare)
	c.JSON(http.StatusOK, breakdown)
}

func (h *ComplianceHandler) GenerateNTSAReport(c *gin.Context) {
	currentUserType := c.GetString("user_type")
	if currentUserType != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start date and end date required"})
		return
	}

	report, err := h.complianceService.GenerateNTSAReport(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate report"})
		return
	}

	c.JSON(http.StatusOK, report)
}

func (h *ComplianceHandler) ValidateVehicle(c *gin.Context) {
	var req struct {
		VehicleYear  int    `json:"vehicle_year" binding:"required"`
		VehicleMake  string `json:"vehicle_make" binding:"required"`
		VehicleModel string `json:"vehicle_model" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := h.complianceService.ValidateVehicleEligibility(req.VehicleYear, req.VehicleMake, req.VehicleModel)
	c.JSON(http.StatusOK, status)
}

