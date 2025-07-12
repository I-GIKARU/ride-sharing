# Kenyan Ride Share Backend

A comprehensive Go-based REST API backend for ride-sharing applications, specifically designed for the Kenyan market with full NTSA compliance and M-Pesa integration.

## 🚀 Features

### Core Features
- **User Management**: Complete registration, authentication, and profile management
- **Ride Management**: Full ride lifecycle from request to completion
- **Driver Matching**: Intelligent location-based driver-passenger matching
- **Real-time Location**: GPS tracking for drivers and rides
- **Payment Integration**: Native M-Pesa STK Push integration
- **Review System**: Bidirectional rating and review system

### Kenya-Specific Features
- **NTSA Compliance**: Full compliance with Kenya's transport network regulations
- **M-Pesa Integration**: Seamless integration with Kenya's leading mobile payment platform
- **Commission Capping**: Automatic 18% commission limit as per Kenya regulations
- **Vehicle Eligibility**: Kenya-specific vehicle age and type validation
- **Phone Validation**: Kenyan phone number format validation and formatting
- **Local Currency**: KES (Kenyan Shilling) support with local fare calculation

## 🏗️ Architecture

```
kenyan-ride-share-backend/
├── cmd/                    # Application entry points
│   └── main.go            # Main application
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   ├── handlers/         # HTTP request handlers
│   ├── middleware/       # HTTP middleware
│   ├── models/          # Data models
│   └── services/        # Business logic services
├── pkg/                  # Public packages
│   ├── database/        # Database connection and utilities
│   └── utils/           # Utility functions
├── migrations/          # Database migrations
├── docs/               # Documentation
└── .env.example        # Environment configuration template
```

## 🛠️ Technology Stack

- **Language**: Go 1.21+
- **Framework**: Gin (HTTP web framework)
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT tokens
- **Payment**: Safaricom Daraja API (M-Pesa)
- **Documentation**: Markdown with comprehensive API docs

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Git
- Safaricom Daraja API credentials (for M-Pesa integration)

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone <repository-url>
cd kenyan-ride-share-backend
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Database Setup

Create a PostgreSQL database:

```sql
CREATE DATABASE kenyan_ride_share_db;
CREATE USER kenyan_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE kenyan_ride_share_db TO kenyan_user;
```

Run migrations:

```bash
psql -U kenyan_user -d kenyan_ride_share_db -f migrations/001_create_tables.sql
```

### 4. Environment Configuration

Copy and configure the environment file:

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```bash
DATABASE_URL=postgres://kenyan_user:your_password@localhost:5432/kenyan_ride_share_db?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key
MPESA_CONSUMER_KEY=your_mpesa_consumer_key
MPESA_CONSUMER_SECRET=your_mpesa_consumer_secret
MPESA_PASSKEY=your_mpesa_passkey
MPESA_SHORTCODE=your_business_shortcode
MPESA_CALLBACK_URL=http://localhost:8080/api/v1/payments/mpesa/callback
```

### 5. Run the Application

```bash
go run cmd/main.go
```

The server will start on `http://localhost:8080`

### 6. Verify Installation

```bash
curl http://localhost:8080/health
```

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication
Include JWT token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

### Key Endpoints

#### User Management
- `POST /register` - Register new user
- `POST /login` - User login
- `GET /users/{id}` - Get user profile
- `POST /drivers/onboard` - Driver onboarding

#### Ride Management
- `POST /ride_requests` - Create ride request
- `GET /ride_requests/nearby_drivers` - Find nearby drivers
- `PUT /ride_requests/{id}/accept` - Accept ride request
- `PUT /rides/{id}/start` - Start ride
- `PUT /rides/{id}/end` - End ride

#### Payments
- `POST /payments/mpesa/stk_push` - Initiate M-Pesa payment
- `POST /payments/mpesa/callback` - M-Pesa callback (webhook)

#### Compliance (Kenya-specific)
- `GET /compliance/drivers/{id}/check` - Check driver compliance
- `GET /compliance/commission/calculate` - Calculate commission
- `GET /compliance/reports/ntsa` - Generate NTSA report
- `POST /compliance/vehicles/validate` - Validate vehicle eligibility

