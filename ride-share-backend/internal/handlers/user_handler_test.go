package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"kenyan-ride-share-backend/internal/handlers"
	"kenyan-ride-share-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	
	// Auto-migrate tables
	db.AutoMigrate(&models.User{}, &models.Driver{}, &models.RideRequest{}, &models.Ride{}, &models.Payment{}, &models.Review{})
	
	return db
}

func TestUserRegistration(t *testing.T) {
	db := setupTestDB()
	userHandler := handlers.NewUserHandler(db)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/register", userHandler.Register)

	// Test valid registration
	user := map[string]interface{}{
		"first_name":   "John",
		"last_name":    "Doe",
		"email":        "john.doe@example.com",
		"phone_number": "254712345678",
		"password":     "password123",
		"user_type":    "passenger",
	}

	jsonData, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "User registered successfully", response["message"])
	assert.NotEmpty(t, response["user_id"])
}

func TestUserLogin(t *testing.T) {
	db := setupTestDB()
	userHandler := handlers.NewUserHandler(db)

	// Create a test user first
	testUser := models.User{
		FirstName:   "Jane",
		LastName:    "Smith",
		Email:       "jane.smith@example.com",
		PhoneNumber: "254712345679",
		UserType:    "passenger",
	}
	// Set password hash manually for testing
	testUser.PasswordHash = "$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj/VjPyV8Paa" // "password123"
	db.Create(&testUser)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/login", userHandler.Login)

	// Test valid login
	loginData := map[string]string{
		"email":    "jane.smith@example.com",
		"password": "password123",
	}

	jsonData, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Login successful", response["message"])
	assert.NotEmpty(t, response["token"])
}

func TestInvalidLogin(t *testing.T) {
	db := setupTestDB()
	userHandler := handlers.NewUserHandler(db)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/login", userHandler.Login)

	// Test invalid login
	loginData := map[string]string{
		"email":    "nonexistent@example.com",
		"password": "wrongpassword",
	}

	jsonData, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

