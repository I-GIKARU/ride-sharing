package database

import (
	"kenyan-ride-share-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Initialize(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(
		&models.User{},
		&models.Driver{},
		&models.RideRequest{},
		&models.Ride{},
		&models.Payment{},
		&models.Review{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

