# Build stage
FROM golang:1.23-alpine AS builder

# Install git and ca-certificates (needed for fetching dependencies)
RUN apk update && apk add --no-cache git ca-certificates tzdata

# Create appuser
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o kenyan-ride-share cmd/main.go

# Final stage
FROM scratch

# Import ca-certificates, timezone data, and user
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd

# Copy the binary from builder
COPY --from=builder /app/kenyan-ride-share /kenyan-ride-share

# Copy migrations
COPY --from=builder /app/migrations /migrations

# Use an unprivileged user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD ["/kenyan-ride-share", "--health-check"] || exit 1

# Run the binary
ENTRYPOINT ["/kenyan-ride-share"]
