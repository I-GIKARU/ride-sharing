package services_test

import (
	"testing"

	"kenyan-ride-share-backend/internal/models"
	"kenyan-ride-share-backend/internal/services"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.User{}, &models.Driver{}, &models.RideRequest{}, &models.Ride{}, &models.Payment{}, &models.Review{})
	return db
}

func TestMpesaService(t *testing.T) {
	db := setupTestDB()
	mpesaService := services.NewMpesaService(db)

	t.Run("GetAccessToken", func(t *testing.T) {
		token, err := mpesaService.GetAccessToken()
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Contains(t, token, "mock_access_token")
	})

	t.Run("InitiateSTKPush", func(t *testing.T) {
		response, err := mpesaService.InitiateSTKPush("254712345678", 100.0, "TEST-REF")
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "0", response.ResponseCode)
		assert.NotEmpty(t, response.CheckoutRequestID)
	})
}

func TestComplianceService(t *testing.T) {
	db := setupTestDB()
	complianceService := services.NewComplianceService(db)

	// Create test user and driver
	testUser := models.User{
		FirstName:   "Test",
		LastName:    "Driver",
		Email:       "test.driver@example.com",
		PhoneNumber: "254712345678",
		UserType:    "driver",
	}
	db.Create(&testUser)

	testDriver := models.Driver{
		DriverID:            testUser.ID,
		DriverLicenseNumber: "DL123456",
		LicensePlate:        "KCA123A",
		VehicleMake:         "Toyota",
		VehicleModel:        "Corolla",
		InsuranceDetails:    "Insurance Company XYZ",
		IsApproved:          true,
	}
	db.Create(&testDriver)

	t.Run("ValidateDriverCompliance", func(t *testing.T) {
		status, err := complianceService.ValidateDriverCompliance(testUser.ID)
		assert.NoError(t, err)
		assert.NotNil(t, status)
		assert.True(t, status.IsCompliant)
		assert.Empty(t, status.Issues)
	})

	t.Run("CalculateCommission", func(t *testing.T) {
		breakdown := complianceService.CalculateCommission(1000.0)
		assert.NotNil(t, breakdown)
		assert.Equal(t, 1000.0, breakdown.TotalFare)
		assert.Equal(t, 0.18, breakdown.CommissionRate)
		assert.Equal(t, 180.0, breakdown.CommissionAmount)
		assert.Equal(t, 820.0, breakdown.DriverEarnings)
		assert.Equal(t, "KES", breakdown.Currency)
	})

	t.Run("ValidateVehicleEligibility", func(t *testing.T) {
		status := complianceService.ValidateVehicleEligibility(2020, "Toyota", "Corolla")
		assert.NotNil(t, status)
		assert.True(t, status.IsEligible)
		assert.Empty(t, status.Issues)

		// Test old vehicle
		oldStatus := complianceService.ValidateVehicleEligibility(2010, "Toyota", "Corolla")
		assert.NotNil(t, oldStatus)
		assert.False(t, oldStatus.IsEligible)
		assert.NotEmpty(t, oldStatus.Issues)
	})
}

