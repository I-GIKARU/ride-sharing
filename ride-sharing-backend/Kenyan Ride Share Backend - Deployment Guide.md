# Kenyan Ride Share Backend - Deployment Guide

## Overview

This guide provides comprehensive instructions for deploying the Kenyan Ride Share Backend in various environments, from local development to production deployment on cloud platforms.

## Prerequisites

### System Requirements
- **Go**: Version 1.21 or higher
- **PostgreSQL**: Version 12 or higher
- **Redis**: Version 6 or higher (optional, for caching)
- **Git**: For version control
- **Docker**: For containerized deployment (optional)

### External Services
- **Safaricom Daraja API**: M-Pesa integration credentials
- **Google Maps API**: For location services (optional)
- **SMTP Server**: For email notifications (optional)

## Local Development Setup

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

#### Install PostgreSQL
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install postgresql postgresql-contrib

# macOS
brew install postgresql
brew services start postgresql

# Windows
# Download and install from https://www.postgresql.org/download/windows/
```

#### Create Database
```bash
sudo -u postgres psql
CREATE DATABASE kenyan_ride_share_db;
CREATE USER kenyan_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE kenyan_ride_share_db TO kenyan_user;
\q
```

#### Run Migrations
```bash
psql -U kenyan_user -d kenyan_ride_share_db -f migrations/001_create_tables.sql
```

### 4. Environment Configuration

Copy the example environment file and configure it:

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```bash
# Database Configuration
DATABASE_URL=postgres://kenyan_user:your_password@localhost:5432/kenyan_ride_share_db?sslmode=disable

# JWT Secret (generate a secure random string)
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Server Configuration
ENVIRONMENT=development
PORT=8080

# M-Pesa Configuration (get from Safaricom Daraja Portal)
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

Test the health endpoint:

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "ok",
  "service": "Kenyan Ride Share Backend",
  "version": "1.0.0",
  "environment": "development"
}
```

## Production Deployment

### Option 1: Traditional Server Deployment

#### 1. Server Setup (Ubuntu 20.04/22.04)

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install required packages
sudo apt install -y git postgresql postgresql-contrib nginx certbot python3-certbot-nginx

# Install Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

#### 2. Database Setup

```bash
# Configure PostgreSQL
sudo -u postgres psql
CREATE DATABASE kenyan_ride_share_prod;
CREATE USER kenyan_prod_user WITH PASSWORD 'secure_production_password';
GRANT ALL PRIVILEGES ON DATABASE kenyan_ride_share_prod TO kenyan_prod_user;
\q

# Run migrations
psql -U kenyan_prod_user -d kenyan_ride_share_prod -f migrations/001_create_tables.sql
```

#### 3. Application Deployment

```bash
# Clone repository
git clone <repository-url> /opt/kenyan-ride-share
cd /opt/kenyan-ride-share

# Build application
go build -o kenyan-ride-share cmd/main.go

# Create production environment file
sudo cp .env.example .env
sudo nano .env  # Configure with production values
```

#### 4. Systemd Service

Create a systemd service file:

```bash
sudo nano /etc/systemd/system/kenyan-ride-share.service
```

```ini
[Unit]
Description=Kenyan Ride Share Backend
After=network.target postgresql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/kenyan-ride-share
ExecStart=/opt/kenyan-ride-share/kenyan-ride-share
Restart=always
RestartSec=5
Environment=PATH=/usr/local/go/bin:/usr/bin:/bin
EnvironmentFile=/opt/kenyan-ride-share/.env

[Install]
WantedBy=multi-user.target
```

Enable and start the service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable kenyan-ride-share
sudo systemctl start kenyan-ride-share
sudo systemctl status kenyan-ride-share
```

#### 5. Nginx Configuration

Create Nginx configuration:

```bash
sudo nano /etc/nginx/sites-available/kenyan-ride-share
```

