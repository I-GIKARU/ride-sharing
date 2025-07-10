package main

import (
	"log"
	"os"

	"kenyan-ride-share-backend/internal/config"
	"kenyan-ride-share-backend/internal/handlers"
	"kenyan-ride-share-backend/internal/middleware"
	"kenyan-ride-share-backend/pkg/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Initialize handlers
	userHandler := handlers.NewUserHandler(db)
	rideHandler := handlers.NewRideHandler(db)
	paymentHandler := handlers.NewPaymentHandler(db)
	complianceHandler := handlers.NewComplianceHandler(db)

	// API routes
	api := r.Group("/api/v1")
	{
		// User management routes
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
		
		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// User routes
			protected.GET("/users/:id", userHandler.GetUser)
			protected.PUT("/users/:id", userHandler.UpdateUser)
			protected.POST("/drivers/onboard", userHandler.OnboardDriver)

			// Ride routes
			protected.POST("/ride_requests", rideHandler.CreateRideRequest)
			protected.GET("/ride_requests/:id", rideHandler.GetRideRequest)
			protected.GET("/ride_requests/nearby_drivers", rideHandler.GetNearbyDrivers)
			protected.PUT("/ride_requests/:id/accept", rideHandler.AcceptRideRequest)
			protected.PUT("/ride_requests/:id/reject", rideHandler.RejectRideRequest)
			protected.PUT("/rides/:id/start", rideHandler.StartRide)
			protected.PUT("/rides/:id/end", rideHandler.EndRide)
			protected.GET("/rides/:id", rideHandler.GetRide)
			protected.GET("/users/:id/rides", rideHandler.GetUserRides)

			// Location routes
			protected.PUT("/drivers/:id/location", rideHandler.UpdateDriverLocation)
			protected.GET("/drivers/location/:id", rideHandler.GetDriverLocation)

			// Payment routes
			protected.POST("/payments/mpesa/stk_push", paymentHandler.InitiateMpesaPayment)
			protected.GET("/payments/:id", paymentHandler.GetPayment)

			// Review routes
			protected.POST("/reviews", rideHandler.CreateReview)
			protected.GET("/users/:id/reviews", rideHandler.GetUserReviews)

			// Compliance routes (Kenya-specific)
			protected.GET("/compliance/drivers/:id/check", complianceHandler.CheckDriverCompliance)
			protected.GET("/compliance/commission/calculate", complianceHandler.CalculateCommission)
			protected.GET("/compliance/reports/ntsa", complianceHandler.GenerateNTSAReport)
			protected.POST("/compliance/vehicles/validate", complianceHandler.ValidateVehicle)
		}

		// M-Pesa callback (no auth required)
		api.POST("/payments/mpesa/callback", paymentHandler.MpesaCallback)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":      "ok",
			"service":     "Kenyan Ride Share Backend",
			"version":     "1.0.0",
			"environment": cfg.Environment,
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Kenyan Ride Share Backend starting on port %s", port)
	log.Printf("🌍 Environment: %s", cfg.Environment)
	log.Printf("📊 Health check: http://localhost:%s/health", port)
	log.Printf("📖 API Documentation: http://localhost:%s/api/v1", port)
	
	if err := r.Run("0.0.0.0:" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

