package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"kenyan-ride-share-backend/internal/models"

	"gorm.io/gorm"
)

type MpesaService struct {
	db           *gorm.DB
	consumerKey  string
	consumerSecret string
	passkey      string
	shortcode    string
	callbackURL  string
	environment  string // "sandbox" or "production"
}

func NewMpesaService(db *gorm.DB) *MpesaService {
	return &MpesaService{
		db:             db,
		consumerKey:    os.Getenv("MPESA_CONSUMER_KEY"),
		consumerSecret: os.Getenv("MPESA_CONSUMER_SECRET"),
		passkey:        os.Getenv("MPESA_PASSKEY"),
		shortcode:      os.Getenv("MPESA_SHORTCODE"),
		callbackURL:    os.Getenv("MPESA_CALLBACK_URL"),
		environment:    os.Getenv("ENVIRONMENT"),
	}
}

type MpesaAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
}

type MpesaSTKPushRequest struct {
	BusinessShortCode string `json:"BusinessShortCode"`
	Password          string `json:"Password"`
	Timestamp         string `json:"Timestamp"`
	TransactionType   string `json:"TransactionType"`
	Amount            int    `json:"Amount"`
	PartyA            string `json:"PartyA"`
	PartyB            string `json:"PartyB"`
	PhoneNumber       string `json:"PhoneNumber"`
	CallBackURL       string `json:"CallBackURL"`
	AccountReference  string `json:"AccountReference"`
	TransactionDesc   string `json:"TransactionDesc"`
}

type MpesaSTKPushResponse struct {
	MerchantRequestID   string `json:"MerchantRequestID"`
	CheckoutRequestID   string `json:"CheckoutRequestID"`
	ResponseCode        string `json:"ResponseCode"`
	ResponseDescription string `json:"ResponseDescription"`
	CustomerMessage     string `json:"CustomerMessage"`
}

func (m *MpesaService) GetAccessToken() (string, error) {
	if m.environment == "development" {
		// Return mock token for development
		return "mock_access_token_" + fmt.Sprintf("%d", time.Now().Unix()), nil
	}

	// Production implementation
	auth := base64.StdEncoding.EncodeToString([]byte(m.consumerKey + ":" + m.consumerSecret))
	
	baseURL := "https://api.safaricom.co.ke"
	if m.environment == "sandbox" {
		baseURL = "https://sandbox.safaricom.co.ke"
	}

	req, err := http.NewRequest("GET", baseURL+"/oauth/v1/generate?grant_type=client_credentials", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenResp MpesaAccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

func (m *MpesaService) InitiateSTKPush(phoneNumber string, amount float64, accountReference string) (*MpesaSTKPushResponse, error) {
	if m.environment == "development" {
		// Return mock response for development
		return &MpesaSTKPushResponse{
			MerchantRequestID:   "mock_merchant_request_" + fmt.Sprintf("%d", time.Now().Unix()),
			CheckoutRequestID:   "mock_checkout_request_" + fmt.Sprintf("%d", time.Now().Unix()),
			ResponseCode:        "0",
			ResponseDescription: "Success. Request accepted for processing",
			CustomerMessage:     "Success. Request accepted for processing",
		}, nil
	}

	accessToken, err := m.GetAccessToken()
	if err != nil {
		return nil, err
	}

	timestamp := time.Now().Format("20060102150405")
	password := base64.StdEncoding.EncodeToString([]byte(m.shortcode + m.passkey + timestamp))

	// Format phone number (ensure it starts with 254)
	if strings.HasPrefix(phoneNumber, "0") {
		phoneNumber = "254" + phoneNumber[1:]
	} else if strings.HasPrefix(phoneNumber, "+254") {
		phoneNumber = phoneNumber[1:]
	} else if !strings.HasPrefix(phoneNumber, "254") {
		phoneNumber = "254" + phoneNumber
	}

	payload := MpesaSTKPushRequest{
		BusinessShortCode: m.shortcode,
		Password:          password,
		Timestamp:         timestamp,
		TransactionType:   "CustomerPayBillOnline",
		Amount:            int(amount),
		PartyA:            phoneNumber,
		PartyB:            m.shortcode,
		PhoneNumber:       phoneNumber,
		CallBackURL:       m.callbackURL,
		AccountReference:  accountReference,
		TransactionDesc:   "Ride payment",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	baseURL := "https://api.safaricom.co.ke"
	if m.environment == "sandbox" {
		baseURL = "https://sandbox.safaricom.co.ke"
	}

	req, err := http.NewRequest("POST", baseURL+"/mpesa/stkpush/v1/processrequest", strings.NewReader(string(jsonPayload)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var stkResp MpesaSTKPushResponse
	if err := json.NewDecoder(resp.Body).Decode(&stkResp); err != nil {
		return nil, err
	}

	return &stkResp, nil
}

func (m *MpesaService) ProcessCallback(callbackData map[string]interface{}) error {
	// Extract callback information
	body, ok := callbackData["Body"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid callback format")
	}

	stkCallback, ok := body["stkCallback"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid stk callback format")
	}

	checkoutRequestID, _ := stkCallback["CheckoutRequestID"].(string)
	resultCode, _ := stkCallback["ResultCode"].(float64)

	// Find payment by checkout request ID
	var payment models.Payment
	if err := m.db.Where("transaction_id = ?", checkoutRequestID).First(&payment).Error; err != nil {
		return fmt.Errorf("payment not found: %v", err)
	}

	// Update payment status
	if resultCode == 0 {
		payment.PaymentStatus = "completed"
		
		// Extract M-Pesa receipt number if available
		if callbackMetadata, ok := stkCallback["CallbackMetadata"].(map[string]interface{}); ok {
			if items, ok := callbackMetadata["Item"].([]interface{}); ok {
				for _, item := range items {
					if itemMap, ok := item.(map[string]interface{}); ok {
						if name, ok := itemMap["Name"].(string); ok && name == "MpesaReceiptNumber" {
							if receiptNumber, ok := itemMap["Value"].(string); ok {
								payment.TransactionID = &receiptNumber
							}
						}
					}
				}
			}
		}
	} else {
		payment.PaymentStatus = "failed"
	}

	return m.db.Save(&payment).Error
}

