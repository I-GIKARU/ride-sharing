package handlers

import (
	"net/http"
	"strconv"
	"time"

	"kenyan-ride-share-backend/internal/models"
	"kenyan-ride-share-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RideHandler struct {
	db *gorm.DB
}

func NewRideHandler(db *gorm.DB) *RideHandler {
	return &RideHandler{db: db}
}

type CreateRideRequestRequest struct {
	PickupLatitude   float64 `json:"pickup_latitude" binding:"required"`
	PickupLongitude  float64 `json:"pickup_longitude" binding:"required"`
	DropoffLatitude  float64 `json:"dropoff_latitude" binding:"required"`
	DropoffLongitude float64 `json:"dropoff_longitude" binding:"required"`
	PickupAddress    string  `json:"pickup_address"`
	DropoffAddress   string  `json:"dropoff_address"`
}

type UpdateLocationRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

type CreateReviewRequest struct {
	RideID     string  `json:"ride_id" binding:"required"`
	ReviewedID string  `json:"reviewed_id" binding:"required"`
	Rating     float64 `json:"rating" binding:"required,min=1,max=5"`
	Comment    string  `json:"comment"`
}

func (h *RideHandler) CreateRideRequest(c *gin.Context) {
	currentUserID := c.GetString("user_id")
	currentUserType := c.GetString("user_type")

	if currentUserType != "passenger" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only passengers can create ride requests"})
		return
	}

	var req CreateRideRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse user ID
	userUUID, err := uuid.Parse(currentUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Calculate estimated fare and distance
	distance := utils.CalculateDistance(req.PickupLatitude, req.PickupLongitude, req.DropoffLatitude, req.DropoffLongitude)
	estimatedFare := calculateFare(distance)
	estimatedDuration := int(distance * 3) // Rough estimate: 3 minutes per km

	// Create ride request
	rideRequest := models.RideRequest{
		PassengerID:              userUUID,
		PickupLatitude:           req.PickupLatitude,
		PickupLongitude:          req.PickupLongitude,
		DropoffLatitude:          req.DropoffLatitude,
		DropoffLongitude:         req.DropoffLongitude,
		PickupAddress:            req.PickupAddress,
		DropoffAddress:           req.DropoffAddress,
		Status:                   "pending",
		EstimatedFare:            &estimatedFare,
		EstimatedDistanceKm:      &distance,
		EstimatedDurationMinutes: &estimatedDuration,
	}

	if err := h.db.Create(&rideRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ride request"})
		return
	}

	// Load passenger details
	if err := h.db.Preload("Passenger").First(&rideRequest, rideRequest.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load ride request details"})
		return
	}

	c.JSON(http.StatusCreated, rideRequest)
}

func (h *RideHandler) GetRideRequest(c *gin.Context) {
	requestID := c.Param("id")

	var rideRequest models.RideRequest
	if err := h.db.Preload("Passenger").Where("id = ?", requestID).First(&rideRequest).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ride request not found"})
		return
	}

	c.JSON(http.StatusOK, rideRequest)
}

func (h *RideHandler) GetNearbyDrivers(c *gin.Context) {
	latStr := c.Query("latitude")
	lonStr := c.Query("longitude")
	radiusStr := c.DefaultQuery("radius", "5") // Default 5km radius

	if latStr == "" || lonStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Latitude and longitude are required"})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
		return
	}

	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid radius"})
		return
	}

	// Find nearby available drivers
	var drivers []models.Driver
	if err := h.db.Preload("User").Where("is_available = ? AND is_approved = ? AND current_latitude IS NOT NULL AND current_longitude IS NOT NULL", true, true).Find(&drivers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find drivers"})
		return
	}

	// Filter drivers within radius
	var nearbyDrivers []models.Driver
	for _, driver := range drivers {
		if driver.CurrentLatitude != nil && driver.CurrentLongitude != nil {
			distance := utils.CalculateDistance(lat, lon, *driver.CurrentLatitude, *driver.CurrentLongitude)
			if distance <= radius {
				nearbyDrivers = append(nearbyDrivers, driver)
			}
		}
	}

	c.JSON(http.StatusOK, nearbyDrivers)
}

