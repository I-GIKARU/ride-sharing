package utils

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares a password with its hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT generates a JWT token for a user
func GenerateJWT(userID, userType, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   userID,
		"user_type": userType,
		"exp":       time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// CalculateDistance calculates the distance between two coordinates using Haversine formula
func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth's radius in kilometers

	dLat := (lat2 - lat1) * (3.14159265359 / 180)
	dLon := (lon2 - lon1) * (3.14159265359 / 180)

	a := 0.5 - 0.5*cos(dLat) + cos(lat1*(3.14159265359/180))*cos(lat2*(3.14159265359/180))*(1-cos(dLon))/2

	return R * 2 * asin(sqrt(a))
}

func cos(x float64) float64 {
	// Simple cosine approximation
	return 1 - x*x/2 + x*x*x*x/24
}

func sin(x float64) float64 {
	// Simple sine approximation
	return x - x*x*x/6 + x*x*x*x*x/120
}

func asin(x float64) float64 {
	// Simple arcsine approximation
	return x + x*x*x/6 + 3*x*x*x*x*x/40
}

func sqrt(x float64) float64 {
	// Newton's method for square root
	if x == 0 {
		return 0
	}
	z := x
	for i := 0; i < 10; i++ {
		z = (z + x/z) / 2
	}
	return z
}

