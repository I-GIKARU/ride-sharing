# Kenyan Ride Share Backend API Documentation

## Overview

The Kenyan Ride Share Backend is a comprehensive Go-based REST API designed specifically for the Kenyan ride-sharing market. It provides full compliance with Kenya's National Transport and Safety Authority (NTSA) regulations, seamless M-Pesa payment integration, and robust user management capabilities.

## Features

### Core Features
- **User Management**: Registration, authentication, and profile management for passengers and drivers
- **Ride Management**: Complete ride lifecycle from request to completion
- **Driver Matching**: Intelligent driver-passenger matching based on location and availability
- **Real-time Location Tracking**: GPS-based location services for drivers and rides
- **Payment Integration**: M-Pesa STK Push integration for seamless payments
- **Review System**: Bidirectional rating and review system

### Kenya-Specific Features
- **NTSA Compliance**: Full compliance with Kenya's transport network regulations
- **M-Pesa Integration**: Native support for Kenya's leading mobile payment platform
- **Commission Capping**: Automatic 18% commission limit as per Kenya regulations
- **Vehicle Eligibility**: Kenya-specific vehicle age and type validation
- **Phone Number Validation**: Kenyan phone number format validation and formatting
- **Local Currency**: KES (Kenyan Shilling) support with local fare calculation

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

## API Endpoints

### User Management

#### Register User
```http
POST /register
```

**Request Body:**
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "phone_number": "254712345678",
  "password": "securepassword123",
  "user_type": "passenger"
}
```

**Response:**
```json
{
  "message": "User registered successfully",
  "user_id": "123e4567-e89b-12d3-a456-426614174000"
}
```

#### Login
```http
POST /login
```

**Request Body:**
```json
{
  "email": "john.doe@example.com",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "user_type": "passenger"
  }
}
```

#### Get User Profile
```http
GET /users/{id}
```

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "phone_number": "254712345678",
  "user_type": "passenger",
  "is_verified": true,
  "created_at": "2024-01-15T10:30:00Z"
}
```

#### Driver Onboarding
```http
POST /drivers/onboard
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "driver_license_number": "DL123456789",
  "license_plate": "KCA123A",
  "vehicle_make": "Toyota",
  "vehicle_model": "Corolla",
  "vehicle_year": 2020,
  "vehicle_color": "White",
  "insurance_details": "Insurance Company XYZ, Policy: INS123456"
}
```

### Ride Management

#### Create Ride Request
```http
POST /ride_requests
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "pickup_latitude": -1.2921,
  "pickup_longitude": 36.8219,
  "pickup_address": "Nairobi CBD, Kenya",
  "dropoff_latitude": -1.3032,
  "dropoff_longitude": 36.8856,
  "dropoff_address": "Westlands, Nairobi, Kenya",
  "special_instructions": "Please call when you arrive"
}
```

**Response:**
```json
{
  "id": "456e7890-e89b-12d3-a456-426614174001",
  "passenger_id": "123e4567-e89b-12d3-a456-426614174000",
  "pickup_address": "Nairobi CBD, Kenya",
  "dropoff_address": "Westlands, Nairobi, Kenya",
  "estimated_distance_km": 8.5,
  "estimated_duration_minutes": 25,
  "estimated_fare": 425.0,
  "status": "pending",
  "created_at": "2024-01-15T11:00:00Z"
}
```

#### Get Nearby Drivers
```http
GET /ride_requests/nearby_drivers?latitude=-1.2921&longitude=36.8219&radius=5
```

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "drivers": [
    {
      "driver_id": "789e0123-e89b-12d3-a456-426614174002",
      "driver_name": "Peter Kamau",
      "vehicle_info": "White Toyota Corolla (KCA123A)",
      "rating": 4.8,
      "distance_km": 1.2,
      "estimated_arrival_minutes": 5,
      "current_latitude": -1.2901,
      "current_longitude": 36.8199
    }
  ]
}
```

#### Accept Ride Request
```http
PUT /ride_requests/{id}/accept
```

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "message": "Ride request accepted",
  "ride_id": "abc12345-e89b-12d3-a456-426614174003"
}
```

#### Start Ride
```http
PUT /rides/{id}/start
```

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "message": "Ride started",
  "start_time": "2024-01-15T11:30:00Z"
}
```

#### End Ride
```http
PUT /rides/{id}/end
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "end_latitude": -1.3032,
  "end_longitude": 36.8856,
  "actual_distance_km": 8.7,
  "actual_duration_minutes": 28
}
```

**Response:**
```json
{
  "message": "Ride completed",
  "ride_id": "abc12345-e89b-12d3-a456-426614174003",
  "final_fare": 450.0,
  "commission_breakdown": {
    "total_fare": 450.0,
    "commission_rate": 0.18,
    "commission_amount": 81.0,
    "driver_earnings": 369.0,
    "currency": "KES"
  }
}
```

### Payment Integration

#### Initiate M-Pesa Payment
```http
POST /payments/mpesa/stk_push
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "ride_id": "abc12345-e89b-12d3-a456-426614174003",
  "phone_number": "254712345678",
  "amount": 450.0
}
```

**Response:**
```json
{
  "message": "M-Pesa payment initiated",
  "payment_id": "def45678-e89b-12d3-a456-426614174004",
  "checkout_request_id": "ws_CO_15012024113000123456789",
  "customer_message": "Success. Request accepted for processing"
}
```

#### M-Pesa Callback (Webhook)
```http
POST /payments/mpesa/callback
```

**Note:** This endpoint is called by Safaricom's servers and doesn't require authentication.

### Location Services

#### Update Driver Location
```http
PUT /drivers/{id}/location
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "latitude": -1.2921,
  "longitude": 36.8219,
  "is_available": true
}
```

### Reviews

#### Create Review
```http
POST /reviews
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "ride_id": "abc12345-e89b-12d3-a456-426614174003",
  "reviewee_id": "789e0123-e89b-12d3-a456-426614174002",
  "rating": 5,
  "comment": "Excellent service, very professional driver!"
}
```

### Compliance (Kenya-Specific)

#### Check Driver Compliance
```http
GET /compliance/drivers/{id}/check
```

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "driver_id": "789e0123-e89b-12d3-a456-426614174002",
  "is_compliant": true,
  "issues": [],
  "last_checked": "2024-01-15T12:00:00Z"
}
```

