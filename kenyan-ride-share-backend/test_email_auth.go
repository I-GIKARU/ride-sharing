package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8080/api/v1"

// Test data structures
type TestResponse struct {
	Message       string      `json:"message"`
	UserID        string      `json:"user_id"`
	EmailVerified bool        `json:"email_verified"`
	Token         string      `json:"token"`
	User          interface{} `json:"user"`
	Error         string      `json:"error"`
}

func main() {
	fmt.Println("🧪 Starting Email Verification and Password Reset Tests")
	fmt.Println("============================================================")
	
	// Test 1: Register a new user
	fmt.Println("\n1. Testing User Registration with Email Verification")
	userEmail := "test@example.com"
	userID := testUserRegistration(userEmail)
	
	// Test 2: Try to login without email verification
	fmt.Println("\n2. Testing Login without Email Verification")
	testLoginWithoutVerification(userEmail)
	
	// Test 3: Resend verification email
	fmt.Println("\n3. Testing Resend Verification Email")
	testResendVerification(userEmail)
	
	// Test 4: Verify email (simulate clicking verification link)
	fmt.Println("\n4. Testing Email Verification")
	// Note: In real testing, you would need to extract the token from the email
	// For this demo, we'll simulate with a mock token
	testEmailVerification("mock-token")
	
	// Test 5: Request password reset
	fmt.Println("\n5. Testing Password Reset Request")
	testPasswordResetRequest(userEmail)
	
	// Test 6: Reset password
	fmt.Println("\n6. Testing Password Reset")
	testPasswordReset("mock-reset-token", "newpassword123")
	
	// Test 7: Login with verified email
	fmt.Println("\n7. Testing Login with Verified Email")
	testLoginWithVerifiedEmail(userEmail)
	
	fmt.Println("\n============================================================")
	fmt.Println("✅ All tests completed!")
	fmt.Println("\nNote: Some tests may fail because they use mock tokens.")
	fmt.Println("In a real scenario, you would extract tokens from the actual emails sent.")
}

func testUserRegistration(email string) string {
	fmt.Printf("📝 Registering user: %s\n", email)
	
	userData := map[string]interface{}{
		"user_type":    "passenger",
		"first_name":   "John",
		"last_name":    "Doe",
		"email":        email,
		"phone_number": "254712345678",
		"password":     "password123",
	}
	
	response := makeHTTPRequest("POST", "/register", userData)
	fmt.Printf("✅ Registration Response: %s\n", response.Message)
	
	if response.UserID != "" {
		fmt.Printf("👤 User ID: %s\n", response.UserID)
		return response.UserID
	}
	
	return ""
}

func testLoginWithoutVerification(email string) {
	fmt.Printf("🔒 Attempting login without email verification: %s\n", email)
	
	loginData := map[string]interface{}{
		"email":    email,
		"password": "password123",
	}
	
	response := makeHTTPRequest("POST", "/login", loginData)
	
	if response.Error != "" {
		fmt.Printf("❌ Login failed as expected: %s\n", response.Error)
	} else {
		fmt.Printf("⚠️  Login succeeded unexpectedly\n")
	}
}

func testResendVerification(email string) {
	fmt.Printf("📧 Resending verification email to: %s\n", email)
	
	resendData := map[string]interface{}{
		"email": email,
	}
	
	response := makeHTTPRequest("POST", "/auth/resend-verification", resendData)
	fmt.Printf("✅ Resend Response: %s\n", response.Message)
}

func testEmailVerification(token string) {
	fmt.Printf("✉️  Verifying email with token: %s\n", token)
	
	url := fmt.Sprintf("%s/auth/verify-email?token=%s", baseURL, token)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("❌ Error making verification request: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	var response TestResponse
	json.Unmarshal(body, &response)
	
	if response.Error != "" {
		fmt.Printf("❌ Verification failed: %s\n", response.Error)
	} else {
		fmt.Printf("✅ Verification Response: %s\n", response.Message)
	}
}

func testPasswordResetRequest(email string) {
	fmt.Printf("🔐 Requesting password reset for: %s\n", email)
	
	resetData := map[string]interface{}{
		"email": email,
	}
	
	response := makeHTTPRequest("POST", "/auth/forgot-password", resetData)
	fmt.Printf("✅ Password Reset Request: %s\n", response.Message)
}

func testPasswordReset(token, newPassword string) {
	fmt.Printf("🔑 Resetting password with token: %s\n", token)
	
	resetData := map[string]interface{}{
		"token":        token,
		"new_password": newPassword,
	}
	
	response := makeHTTPRequest("POST", "/auth/reset-password", resetData)
	
	if response.Error != "" {
		fmt.Printf("❌ Password reset failed: %s\n", response.Error)
	} else {
		fmt.Printf("✅ Password Reset: %s\n", response.Message)
	}
}

func testLoginWithVerifiedEmail(email string) {
	fmt.Printf("🔓 Attempting login with verified email: %s\n", email)
	
	loginData := map[string]interface{}{
		"email":    email,
		"password": "password123",
	}
	
	response := makeHTTPRequest("POST", "/login", loginData)
	
	if response.Token != "" {
		fmt.Printf("✅ Login successful! Token received.\n")
	} else if response.Error != "" {
		fmt.Printf("❌ Login failed: %s\n", response.Error)
	}
}

func makeHTTPRequest(method, endpoint string, data map[string]interface{}) TestResponse {
	jsonData, _ := json.Marshal(data)
	
	req, err := http.NewRequest(method, baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Error creating request: %v\n", err)
		return TestResponse{}
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Error making request: %v\n", err)
		return TestResponse{}
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	var response TestResponse
	json.Unmarshal(body, &response)
	
	return response
}
