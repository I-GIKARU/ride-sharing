package handlers

import (
	"net/http"
	"time"

	"kenyan-ride-share-backend/internal/models"
	"kenyan-ride-share-backend/pkg/email"
	"kenyan-ride-share-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandler struct {
	db           *gorm.DB
	emailService *email.EmailService
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		db:           db,
		emailService: email.NewEmailService(),
	}
}

type RegisterRequest struct {
	UserType    string `json:"user_type" binding:"required,oneof=driver passenger"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type OnboardDriverRequest struct {
	VehicleMake           string `json:"vehicle_make" binding:"required"`
	VehicleModel          string `json:"vehicle_model" binding:"required"`
	LicensePlate          string `json:"license_plate" binding:"required"`
	DriverLicenseNumber   string `json:"driver_license_number" binding:"required"`
	InsuranceDetails      string `json:"insurance_details"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := h.db.Where("email = ? OR phone_number = ?", req.Email, req.PhoneNumber).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email or phone number already exists"})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Generate email verification token
	verificationToken, err := utils.GenerateRandomString(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate verification token"})
		return
	}

	// Set verification token expiry (24 hours)
	verificationExpiry := time.Now().Add(24 * time.Hour)

	// Create user with email verification fields
	user := models.User{
		UserType:                req.UserType,
		FirstName:               req.FirstName,
		LastName:                req.LastName,
		Email:                   req.Email,
		PhoneNumber:             req.PhoneNumber,
		PasswordHash:            hashedPassword,
		IsEmailVerified:         false,
		EmailVerificationToken:  &verificationToken,
		EmailVerificationExpiry: &verificationExpiry,
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Send verification email
	if err := h.emailService.SendVerificationEmail(user.Email, user.FirstName, verificationToken); err != nil {
		// Log error but don't fail registration
		c.JSON(http.StatusCreated, gin.H{
			"message": "User registered successfully. Please check your email to verify your account. Note: Email sending failed, please contact support.",
			"user_id": user.ID,
			"email_verified": false,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully. Please check your email to verify your account.",
		"user_id": user.ID,
		"email_verified": false,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check if email is verified
	if !user.IsEmailVerified {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Email not verified. Please check your email and verify your account before logging in.",
			"email_verified": false,
		})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID.String(), user.UserType, "your-secret-key") // TODO: Use config
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    user,
		"token":   token,
		"email_verified": true,
	})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	
	var user models.User
	if err := h.db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	currentUserID := c.GetString("user_id")

	// Check if user is updating their own profile
	if userID != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own profile"})
		return
	}

	var user models.User
	if err := h.db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Remove sensitive fields that shouldn't be updated directly
	delete(updateData, "id")
	delete(updateData, "password_hash")
	delete(updateData, "user_type")
	delete(updateData, "created_at")

	updateData["updated_at"] = time.Now()

	if err := h.db.Model(&user).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) OnboardDriver(c *gin.Context) {
	currentUserID := c.GetString("user_id")
	currentUserType := c.GetString("user_type")

	if currentUserType != "driver" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only drivers can be onboarded"})
		return
	}

	var req OnboardDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if driver already exists
	var existingDriver models.Driver
	if err := h.db.Where("driver_id = ?", currentUserID).First(&existingDriver).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Driver already onboarded"})
		return
	}

	// Parse user ID
	userUUID, err := uuid.Parse(currentUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Create driver profile
	driver := models.Driver{
		DriverID:            userUUID,
		VehicleMake:         req.VehicleMake,
		VehicleModel:        req.VehicleModel,
		LicensePlate:        req.LicensePlate,
		DriverLicenseNumber: req.DriverLicenseNumber,
		InsuranceDetails:    req.InsuranceDetails,
		IsApproved:          false, // Requires admin approval
		IsAvailable:         false,
	}

	if err := h.db.Create(&driver).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to onboard driver"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Driver onboarded successfully. Awaiting approval.",
		"driver":  driver,
	})
}

// VerifyEmail handles email verification
func (h *UserHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification token is required"})
		return
	}

	// Find user by verification token
	var user models.User
	if err := h.db.Where("email_verification_token = ?", token).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification token"})
		return
	}

	// Check if already verified
	if user.IsEmailVerified {
		c.JSON(http.StatusOK, gin.H{"message": "Email already verified"})
		return
	}

	// Check if token has expired
	if user.EmailVerificationExpiry != nil && time.Now().After(*user.EmailVerificationExpiry) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification token has expired"})
		return
	}

	// Update user as verified
	if err := h.db.Model(&user).Updates(map[string]interface{}{
		"is_email_verified":         true,
		"email_verification_token":  nil,
		"email_verification_expiry": nil,
		"updated_at":                time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Email verified successfully",
		"email_verified": true,
	})
}

// ResendVerificationEmail resends the verification email
func (h *UserHandler) ResendVerificationEmail(c *gin.Context) {
	var req ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if already verified
	if user.IsEmailVerified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already verified"})
		return
	}

	// Generate new verification token
	verificationToken, err := utils.GenerateRandomString(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate verification token"})
		return
	}

	// Set new verification token expiry (24 hours)
	verificationExpiry := time.Now().Add(24 * time.Hour)

	// Update user with new verification token
	if err := h.db.Model(&user).Updates(map[string]interface{}{
		"email_verification_token":  verificationToken,
		"email_verification_expiry": verificationExpiry,
		"updated_at":                time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update verification token"})
		return
	}

	// Send verification email
	if err := h.emailService.SendVerificationEmail(user.Email, user.FirstName, verificationToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Verification email sent successfully",
	})
}

// ForgotPassword handles password reset requests
func (h *UserHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// Don't reveal if user exists or not for security
		c.JSON(http.StatusOK, gin.H{
			"message": "If the email exists, a password reset link has been sent",
		})
		return
	}

	// Generate password reset token
	resetToken, err := utils.GenerateRandomString(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset token"})
		return
	}

	// Set reset token expiry (1 hour)
	resetExpiry := time.Now().Add(1 * time.Hour)

	// Update user with reset token
	if err := h.db.Model(&user).Updates(map[string]interface{}{
		"password_reset_token":  resetToken,
		"password_reset_expiry": resetExpiry,
		"updated_at":            time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reset token"})
		return
	}

	// Send password reset email
	if err := h.emailService.SendPasswordResetEmail(user.Email, user.FirstName, resetToken); err != nil {
		// Log error but don't reveal to user
		c.JSON(http.StatusOK, gin.H{
			"message": "If the email exists, a password reset link has been sent",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "If the email exists, a password reset link has been sent",
	})
}

// ResetPassword handles password reset with token
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by reset token
	var user models.User
	if err := h.db.Where("password_reset_token = ?", req.Token).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reset token"})
		return
	}

	// Check if token has expired
	if user.PasswordResetExpiry != nil && time.Now().After(*user.PasswordResetExpiry) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reset token has expired"})
		return
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Update user password and clear reset token
	if err := h.db.Model(&user).Updates(map[string]interface{}{
		"password_hash":         hashedPassword,
		"password_reset_token":  nil,
		"password_reset_expiry": nil,
		"updated_at":            time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully",
	})
}