#### Calculate Commission
```http
GET /compliance/commission/calculate?fare=1000
```

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "total_fare": 1000.0,
  "commission_rate": 0.18,
  "commission_amount": 180.0,
  "driver_earnings": 820.0,
  "currency": "KES"
}
```

#### Generate NTSA Report
```http
GET /compliance/reports/ntsa?start_date=2024-01-01&end_date=2024-01-31
```

**Headers:** `Authorization: Bearer <token>` (Admin only)

#### Validate Vehicle
```http
POST /compliance/vehicles/validate
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "vehicle_year": 2020,
  "vehicle_make": "Toyota",
  "vehicle_model": "Corolla"
}
```

**Response:**
```json
{
  "is_eligible": true,
  "issues": [],
  "max_age_years": 10
}
```

## Error Responses

The API returns standard HTTP status codes and error messages:

### 400 Bad Request
```json
{
  "error": "Invalid request data"
}
```

### 401 Unauthorized
```json
{
  "error": "Invalid or missing authentication token"
}
```

### 403 Forbidden
```json
{
  "error": "Access denied"
}
```

### 404 Not Found
```json
{
  "error": "Resource not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error"
}
```

## Rate Limiting

The API implements rate limiting to prevent abuse:
- **100 requests per minute** per IP address
- Rate limit headers are included in responses:
  - `X-RateLimit-Limit`: Maximum requests allowed
  - `X-RateLimit-Remaining`: Remaining requests in current window
  - `X-RateLimit-Reset`: Time when the rate limit resets

## Data Models

### User
```json
{
  "id": "uuid",
  "first_name": "string",
  "last_name": "string",
  "email": "string",
  "phone_number": "string",
  "user_type": "passenger|driver|admin",
  "is_verified": "boolean",
  "profile_picture_url": "string",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Driver
```json
{
  "driver_id": "uuid",
  "driver_license_number": "string",
  "license_plate": "string",
  "vehicle_make": "string",
  "vehicle_model": "string",
  "vehicle_year": "integer",
  "vehicle_color": "string",
  "insurance_details": "string",
  "is_approved": "boolean",
  "is_available": "boolean",
  "current_latitude": "decimal",
  "current_longitude": "decimal",
  "rating": "decimal",
  "total_rides": "integer"
}
```

### Ride Request
```json
{
  "id": "uuid",
  "passenger_id": "uuid",
  "pickup_latitude": "decimal",
  "pickup_longitude": "decimal",
  "pickup_address": "string",
  "dropoff_latitude": "decimal",
  "dropoff_longitude": "decimal",
  "dropoff_address": "string",
  "estimated_distance_km": "decimal",
  "estimated_duration_minutes": "integer",
  "estimated_fare": "decimal",
  "status": "pending|accepted|rejected|cancelled",
  "special_instructions": "string"
}
```

### Ride
```json
{
  "id": "uuid",
  "ride_request_id": "uuid",
  "driver_id": "uuid",
  "passenger_id": "uuid",
  "status": "accepted|started|completed|cancelled",
  "start_time": "timestamp",
  "end_time": "timestamp",
  "actual_distance_km": "decimal",
  "actual_duration_minutes": "integer",
  "actual_fare": "decimal"
}
```

### Payment
```json
{
  "id": "uuid",
  "ride_id": "uuid",
  "amount": "decimal",
  "currency": "string",
  "payment_method": "string",
  "transaction_id": "string",
  "payment_status": "pending|completed|failed|refunded"
}
```

## Environment Configuration

The following environment variables are required:

```bash
# Database
DATABASE_URL=postgres://username:password@localhost:5432/kenyan_ride_share_db

# JWT
JWT_SECRET=your-super-secret-jwt-key

# M-Pesa (Safaricom Daraja API)
MPESA_CONSUMER_KEY=your_consumer_key
MPESA_CONSUMER_SECRET=your_consumer_secret
MPESA_PASSKEY=your_passkey
MPESA_SHORTCODE=your_shortcode
MPESA_CALLBACK_URL=https://yourdomain.com/api/v1/payments/mpesa/callback

# Server
PORT=8080
ENVIRONMENT=development
```

## Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## Health Check

Check API health status:

```http
GET /health
```

**Response:**
```json
{
  "status": "ok",
  "service": "Kenyan Ride Share Backend",
  "version": "1.0.0",
  "environment": "development"
}
```

