package services

import (
	"fmt"
	"kenyan-ride-share-backend/internal/models"
	"kenyan-ride-share-backend/pkg/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ComplianceService struct {
	db *gorm.DB
}

func NewComplianceService(db *gorm.DB) *ComplianceService {
	return &ComplianceService{db: db}
}

// ValidateDriverCompliance checks if driver meets Kenya NTSA requirements
func (c *ComplianceService) ValidateDriverCompliance(driverID uuid.UUID) (*DriverComplianceStatus, error) {
	var driver models.Driver
	if err := c.db.Where("driver_id = ?", driverID).First(&driver).Error; err != nil {
		return nil, err
	}

	status := &DriverComplianceStatus{
		DriverID:     driverID,
		IsCompliant:  true,
		Issues:       []string{},
		LastChecked:  utils.TimeNow(),
	}

	// Check if driver is approved
	if !driver.IsApproved {
		status.IsCompliant = false
		status.Issues = append(status.Issues, "Driver not approved by NTSA")
	}

	// Check if driver license is provided
	if driver.DriverLicenseNumber == "" {
		status.IsCompliant = false
		status.Issues = append(status.Issues, "Driver license number missing")
	}

	// Check if vehicle details are complete
	if driver.LicensePlate == "" {
		status.IsCompliant = false
		status.Issues = append(status.Issues, "Vehicle license plate missing")
	}

	if driver.VehicleMake == "" || driver.VehicleModel == "" {
		status.IsCompliant = false
		status.Issues = append(status.Issues, "Vehicle details incomplete")
	}

	// Check if insurance details are provided
	if driver.InsuranceDetails == "" {
		status.IsCompliant = false
		status.Issues = append(status.Issues, "Insurance details missing")
	}

	return status, nil
}

// CalculateCommission ensures compliance with 18% service fee cap
func (c *ComplianceService) CalculateCommission(fareAmount float64) *CommissionBreakdown {
	const maxCommissionRate = 0.18 // 18% as per Kenya regulations

	commission := fareAmount * maxCommissionRate
	driverEarnings := fareAmount - commission

	return &CommissionBreakdown{
		TotalFare:        fareAmount,
		CommissionRate:   maxCommissionRate,
		CommissionAmount: commission,
		DriverEarnings:   driverEarnings,
		Currency:         "KES",
	}
}

// GenerateNTSAReport creates compliance report for NTSA
func (c *ComplianceService) GenerateNTSAReport(startDate, endDate string) (*NTSAReport, error) {
	var rides []models.Ride
	if err := c.db.Where("created_at BETWEEN ? AND ? AND status = ?", startDate, endDate, "completed").
		Find(&rides).Error; err != nil {
		return nil, err
	}

	report := &NTSAReport{
		ReportPeriod: fmt.Sprintf("%s to %s", startDate, endDate),
		GeneratedAt:  utils.TimeNow(),
		TotalRides:   len(rides),
		TotalRevenue: 0,
		Rides:        []NTSARideData{},
	}

	for _, ride := range rides {
		if ride.ActualFare != nil {
			report.TotalRevenue += *ride.ActualFare
		}

		// Get driver details
		var driver models.Driver
		var user models.User
		driverName := "Unknown Driver"
		if err := c.db.Where("driver_id = ?", ride.DriverID).First(&driver).Error; err == nil {
			if err := c.db.Where("id = ?", ride.DriverID).First(&user).Error; err == nil {
				driverName = user.FirstName + " " + user.LastName
			}
		}

		// Get ride request details
		var rideRequest models.RideRequest
		pickupLocation := "Unknown"
		dropoffLocation := "Unknown"
		if err := c.db.Where("id = ?", ride.RequestID).First(&rideRequest).Error; err == nil {
			pickupLocation = fmt.Sprintf("%.6f,%.6f", rideRequest.PickupLatitude, rideRequest.PickupLongitude)
			dropoffLocation = fmt.Sprintf("%.6f,%.6f", rideRequest.DropoffLatitude, rideRequest.DropoffLongitude)
		}

		rideData := NTSARideData{
			RideID:          ride.ID.String(),
			DriverID:        ride.DriverID.String(),
			DriverName:      driverName,
			PassengerID:     ride.PassengerID.String(),
			StartTime:       ride.StartTime,
			EndTime:         ride.EndTime,
			PickupLocation:  pickupLocation,
			DropoffLocation: dropoffLocation,
			Distance:        ride.ActualDistanceKm,
			Duration:        ride.ActualDurationMinutes,
			Fare:            ride.ActualFare,
		}

		if ride.ActualFare != nil {
			commission := c.CalculateCommission(*ride.ActualFare)
			rideData.CommissionAmount = &commission.CommissionAmount
		}

		report.Rides = append(report.Rides, rideData)
	}

	return report, nil
}

// ValidateVehicleEligibility checks vehicle against Kenya requirements
func (c *ComplianceService) ValidateVehicleEligibility(vehicleYear int, vehicleMake, vehicleModel string) *VehicleEligibilityStatus {
	currentYear := utils.CurrentYear()
	vehicleAge := currentYear - vehicleYear

	status := &VehicleEligibilityStatus{
		IsEligible: true,
		Issues:     []string{},
		MaxAge:     10, // Kenya typically allows vehicles up to 10 years old
	}

	// Check vehicle age (as per Uber Kenya 2025 requirements: 2015 and newer)
	if vehicleYear < 2015 {
		status.IsEligible = false
		status.Issues = append(status.Issues, fmt.Sprintf("Vehicle too old: %d years (max 10 years)", vehicleAge))
	}

	// Check if vehicle details are provided
	if vehicleMake == "" || vehicleModel == "" {
		status.IsEligible = false
		status.Issues = append(status.Issues, "Vehicle make and model required")
	}

	return status
}

// Data structures for compliance
type DriverComplianceStatus struct {
	DriverID    uuid.UUID `json:"driver_id"`
	IsCompliant bool      `json:"is_compliant"`
	Issues      []string  `json:"issues"`
	LastChecked string    `json:"last_checked"`
}

type CommissionBreakdown struct {
	TotalFare        float64 `json:"total_fare"`
	CommissionRate   float64 `json:"commission_rate"`
	CommissionAmount float64 `json:"commission_amount"`
	DriverEarnings   float64 `json:"driver_earnings"`
	Currency         string  `json:"currency"`
}

type NTSAReport struct {
	ReportPeriod string         `json:"report_period"`
	GeneratedAt  string         `json:"generated_at"`
	TotalRides   int            `json:"total_rides"`
	TotalRevenue float64        `json:"total_revenue"`
	Rides        []NTSARideData `json:"rides"`
}

type NTSARideData struct {
	RideID           string     `json:"ride_id"`
	DriverID         string     `json:"driver_id"`
	DriverName       string     `json:"driver_name"`
	PassengerID      string     `json:"passenger_id"`
	StartTime        *time.Time `json:"start_time"`
	EndTime          *time.Time `json:"end_time"`
	PickupLocation   string     `json:"pickup_location"`
	DropoffLocation  string     `json:"dropoff_location"`
	Distance         *float64   `json:"distance_km"`
	Duration         *int       `json:"duration_minutes"`
	Fare             *float64   `json:"fare"`
	CommissionAmount *float64   `json:"commission_amount"`
}

type VehicleEligibilityStatus struct {
	IsEligible bool     `json:"is_eligible"`
	Issues     []string `json:"issues"`
	MaxAge     int      `json:"max_age_years"`
}