```nginx
server {
    listen 80;
    server_name yourdomain.com www.yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Enable the site:

```bash
sudo ln -s /etc/nginx/sites-available/kenyan-ride-share /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

#### 6. SSL Certificate

```bash
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com
```

### Option 2: Docker Deployment

#### 1. Create Dockerfile

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o kenyan-ride-share cmd/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/kenyan-ride-share .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./kenyan-ride-share"]
```

#### 2. Create Docker Compose

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://kenyan_user:password@db:5432/kenyan_ride_share_db?sslmode=disable
      - JWT_SECRET=your-super-secret-jwt-key
      - ENVIRONMENT=production
      - MPESA_CONSUMER_KEY=${MPESA_CONSUMER_KEY}
      - MPESA_CONSUMER_SECRET=${MPESA_CONSUMER_SECRET}
      - MPESA_PASSKEY=${MPESA_PASSKEY}
      - MPESA_SHORTCODE=${MPESA_SHORTCODE}
      - MPESA_CALLBACK_URL=${MPESA_CALLBACK_URL}
    depends_on:
      - db
      - redis

  db:
    image: postgres:14
    environment:
      - POSTGRES_DB=kenyan_ride_share_db
      - POSTGRES_USER=kenyan_user
      - POSTGRES_PASSWORD=password
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d
    ports:
      - "5432:5432"

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - app

volumes:
  postgres_data:
```

#### 3. Deploy with Docker

```bash
# Build and start services
docker-compose up -d

# View logs
docker-compose logs -f app

# Scale the application
docker-compose up -d --scale app=3
```

### Option 3: Cloud Platform Deployment

#### AWS Deployment

1. **EC2 Instance Setup**
   - Launch Ubuntu 20.04 LTS instance
   - Configure security groups (ports 22, 80, 443, 8080)
   - Follow traditional server deployment steps

2. **RDS Database**
   - Create PostgreSQL RDS instance
   - Configure security groups for database access
   - Update DATABASE_URL in environment

3. **Application Load Balancer**
   - Create ALB for high availability
   - Configure target groups
   - Set up health checks

4. **Auto Scaling**
   - Create launch template
   - Configure auto scaling group
   - Set up CloudWatch monitoring

#### Google Cloud Platform

1. **Compute Engine**
   - Create VM instance with Ubuntu
   - Follow traditional deployment steps

2. **Cloud SQL**
   - Create PostgreSQL instance
   - Configure connection and security

3. **Load Balancer**
   - Set up HTTP(S) load balancer
   - Configure backend services

#### DigitalOcean

1. **Droplet Creation**
   - Create Ubuntu droplet
   - Follow traditional deployment steps

2. **Managed Database**
   - Create PostgreSQL cluster
   - Configure connection strings

3. **Load Balancer**
   - Set up load balancer
   - Configure health checks

## Environment-Specific Configurations

### Development
```bash
ENVIRONMENT=development
LOG_LEVEL=debug
MPESA_CALLBACK_URL=http://localhost:8080/api/v1/payments/mpesa/callback
```

### Staging
```bash
ENVIRONMENT=staging
LOG_LEVEL=info
MPESA_CALLBACK_URL=https://staging.yourdomain.com/api/v1/payments/mpesa/callback
```

### Production
```bash
ENVIRONMENT=production
LOG_LEVEL=warn
MPESA_CALLBACK_URL=https://yourdomain.com/api/v1/payments/mpesa/callback
```

## Security Considerations

### 1. Environment Variables
- Never commit `.env` files to version control
- Use secure, randomly generated JWT secrets
- Rotate secrets regularly

### 2. Database Security
- Use strong passwords
- Enable SSL connections in production
- Restrict database access to application servers only
- Regular backups and security updates

### 3. API Security
- Implement rate limiting
- Use HTTPS in production
- Validate all input data
- Implement proper authentication and authorization

### 4. M-Pesa Security
- Secure callback URLs with HTTPS
- Validate callback signatures
- Store credentials securely
- Monitor for suspicious transactions

