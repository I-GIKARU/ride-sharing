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
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to Kenyan Ride Share</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; line-height: 1.6; color: #333; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); min-height: 100vh; padding: 20px; }
        .email-container { max-width: 600px; margin: 0 auto; background: white; border-radius: 20px; overflow: hidden; box-shadow: 0 20px 40px rgba(0,0,0,0.1); }
        .header { background: linear-gradient(135deg, #FF6B35 0%%, #F7931E 50%%, #FFD23F 100%%); padding: 40px 30px; text-align: center; position: relative; }
        .header::before { content: ''; position: absolute; top: 0; left: 0; right: 0; bottom: 0; background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><path d="M20,20 Q50,5 80,20 Q95,50 80,80 Q50,95 20,80 Q5,50 20,20" fill="none" stroke="rgba(255,255,255,0.1)" stroke-width="2"/></svg>') center/contain no-repeat; }
        .logo { color: white; font-size: 32px; font-weight: bold; margin-bottom: 10px; text-shadow: 0 2px 4px rgba(0,0,0,0.3); z-index: 2; position: relative; }
        .logo::before { content: '🚗'; margin-right: 10px; }
        .tagline { color: rgba(255,255,255,0.9); font-size: 16px; z-index: 2; position: relative; }
        .content { padding: 50px 30px; }
        .welcome-title { font-size: 28px; color: #2c3e50; margin-bottom: 20px; text-align: center; }
        .greeting { font-size: 20px; color: #34495e; margin-bottom: 25px; }
        .message { font-size: 16px; color: #555; margin-bottom: 30px; line-height: 1.8; }
        .cta-container { text-align: center; margin: 40px 0; }
        .verify-button { display: inline-block; padding: 18px 40px; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; text-decoration: none; border-radius: 50px; font-size: 16px; font-weight: 600; box-shadow: 0 8px 25px rgba(102, 126, 234, 0.4); transition: all 0.3s ease; text-transform: uppercase; letter-spacing: 1px; }
        .verify-button:hover { transform: translateY(-2px); box-shadow: 0 12px 35px rgba(102, 126, 234, 0.6); }
        .link-fallback { background: #f8f9fa; padding: 20px; border-radius: 10px; margin: 30px 0; border-left: 4px solid #667eea; }
        .link-fallback p { margin-bottom: 10px; color: #666; font-size: 14px; }
        .link-text { word-break: break-all; color: #667eea; font-family: monospace; background: white; padding: 10px; border-radius: 5px; border: 1px solid #e9ecef; }
        .warning-box { background: linear-gradient(135deg, #ffeaa7 0%%, #fab1a0 100%%); padding: 20px; border-radius: 10px; margin: 30px 0; text-align: center; }
        .warning-box strong { color: #d63031; }
        .features { display: flex; justify-content: space-around; margin: 40px 0; flex-wrap: wrap; }
        .feature { text-align: center; flex: 1; min-width: 150px; margin: 10px; }
        .feature-icon { font-size: 30px; margin-bottom: 10px; }
        .feature-text { font-size: 14px; color: #666; }
        .footer { background: #2c3e50; color: white; padding: 30px; text-align: center; }
        .footer-content { margin-bottom: 20px; }
        .social-links { margin: 20px 0; }
        .social-links a { color: white; text-decoration: none; margin: 0 15px; font-size: 18px; }
        .footer-note { font-size: 12px; color: #bdc3c7; margin-top: 20px; }
        @media (max-width: 600px) {
            .email-container { margin: 10px; border-radius: 15px; }
            .header { padding: 30px 20px; }
            .content { padding: 30px 20px; }
            .logo { font-size: 28px; }
            .welcome-title { font-size: 24px; }
            .features { flex-direction: column; }
        }
    </style>
</head>
<body>
    <div class="email-container">
        <div class="header">
            <div class="logo">Kenyan Ride Share</div>
            <div class="tagline">Your Journey, Our Priority</div>
        </div>
        
        <div class="content">
            <h1 class="welcome-title">Welcome Aboard! 🎉</h1>
            <p class="greeting">Hey %s! 👋</p>
            <p class="message">
                We're thrilled to have you join the Kenyan Ride Share family! You're just one step away from experiencing seamless, safe, and affordable transportation across Kenya.
            </p>
            
            <div class="features">
                <div class="feature">
                    <div class="feature-icon">🚗</div>
                    <div class="feature-text">Reliable Rides</div>
                </div>
                <div class="feature">
                    <div class="feature-icon">💰</div>
                    <div class="feature-text">Affordable Rates</div>
                </div>
                <div class="feature">
                    <div class="feature-icon">🛡️</div>
                    <div class="feature-text">Safe & Secure</div>
                </div>
            </div>
            
            <p class="message">
                To get started and unlock all features, please verify your email address by clicking the button below:
            </p>
            
            <div class="cta-container">
                <a href="%s" class="verify-button">Verify My Email</a>
            </div>
            
            <div class="link-fallback">
                <p><strong>Button not working?</strong> Copy and paste this link into your browser:</p>
                <div class="link-text">%s</div>
            </div>
            
            <div class="warning-box">
                <p><strong>⚠️ Important:</strong> This verification link expires in 24 hours for your security.</p>
            </div>
            
            <p class="message" style="font-size: 14px; color: #888;">
                Didn't create an account? You can safely ignore this email. If you have any questions, feel free to contact our support team.
            </p>
        </div>
        
        <div class="footer">
            <div class="footer-content">
                <p><strong>Kenyan Ride Share</strong></p>
                <p>Connecting Kenya, One Ride at a Time</p>
            </div>
            <div class="social-links">
                <a href="#">📱 App Store</a>
                <a href="#">🤖 Google Play</a>
                <a href="#">📧 Support</a>
            </div>
            <p class="footer-note">
                © 2024 Kenyan Ride Share. All rights reserved.<br>
                This email was sent to you because you signed up for our service.
            </p>
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
