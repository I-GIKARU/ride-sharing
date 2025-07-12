package utils

import (
	"time"
)

// TimeNow returns current time as string
func TimeNow() string {
	return time.Now().Format(time.RFC3339)
}

// CurrentYear returns current year
func CurrentYear() int {
	return time.Now().Year()
}

// FormatKenyanPhoneNumber formats phone number to Kenya format
func FormatKenyanPhoneNumber(phoneNumber string) string {
	// Remove any spaces or special characters
	cleaned := ""
	for _, char := range phoneNumber {
		if char >= '0' && char <= '9' {
			cleaned += string(char)
		}
	}

	// Convert to Kenya format (254XXXXXXXXX)
	if len(cleaned) == 10 && cleaned[0] == '0' {
		// 0712345678 -> 254712345678
		return "254" + cleaned[1:]
	} else if len(cleaned) == 9 {
		// 712345678 -> 254712345678
		return "254" + cleaned
	} else if len(cleaned) == 12 && cleaned[:3] == "254" {
		// Already in correct format
		return cleaned
	} else if len(cleaned) == 13 && cleaned[0] == '+' && cleaned[1:4] == "254" {
		// +254712345678 -> 254712345678
		return cleaned[1:]
	}

	// Return as is if format is unclear
	return phoneNumber
}

// ValidateKenyanPhoneNumber checks if phone number is valid Kenya format
func ValidateKenyanPhoneNumber(phoneNumber string) bool {
	formatted := FormatKenyanPhoneNumber(phoneNumber)
	
	// Should be 12 digits starting with 254
	if len(formatted) != 12 {
		return false
	}
	
	if formatted[:3] != "254" {
		return false
	}
	
	// Check if the next digit is valid (7, 1, or 0 for Kenya mobile networks)
	fourthDigit := formatted[3]
	return fourthDigit == '7' || fourthDigit == '1' || fourthDigit == '0'
}

// CalculateFareKenyan calculates fare using Kenya-specific rates
func CalculateFareKenyan(distanceKm float64, durationMinutes int, isRushHour bool) float64 {
	baseFare := 50.0  // KES 50 base fare
	perKmRate := 25.0 // KES 25 per km
	perMinuteRate := 2.0 // KES 2 per minute
	
	fare := baseFare + (distanceKm * perKmRate) + (float64(durationMinutes) * perMinuteRate)
	
	// Apply surge pricing during rush hours
	if isRushHour {
		fare *= 1.5 // 50% surge
	}
	
	// Minimum fare
	if fare < 100.0 {
		fare = 100.0
	}
	
	return fare
}

// IsRushHour determines if current time is rush hour in Kenya
func IsRushHour() bool {
	now := time.Now()
	hour := now.Hour()
	
	// Morning rush: 7-9 AM, Evening rush: 5-7 PM
	return (hour >= 7 && hour <= 9) || (hour >= 17 && hour <= 19)
}

// GetKenyanTimezone returns Kenya timezone
func GetKenyanTimezone() *time.Location {
	loc, _ := time.LoadLocation("Africa/Nairobi")
	return loc
}

// ConvertToKenyanTime converts time to Kenya timezone
func ConvertToKenyanTime(t time.Time) time.Time {
	return t.In(GetKenyanTimezone())
}

