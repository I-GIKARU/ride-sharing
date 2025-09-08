package com.ridesapp.models

data class User(
    val id: String,
    val firstName: String,
    val lastName: String,
    val email: String,
    val phoneNumber: String,
    val userType: String
)

data class LoginRequest(
    val email: String,
    val password: String
)

data class LoginResponse(
    val message: String,
    val token: String,
    val user: User
)

data class RideRequest(
    val pickupLatitude: Double,
    val pickupLongitude: Double,
    val pickupAddress: String,
    val dropoffLatitude: Double,
    val dropoffLongitude: Double,
    val dropoffAddress: String,
    val specialInstructions: String? = null
)

data class Driver(
    val driverId: String,
    val driverName: String,
    val vehicleInfo: String,
    val rating: Double,
    val distanceKm: Double,
    val estimatedArrivalMinutes: Int
)