func (h *RideHandler) AcceptRideRequest(c *gin.Context) {
	requestID := c.Param("id")
	currentUserID := c.GetString("user_id")
	currentUserType := c.GetString("user_type")

	if currentUserType != "driver" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only drivers can accept ride requests"})
		return
	}

	// Parse user ID
	driverUUID, err := uuid.Parse(currentUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if driver exists and is approved
	var driver models.Driver
	if err := h.db.Where("driver_id = ? AND is_approved = ? AND is_available = ?", driverUUID, true, true).First(&driver).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Driver not found, not approved, or not available"})
		return
	}

	// Get ride request
	var rideRequest models.RideRequest
	if err := h.db.Where("id = ? AND status = ?", requestID, "pending").First(&rideRequest).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ride request not found or already processed"})
		return
	}

	// Start transaction
	tx := h.db.Begin()

	// Update ride request status
	if err := tx.Model(&rideRequest).Update("status", "accepted").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ride request"})
		return
	}

	// Create ride
	ride := models.Ride{
		RequestID:   rideRequest.ID,
		DriverID:    driverUUID,
		PassengerID: rideRequest.PassengerID,
		Status:      "in_progress",
	}

	if err := tx.Create(&ride).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ride"})
		return
	}

	// Update driver availability
	if err := tx.Model(&driver).Update("is_available", false).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update driver availability"})
		return
	}

	tx.Commit()

	// Load ride details
	if err := h.db.Preload("RideRequest").Preload("Driver.User").Preload("Passenger").First(&ride, ride.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load ride details"})
		return
	}

	c.JSON(http.StatusOK, ride)
}

func (h *RideHandler) RejectRideRequest(c *gin.Context) {
	requestID := c.Param("id")
	currentUserType := c.GetString("user_type")

	if currentUserType != "driver" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only drivers can reject ride requests"})
		return
	}

	// Get ride request
	var rideRequest models.RideRequest
	if err := h.db.Where("id = ? AND status = ?", requestID, "pending").First(&rideRequest).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ride request not found or already processed"})
		return
	}

	// Update ride request status
	if err := h.db.Model(&rideRequest).Update("status", "rejected").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ride request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ride request rejected"})
}

func (h *RideHandler) StartRide(c *gin.Context) {
	rideID := c.Param("id")
	currentUserID := c.GetString("user_id")
	currentUserType := c.GetString("user_type")

	if currentUserType != "driver" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only drivers can start rides"})
		return
	}

	// Parse user ID
	driverUUID, err := uuid.Parse(currentUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get ride
	var ride models.Ride
	if err := h.db.Where("id = ? AND driver_id = ? AND status = ?", rideID, driverUUID, "in_progress").First(&ride).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ride not found or not authorized"})
		return
	}

	// Update ride start time
	now := time.Now()
	if err := h.db.Model(&ride).Update("start_time", now).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start ride"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ride started", "start_time": now})
}

func (h *RideHandler) EndRide(c *gin.Context) {
	rideID := c.Param("id")
	currentUserID := c.GetString("user_id")
	currentUserType := c.GetString("user_type")

	if currentUserType != "driver" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only drivers can end rides"})
		return
	}

	// Parse user ID
	driverUUID, err := uuid.Parse(currentUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get ride
	var ride models.Ride
	if err := h.db.Preload("RideRequest").Where("id = ? AND driver_id = ? AND status = ?", rideID, driverUUID, "in_progress").First(&ride).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ride not found or not authorized"})
		return
	}

	if ride.StartTime == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ride has not been started yet"})
		return
	}

	// Calculate actual fare and duration
	now := time.Now()
	duration := int(now.Sub(*ride.StartTime).Minutes())
	actualFare := *ride.RideRequest.EstimatedFare // Use estimated fare for now
	actualDistance := *ride.RideRequest.EstimatedDistanceKm

	// Start transaction
	tx := h.db.Begin()

	// Update ride
	updates := map[string]interface{}{
		"end_time":                now,
		"status":                  "completed",
		"actual_fare":             actualFare,
		"actual_distance_km":      actualDistance,
		"actual_duration_minutes": duration,
	}

	if err := tx.Model(&ride).Updates(updates).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to end ride"})
		return
	}

	// Update ride request status
	if err := tx.Model(&ride.RideRequest).Update("status", "completed").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ride request"})
		return
	}

	// Update driver availability
	if err := tx.Model(&models.Driver{}).Where("driver_id = ?", driverUUID).Update("is_available", true).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update driver availability"})
		return
	}

	// Create payment record
	payment := models.Payment{
		RideID:        ride.ID,
		Amount:        actualFare,
		Currency:      "KES",
		PaymentMethod: "mpesa", // Default to M-Pesa
		PaymentStatus: "pending",
	}

	if err := tx.Create(&payment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment record"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message":     "Ride completed",
		"ride":        ride,
		"payment_id":  payment.ID,
		"total_fare":  actualFare,
	})
}

