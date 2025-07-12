package email

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

type EmailService struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

func NewEmailService() *EmailService {
	return &EmailService{
		SMTPHost:     getEnv("EMAIL_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnv("EMAIL_PORT", "587"),
		SMTPUsername: getEnv("EMAIL_USERNAME", ""),
		SMTPPassword: getEnv("EMAIL_PASSWORD", ""),
		FromEmail:    getEnv("EMAIL_FROM", "noreply@kenyanrideshare.com"),
		FromName:     getEnv("EMAIL_FROM_NAME", "Kenyan Ride Share"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (es *EmailService) SendVerificationEmail(to, firstName, verificationToken string) error {
	subject := "Verify Your Email Address"
	verificationURL := fmt.Sprintf("%s/api/v1/auth/verify-email?token=%s", 
		getEnv("BASE_URL", "http://localhost:8080"), verificationToken)
	
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Email Verification</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #007bff; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .button { display: inline-block; padding: 12px 24px; background-color: #007bff; color: white; text-decoration: none; border-radius: 5px; margin: 20px 0; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to Kenyan Ride Share!</h1>
        </div>
        <div class="content">
            <h2>Hi %s,</h2>
            <p>Thank you for registering with Kenyan Ride Share. To complete your registration, please verify your email address by clicking the button below:</p>
            <p style="text-align: center;">
                <a href="%s" class="button">Verify Email Address</a>
            </p>
            <p>If the button doesn't work, you can also copy and paste this link into your browser:</p>
            <p style="word-break: break-all; color: #007bff;">%s</p>
            <p><strong>Important:</strong> This verification link will expire in 24 hours.</p>
            <p>If you didn't create an account with us, please ignore this email.</p>
        </div>
        <div class="footer">
            <p>© 2024 Kenyan Ride Share. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`, firstName, verificationURL, verificationURL)

	return es.sendEmail(to, subject, body)
}

func (es *EmailService) SendPasswordResetEmail(to, firstName, resetToken string) error {
	subject := "Reset Your Password"
	resetURL := fmt.Sprintf("%s/api/v1/auth/reset-password?token=%s", 
		getEnv("BASE_URL", "http://localhost:8080"), resetToken)
	
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Password Reset</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #dc3545; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .button { display: inline-block; padding: 12px 24px; background-color: #dc3545; color: white; text-decoration: none; border-radius: 5px; margin: 20px 0; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Password Reset Request</h1>
        </div>
        <div class="content">
            <h2>Hi %s,</h2>
            <p>We received a request to reset your password for your Kenyan Ride Share account.</p>
            <p>If you made this request, click the button below to reset your password:</p>
            <p style="text-align: center;">
                <a href="%s" class="button">Reset Password</a>
            </p>
            <p>If the button doesn't work, you can also copy and paste this link into your browser:</p>
            <p style="word-break: break-all; color: #dc3545;">%s</p>
            <p><strong>Important:</strong> This password reset link will expire in 1 hour.</p>
            <p>If you didn't request a password reset, please ignore this email. Your password will remain unchanged.</p>
        </div>
        <div class="footer">
            <p>© 2024 Kenyan Ride Share. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`, firstName, resetURL, resetURL)

	return es.sendEmail(to, subject, body)
}

func (es *EmailService) sendEmail(to, subject, body string) error {
	// Skip sending emails if SMTP credentials are not configured
	if es.SMTPUsername == "" || es.SMTPPassword == "" {
		fmt.Printf("Email would be sent to %s: %s\n", to, subject)
		fmt.Printf("Body: %s\n", body)
		return nil
	}

	auth := smtp.PlainAuth("", es.SMTPUsername, es.SMTPPassword, es.SMTPHost)

	msg := []string{
		fmt.Sprintf("From: %s <%s>", es.FromName, es.FromEmail),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=UTF-8",
		"",
		body,
	}

	message := strings.Join(msg, "\r\n")
	addr := fmt.Sprintf("%s:%s", es.SMTPHost, es.SMTPPort)

	return smtp.SendMail(addr, auth, es.FromEmail, []string{to}, []byte(message))
}
