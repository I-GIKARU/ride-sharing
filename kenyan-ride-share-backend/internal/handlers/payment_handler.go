package handlers

import (
	"net/http"

	"kenyan-ride-share-backend/internal/models"
	"kenyan-ride-share-backend/internal/services"
	"kenyan-ride-share-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentHandler struct {
	db           *gorm.DB
	mpesaService *services.MpesaService
}

func NewPaymentHandler(db *gorm.DB) *PaymentHandler {
	return &PaymentHandler{
		db:           db,
		mpesaService: services.NewMpesaService(db),
	}
}

type MpesaSTKPushRequest struct {
	RideID      string  `json:"ride_id" binding:"required"`
	PhoneNumber string  `json:"phone_number" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
}

func (h *PaymentHandler) InitiateMpesaPayment(c *gin.Context) {
	currentUserID := c.GetString("user_id")

	var req MpesaSTKPushRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate and format phone number
	if !utils.ValidateKenyanPhoneNumber(req.PhoneNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Kenyan phone number"})
		return
	}
	req.PhoneNumber = utils.FormatKenyanPhoneNumber(req.PhoneNumber)

	// Parse ride ID
	rideUUID, err := uuid.Parse(req.RideID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ride ID"})
		return
	}

	// Verify ride exists and user is the passenger
	var ride models.Ride
	if err := h.db.Preload("RideRequest").Where("id = ? AND passenger_id = ? AND status = ?", rideUUID, currentUserID, "completed").First(&ride).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ride not found or not authorized"})
		return
	}

	// Check if payment already exists
	var existingPayment models.Payment
	if err := h.db.Where("ride_id = ?", rideUUID).First(&existingPayment).Error; err == nil {
		if existingPayment.PaymentStatus == "completed" {
			c.JSON(http.StatusConflict, gin.H{"error": "Payment already completed"})
			return
		}
	}

	// Initiate M-Pesa STK Push
	accountReference := "RIDE-" + ride.ID.String()[:8]
	stkResponse, err := h.mpesaService.InitiateSTKPush(req.PhoneNumber, req.Amount, accountReference)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate M-Pesa payment: " + err.Error()})
		return
	}

	// Create or update payment record
	payment := models.Payment{
		RideID:        rideUUID,
		Amount:        req.Amount,
		Currency:      "KES",
		PaymentMethod: "mpesa",
		TransactionID: &stkResponse.CheckoutRequestID,
		PaymentStatus: "pending",
	}

	if existingPayment.ID != uuid.Nil {
		// Update existing payment
		if err := h.db.Model(&existingPayment).Updates(payment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment record"})
			return
		}
		payment.ID = existingPayment.ID
	} else {
		// Create new payment
		if err := h.db.Create(&payment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment record"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":            "M-Pesa payment initiated",
		"payment_id":         payment.ID,
		"checkout_request_id": stkResponse.CheckoutRequestID,
		"customer_message":   stkResponse.CustomerMessage,
	})
}

func (h *PaymentHandler) MpesaCallback(c *gin.Context) {
	var callbackData map[string]interface{}
	if err := c.ShouldBindJSON(&callbackData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.mpesaService.ProcessCallback(callbackData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process callback: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Callback processed successfully"})
}

func (h *PaymentHandler) GetPayment(c *gin.Context) {
	paymentID := c.Param("id")

	var payment models.Payment
	if err := h.db.Preload("Ride").Where("id = ?", paymentID).First(&payment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	c.JSON(http.StatusOK, payment)
}

