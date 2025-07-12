package handlers

import (
	"net/http"
	"time"

	"kenyan-ride-share-backend/internal/config"
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
	config       *config.Config
}

func NewUserHandler(db *gorm.DB, cfg *config.Config) *UserHandler {
	return &UserHandler{
		db:           db,
		emailService: email.NewEmailService(),
		config:       cfg,
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
	token, err := utils.GenerateJWT(user.ID.String(), user.UserType, h.config.JWTSecret)
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

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, `
	<!DOCTYPE html>
	<html lang="en">
	<head>
	    <meta charset="UTF-8">
	    <meta name="viewport" content="width=device-width, initial-scale=1.0">
	    <title>Email Verification</title>
	    <style>
	        body { font-family: Arial, sans-serif; text-align: center; padding: 50px; }
	        .content { max-width: 500px; margin: auto; }
	        h1 { color: #4CAF50; }
	        p { font-size: 1.2em; }
	    </style>
	</head>
	<body>
	    <div class="content">
	        <h1>Success!</h1>
	        <p>Your email has been verified successfully.</p>
	        <p><strong>Please return to the mobile app to continue.</strong></p>
	        <p style="font-size: 0.9em; color: #666;">You can close this browser window and open the Kenyan Ride Share app on your device.</p>
	    </div>
	</body>
	</html>
	`)
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

// ShowResetPasswordForm displays the password reset form
func (h *UserHandler) ShowResetPasswordForm(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusBadRequest, `
		<!DOCTYPE html>
		<html lang="en">
		<head>
		    <meta charset="UTF-8">
		    <meta name="viewport" content="width=device-width, initial-scale=1.0">
		    <title>Password Reset</title>
		    <style>
		        body { font-family: Arial, sans-serif; text-align: center; padding: 50px; }
		        .content { max-width: 500px; margin: auto; }
		        h1 { color: #f44336; }
		        p { font-size: 1.2em; }
		    </style>
		</head>
		<body>
		    <div class="content">
		        <h1>Error</h1>
		        <p>Invalid password reset link. Please request a new one.</p>
		    </div>
		</body>
		</html>
		`)
		return
	}

	// Check if token is valid
	var user models.User
	if err := h.db.Where("password_reset_token = ?", token).First(&user).Error; err != nil {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusBadRequest, `
		<!DOCTYPE html>
		<html lang="en">
		<head>
		    <meta charset="UTF-8">
		    <meta name="viewport" content="width=device-width, initial-scale=1.0">
		    <title>Password Reset</title>
		    <style>
		        body { font-family: Arial, sans-serif; text-align: center; padding: 50px; }
		        .content { max-width: 500px; margin: auto; }
		        h1 { color: #f44336; }
		        p { font-size: 1.2em; }
		    </style>
		</head>
		<body>
		    <div class="content">
		        <h1>Error</h1>
		        <p>Invalid or expired password reset link. Please request a new one.</p>
		    </div>
		</body>
		</html>
		`)
		return
	}

	// Check if token has expired
	if user.PasswordResetExpiry != nil && time.Now().After(*user.PasswordResetExpiry) {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusBadRequest, `
		<!DOCTYPE html>
		<html lang="en">
		<head>
		    <meta charset="UTF-8">
		    <meta name="viewport" content="width=device-width, initial-scale=1.0">
		    <title>Password Reset</title>
		    <style>
		        body { font-family: Arial, sans-serif; text-align: center; padding: 50px; }
		        .content { max-width: 500px; margin: auto; }
		        h1 { color: #f44336; }
		        p { font-size: 1.2em; }
		    </style>
		</head>
		<body>
		    <div class="content">
		        <h1>Error</h1>
		        <p>This password reset link has expired. Please request a new one.</p>
		    </div>
		</body>
		</html>
		`)
		return
	}

	// Display password reset form
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, `
	<!DOCTYPE html>
	<html lang="en">
	<head>
	    <meta charset="UTF-8">
	    <meta name="viewport" content="width=device-width, initial-scale=1.0">
	    <title>Reset Password</title>
	    <style>
	        body { font-family: Arial, sans-serif; text-align: center; padding: 50px; background-color: #f5f5f5; }
	        .content { max-width: 400px; margin: auto; background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
	        h1 { color: #333; margin-bottom: 30px; }
	        .form-group { margin-bottom: 20px; text-align: left; }
	        label { display: block; margin-bottom: 5px; font-weight: bold; color: #555; }
	        input[type="password"] { width: 100%%; padding: 12px; border: 1px solid #ddd; border-radius: 5px; font-size: 16px; }
	        button { width: 100%%; padding: 12px; background-color: #4CAF50; color: white; border: none; border-radius: 5px; font-size: 16px; cursor: pointer; }
	        button:hover { background-color: #45a049; }
	        .error { color: #f44336; margin-top: 10px; }
	        .success { color: #4CAF50; margin-top: 10px; }
	    </style>
	</head>
	<body>
	    <div class="content">
	        <h1>Reset Your Password</h1>
	        <form id="resetForm" method="POST" action="`+h.config.APIBasePath+`/auth/reset-password">
	            <input type="hidden" name="token" value="`+token+`">
	            <div class="form-group">
	                <label for="new_password">New Password:</label>
	                <input type="password" id="new_password" name="new_password" required minlength="6" placeholder="Enter your new password">
	            </div>
	            <div class="form-group">
	                <label for="confirm_password">Confirm Password:</label>
	                <input type="password" id="confirm_password" name="confirm_password" required minlength="6" placeholder="Confirm your new password">
	            </div>
	            <button type="submit">Reset Password</button>
	            <div id="message"></div>
	        </form>
	        <script>
	            document.getElementById('resetForm').addEventListener('submit', async function(e) {
	                e.preventDefault();
	                const formData = new FormData(this);
	                const password = formData.get('new_password');
	                const confirmPassword = formData.get('confirm_password');
	                const token = formData.get('token');
	                const messageDiv = document.getElementById('message');
	                
	                if (password !== confirmPassword) {
	                    messageDiv.innerHTML = '<div class="error">Passwords do not match!</div>';
	                    return;
	                }
	                
	                try {
	                    const response = await fetch('`+h.config.APIBasePath+`/auth/reset-password', {
	                        method: 'POST',
	                        headers: { 'Content-Type': 'application/json' },
	                        body: JSON.stringify({ token: token, new_password: password })
	                    });
	                    
	                    if (response.headers.get('content-type').includes('text/html')) {
	                        // Success - redirect to success page
	                        document.body.innerHTML = await response.text();
	                    } else {
	                        const result = await response.json();
	                        if (response.ok) {
	                            messageDiv.innerHTML = '<div class="success">' + result.message + '</div>';
	                        } else {
	                            messageDiv.innerHTML = '<div class="error">' + result.error + '</div>';
	                        }
	                    }
	                } catch (error) {
	                    messageDiv.innerHTML = '<div class="error">An error occurred. Please try again.</div>';
	                }
	            });
	        </script>
	    </div>
	</body>
	</html>
	`)
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

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, `
	<!DOCTYPE html>
	<html lang="en">
	<head>
	    <meta charset="UTF-8">
	    <meta name="viewport" content="width=device-width, initial-scale=1.0">
	    <title>Password Reset</title>
	    <style>
	        body { font-family: Arial, sans-serif; text-align: center; padding: 50px; }
	        .content { max-width: 500px; margin: auto; }
	        h1 { color: #4CAF50; }
	        p { font-size: 1.2em; }
	        form { margin-top: 20px; }
	        input { padding: 10px; width: 100%; margin-bottom: 10px; }
	        button { padding: 10px 20px; background-color: #4CAF50; color: white; border: none; }
	    </style>
	</head>
	<body>
	    <div class="content">
	        <h1>Password Reset</h1>
	        <p>Your password has been reset successfully.</p>
	        <p><strong>Please return to the mobile app to log in with your new password.</strong></p>
	        <p style="font-size: 0.9em; color: #666;">You can close this browser window and open the Kenyan Ride Share app on your device.</p>
	    </div>
	</body>
	</html>
	`)
}