func (h *RideHandler) GetRide(c *gin.Context) {
	rideID := c.Param("id")

	var ride models.Ride
	if err := h.db.Preload("RideRequest").Preload("Driver.User").Preload("Passenger").Where("id = ?", rideID).First(&ride).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ride not found"})
		return
	}

	c.JSON(http.StatusOK, ride)
}

func (h *RideHandler) GetUserRides(c *gin.Context) {
	userID := c.Param("id")
	currentUserID := c.GetString("user_id")

	// Users can only view their own rides
	if userID != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only view your own rides"})
		return
	}

	var rides []models.Ride
	if err := h.db.Preload("RideRequest").Preload("Driver.User").Preload("Passenger").Where("passenger_id = ? OR driver_id = ?", userID, userID).Order("created_at DESC").Find(&rides).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rides"})
		return
	}

	c.JSON(http.StatusOK, rides)
}

func (h *RideHandler) UpdateDriverLocation(c *gin.Context) {
	driverID := c.Param("id")
	currentUserID := c.GetString("user_id")
	currentUserType := c.GetString("user_type")

	if currentUserType != "driver" || driverID != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own location"})
		return
	}

	var req UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse user ID
	driverUUID, err := uuid.Parse(driverID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid driver ID"})
		return
	}

	// Update driver location
	now := time.Now()
	updates := map[string]interface{}{
		"current_latitude":       req.Latitude,
		"current_longitude":      req.Longitude,
		"last_location_update":   now,
	}

	if err := h.db.Model(&models.Driver{}).Where("driver_id = ?", driverUUID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update location"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Location updated successfully"})
}

func (h *RideHandler) GetDriverLocation(c *gin.Context) {
	driverID := c.Param("id")

	var driver models.Driver
	if err := h.db.Where("driver_id = ?", driverID).First(&driver).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
		return
	}

	if driver.CurrentLatitude == nil || driver.CurrentLongitude == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver location not available"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"latitude":            *driver.CurrentLatitude,
		"longitude":           *driver.CurrentLongitude,
		"last_location_update": driver.LastLocationUpdate,
	})
}

func (h *RideHandler) CreateReview(c *gin.Context) {
	currentUserID := c.GetString("user_id")

	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse UUIDs
	rideUUID, err := uuid.Parse(req.RideID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ride ID"})
		return
	}

	reviewerUUID, err := uuid.Parse(currentUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reviewer ID"})
		return
	}

	reviewedUUID, err := uuid.Parse(req.ReviewedID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reviewed user ID"})
		return
	}

	// Verify ride exists and user was part of it
	var ride models.Ride
	if err := h.db.Where("id = ? AND (passenger_id = ? OR driver_id = ?) AND status = ?", rideUUID, reviewerUUID, reviewerUUID, "completed").First(&ride).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ride not found or not authorized to review"})
		return
	}

	// Check if review already exists
	var existingReview models.Review
	if err := h.db.Where("ride_id = ? AND reviewer_id = ?", rideUUID, reviewerUUID).First(&existingReview).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Review already exists for this ride"})
		return
	}

	// Create review
	review := models.Review{
		RideID:     rideUUID,
		ReviewerID: reviewerUUID,
		ReviewedID: reviewedUUID,
		Rating:     req.Rating,
		Comment:    req.Comment,
	}

	if err := h.db.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	// Update user rating
	h.updateUserRating(reviewedUUID)

	c.JSON(http.StatusCreated, review)
}

func (h *RideHandler) GetUserReviews(c *gin.Context) {
	userID := c.Param("id")

	var reviews []models.Review
	if err := h.db.Preload("Reviewer").Preload("Ride").Where("reviewed_id = ?", userID).Order("created_at DESC").Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// Helper functions
func calculateFare(distanceKm float64) float64 {
	// Basic fare calculation for Kenya
	baseFare := 50.0  // KES 50 base fare
	perKmRate := 25.0 // KES 25 per km
	return baseFare + (distanceKm * perKmRate)
}

func (h *RideHandler) updateUserRating(userID uuid.UUID) {
	var avgRating float64
	h.db.Model(&models.Review{}).Where("reviewed_id = ?", userID).Select("AVG(rating)").Scan(&avgRating)
	h.db.Model(&models.User{}).Where("id = ?", userID).Update("rating", avgRating)
}

