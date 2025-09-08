package com.ridesapp

import android.os.Bundle
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.lifecycleScope
import com.ridesapp.api.ApiClient
import com.ridesapp.databinding.ActivityRideRequestBinding
import com.ridesapp.models.RideRequest
import kotlinx.coroutines.launch

class RideRequestActivity : AppCompatActivity() {
    
    private lateinit var binding: ActivityRideRequestBinding
    
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityRideRequestBinding.inflate(layoutInflater)
        setContentView(binding.root)
        
        binding.btnFindDrivers.setOnClickListener {
            val pickup = binding.etPickup.text.toString()
            val destination = binding.etDestination.text.toString()
            
            if (pickup.isNotEmpty() && destination.isNotEmpty()) {
                requestRide(pickup, destination)
            } else {
                Toast.makeText(this, "Please fill pickup and destination", Toast.LENGTH_SHORT).show()
            }
        }
        
        binding.btnBack.setOnClickListener {
            finish()
        }
    }
    
    private fun requestRide(pickup: String, destination: String) {
        lifecycleScope.launch {
            try {
                val rideRequest = RideRequest(
                    pickupLatitude = -1.2921,
                    pickupLongitude = 36.8219,
                    pickupAddress = pickup,
                    dropoffLatitude = -1.3032,
                    dropoffLongitude = 36.8856,
                    dropoffAddress = destination
                )
                
                val response = ApiClient.apiService.createRideRequest(rideRequest)
                if (response.isSuccessful) {
                    Toast.makeText(this@RideRequestActivity, "Ride requested successfully!", Toast.LENGTH_SHORT).show()
                    getNearbyDrivers()
                } else {
                    Toast.makeText(this@RideRequestActivity, "Failed to request ride", Toast.LENGTH_SHORT).show()
                }
            } catch (e: Exception) {
                Toast.makeText(this@RideRequestActivity, "Network error: ${e.message}", Toast.LENGTH_SHORT).show()
            }
        }
    }
    
    private fun getNearbyDrivers() {
        lifecycleScope.launch {
            try {
                val response = ApiClient.apiService.getNearbyDrivers(-1.2921, 36.8219)
                if (response.isSuccessful) {
                    val drivers = response.body() ?: emptyList()
                    binding.tvDrivers.text = "Found ${drivers.size} nearby drivers"
                }
            } catch (e: Exception) {
                binding.tvDrivers.text = "Error finding drivers: ${e.message}"
            }
        }
    }
}