## Monitoring and Logging

### 1. Application Logs
```bash
# View systemd logs
sudo journalctl -u kenyan-ride-share -f

# View Docker logs
docker-compose logs -f app
```

### 2. Database Monitoring
```bash
# PostgreSQL logs
sudo tail -f /var/log/postgresql/postgresql-14-main.log

# Connection monitoring
sudo -u postgres psql -c "SELECT * FROM pg_stat_activity;"
```

### 3. Performance Monitoring
- Set up Prometheus and Grafana
- Monitor API response times
- Track database performance
- Monitor M-Pesa transaction success rates

## Backup and Recovery

### 1. Database Backup
```bash
# Create backup
pg_dump -U kenyan_user -h localhost kenyan_ride_share_db > backup_$(date +%Y%m%d_%H%M%S).sql

# Automated backup script
#!/bin/bash
BACKUP_DIR="/opt/backups"
DATE=$(date +%Y%m%d_%H%M%S)
pg_dump -U kenyan_user -h localhost kenyan_ride_share_db > $BACKUP_DIR/backup_$DATE.sql
find $BACKUP_DIR -name "backup_*.sql" -mtime +7 -delete
```

### 2. Application Backup
```bash
# Backup application files
tar -czf app_backup_$(date +%Y%m%d_%H%M%S).tar.gz /opt/kenyan-ride-share
```

### 3. Recovery Procedures
```bash
# Restore database
psql -U kenyan_user -d kenyan_ride_share_db < backup_20240115_120000.sql

# Restore application
tar -xzf app_backup_20240115_120000.tar.gz -C /
sudo systemctl restart kenyan-ride-share
```

## Troubleshooting

### Common Issues

1. **Database Connection Errors**
   - Check PostgreSQL service status
   - Verify connection string
   - Check firewall settings

2. **M-Pesa Integration Issues**
   - Verify Daraja API credentials
   - Check callback URL accessibility
   - Monitor Safaricom API status

3. **High Memory Usage**
   - Monitor Go garbage collection
   - Check for memory leaks
   - Optimize database queries

4. **Performance Issues**
   - Enable database query logging
   - Monitor API response times
   - Check system resources

### Debug Commands

```bash
# Check service status
sudo systemctl status kenyan-ride-share

# View recent logs
sudo journalctl -u kenyan-ride-share --since "1 hour ago"

# Test database connection
psql -U kenyan_user -d kenyan_ride_share_db -c "SELECT version();"

# Test API endpoints
curl -X GET http://localhost:8080/health
curl -X POST http://localhost:8080/api/v1/login -H "Content-Type: application/json" -d '{"email":"test@example.com","password":"password"}'
```

## Scaling Considerations

### Horizontal Scaling
- Use load balancers
- Implement session management
- Database connection pooling
- Stateless application design

### Vertical Scaling
- Monitor CPU and memory usage
- Optimize database queries
- Implement caching strategies
- Use connection pooling

### Database Scaling
- Read replicas for read-heavy workloads
- Database sharding for large datasets
- Connection pooling
- Query optimization

## Maintenance

### Regular Tasks
- Update dependencies monthly
- Security patches weekly
- Database maintenance weekly
- Log rotation daily
- Backup verification weekly

### Update Procedure
```bash
# Stop service
sudo systemctl stop kenyan-ride-share

# Backup current version
cp /opt/kenyan-ride-share/kenyan-ride-share /opt/kenyan-ride-share/kenyan-ride-share.backup

# Update code
git pull origin main
go build -o kenyan-ride-share cmd/main.go

# Run migrations if needed
psql -U kenyan_user -d kenyan_ride_share_db -f migrations/new_migration.sql

# Start service
sudo systemctl start kenyan-ride-share

# Verify deployment
curl http://localhost:8080/health
```

This deployment guide provides comprehensive instructions for setting up the Kenyan Ride Share Backend in various environments, ensuring reliable and secure operation in production.

