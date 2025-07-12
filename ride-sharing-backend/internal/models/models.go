package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                    uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserType              string     `json:"user_type" gorm:"not null"` // 'driver' or 'passenger'
	FirstName             string     `json:"first_name" gorm:"not null"`
	LastName              string     `json:"last_name" gorm:"not null"`
	Email                 string     `json:"email" gorm:"unique;not null"`
	PhoneNumber           string     `json:"phone_number" gorm:"unique;not null"`
	PasswordHash          string     `json:"-" gorm:"not null"`
	IsEmailVerified       bool       `json:"is_email_verified" gorm:"default:false"`
	EmailVerificationToken *string   `json:"-" gorm:"unique"`
	EmailVerificationExpiry *time.Time `json:"-"`
	PasswordResetToken     *string   `json:"-" gorm:"unique"`
	PasswordResetExpiry    *time.Time `json:"-"`
	Rating                float64    `json:"rating" gorm:"default:0.0"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

type Driver struct {
	DriverID              uuid.UUID  `json:"driver_id" gorm:"type:uuid;primary_key"`
	VehicleMake           string     `json:"vehicle_make"`
	VehicleModel          string     `json:"vehicle_model"`
	LicensePlate          string     `json:"license_plate" gorm:"unique;not null"`
	DriverLicenseNumber   string     `json:"driver_license_number" gorm:"unique;not null"`
	InsuranceDetails      string     `json:"insurance_details"`
	IsApproved            bool       `json:"is_approved" gorm:"default:false"`
	IsAvailable           bool       `json:"is_available" gorm:"default:false"`
	CurrentLatitude       *float64   `json:"current_latitude"`
	CurrentLongitude      *float64   `json:"current_longitude"`
	LastLocationUpdate    *time.Time `json:"last_location_update"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

type RideRequest struct {
	ID                      uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	PassengerID             uuid.UUID `json:"passenger_id" gorm:"not null"`
	PickupLatitude          float64   `json:"pickup_latitude" gorm:"not null"`
	PickupLongitude         float64   `json:"pickup_longitude" gorm:"not null"`
	DropoffLatitude         float64   `json:"dropoff_latitude" gorm:"not null"`
	DropoffLongitude        float64   `json:"dropoff_longitude" gorm:"not null"`
	PickupAddress           string    `json:"pickup_address"`
	DropoffAddress          string    `json:"dropoff_address"`
	RequestedAt             time.Time `json:"requested_at" gorm:"default:CURRENT_TIMESTAMP"`
	Status                  string    `json:"status" gorm:"not null"` // 'pending', 'accepted', 'rejected', 'cancelled', 'completed'
	EstimatedFare           *float64  `json:"estimated_fare"`
	EstimatedDistanceKm     *float64  `json:"estimated_distance_km"`
	EstimatedDurationMinutes *int     `json:"estimated_duration_minutes"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

type Ride struct {
	ID                     uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RequestID              uuid.UUID    `json:"request_id" gorm:"unique;not null"`
	DriverID               uuid.UUID    `json:"driver_id" gorm:"not null"`
	PassengerID            uuid.UUID    `json:"passenger_id" gorm:"not null"`
	StartTime              *time.Time   `json:"start_time"`
	EndTime                *time.Time   `json:"end_time"`
	ActualFare             *float64     `json:"actual_fare"`
	ActualDistanceKm       *float64     `json:"actual_distance_km"`
	ActualDurationMinutes  *int         `json:"actual_duration_minutes"`
	RouteGeoJSON           string       `json:"route_geojson"`
	Status                 string       `json:"status" gorm:"not null"` // 'in_progress', 'completed', 'cancelled'
	CreatedAt              time.Time    `json:"created_at"`
	UpdatedAt              time.Time    `json:"updated_at"`
}

type Payment struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RideID        uuid.UUID  `json:"ride_id" gorm:"unique;not null"`
	Amount        float64    `json:"amount" gorm:"not null"`
	Currency      string     `json:"currency" gorm:"default:'KES'"`
	PaymentMethod string     `json:"payment_method" gorm:"not null"` // 'mpesa', 'card', 'cash'
	TransactionID *string    `json:"transaction_id" gorm:"unique"`   // M-Pesa transaction ID
	PaymentStatus string     `json:"payment_status" gorm:"not null"` // 'pending', 'completed', 'failed'
	PaymentDate   time.Time  `json:"payment_date" gorm:"default:CURRENT_TIMESTAMP"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type Review struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RideID     uuid.UUID `json:"ride_id" gorm:"not null"`
	ReviewerID uuid.UUID `json:"reviewer_id" gorm:"not null"` // Passenger or Driver
	ReviewedID uuid.UUID `json:"reviewed_id" gorm:"not null"` // Driver or Passenger
	Rating     float64   `json:"rating" gorm:"not null;check:rating >= 1.0 AND rating <= 5.0"`
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// BeforeCreate hook to generate UUID for models
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (d *Driver) BeforeCreate(tx *gorm.DB) error {
	if d.DriverID == uuid.Nil {
		d.DriverID = uuid.New()
	}
	return nil
}

func (rr *RideRequest) BeforeCreate(tx *gorm.DB) error {
	if rr.ID == uuid.Nil {
		rr.ID = uuid.New()
	}
	return nil
}

func (r *Ride) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (r *Review) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