For complete API documentation, see [API_DOCUMENTATION.md](docs/API_DOCUMENTATION.md)

## 🧪 Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## 🚀 Deployment

### Development
```bash
go run cmd/main.go
```

### Production
See the comprehensive [Deployment Guide](docs/DEPLOYMENT_GUIDE.md) for detailed production deployment instructions including:
- Traditional server deployment
- Docker containerization
- Cloud platform deployment (AWS, GCP, DigitalOcean)
- Security considerations
- Monitoring and logging
- Backup and recovery

## 🏛️ Kenya Compliance

This backend is designed to comply with Kenya's transport regulations:

### NTSA Regulations
- **Legal Notice 120 of 2022**: Full compliance with transport network company regulations
- **Driver Requirements**: License verification and approval workflow
- **Vehicle Standards**: Age and type validation (2015+ vehicles)
- **Commission Limits**: Maximum 18% service fee as per regulations
- **Reporting**: NTSA-compliant trip and revenue reporting

### M-Pesa Integration
- **Daraja API**: Official Safaricom integration
- **STK Push**: Customer-initiated payments
- **Callback Handling**: Secure payment confirmation
- **Phone Validation**: Kenya-specific number formatting

## 🔧 Configuration

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `DATABASE_URL` | PostgreSQL connection string | Yes |
| `JWT_SECRET` | Secret key for JWT tokens | Yes |
| `MPESA_CONSUMER_KEY` | M-Pesa API consumer key | Yes |
| `MPESA_CONSUMER_SECRET` | M-Pesa API consumer secret | Yes |
| `MPESA_PASSKEY` | M-Pesa API passkey | Yes |
| `MPESA_SHORTCODE` | M-Pesa business shortcode | Yes |
| `MPESA_CALLBACK_URL` | M-Pesa callback URL | Yes |
| `PORT` | Server port (default: 8080) | No |
| `ENVIRONMENT` | Environment (development/staging/production) | No |

### Database Schema

The application uses PostgreSQL with the following main tables:
- `users` - User accounts (passengers, drivers, admins)
- `drivers` - Driver-specific information and vehicle details
- `ride_requests` - Ride requests from passengers
- `rides` - Active and completed rides
- `payments` - Payment transactions
- `reviews` - User reviews and ratings

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines
- Follow Go best practices and conventions
- Write comprehensive tests for new features
- Update documentation for API changes
- Ensure Kenya compliance for transport-related features

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

For support and questions:
- Create an issue in the repository
- Check the [API Documentation](docs/API_DOCUMENTATION.md)
- Review the [Deployment Guide](docs/DEPLOYMENT_GUIDE.md)

## 🗺️ Roadmap

### Phase 1 (Current)
- ✅ Core ride-sharing functionality
- ✅ M-Pesa payment integration
- ✅ NTSA compliance features
- ✅ Basic API documentation

### Phase 2 (Planned)
- [ ] Real-time notifications (WebSocket/SSE)
- [ ] Advanced analytics and reporting
- [ ] Multi-language support (Swahili)
- [ ] Mobile app SDK

### Phase 3 (Future)
- [ ] Machine learning for demand prediction
- [ ] Integration with other payment providers
- [ ] Advanced fraud detection
- [ ] API rate limiting and throttling

## 🌍 Kenya Market Focus

This backend is specifically designed for the Kenyan ride-sharing market:

### Local Considerations
- **Currency**: All pricing in Kenyan Shillings (KES)
- **Phone Numbers**: Kenya-specific validation (+254 format)
- **Regulations**: NTSA compliance built-in
- **Payment**: M-Pesa as primary payment method
- **Geography**: Optimized for Kenyan cities and road networks

### Market Features
- **Rush Hour Pricing**: Nairobi-specific rush hour detection
- **Local Fare Calculation**: Kenya-appropriate base fares and per-km rates
- **Vehicle Standards**: Kenya vehicle age and type requirements
- **Driver Verification**: PSV license and NTSA approval workflow

---

**Built with ❤️ for the Kenyan ride-sharing market**

