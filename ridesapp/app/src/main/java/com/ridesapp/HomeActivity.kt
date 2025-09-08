package com.ridesapp

import android.content.Intent
import android.os.Bundle
import androidx.appcompat.app.AppCompatActivity
import com.ridesapp.databinding.ActivityHomeBinding

class HomeActivity : AppCompatActivity() {
    
    private lateinit var binding: ActivityHomeBinding
    
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityHomeBinding.inflate(layoutInflater)
        setContentView(binding.root)
        
        val userName = getSharedPreferences("rides_prefs", MODE_PRIVATE)
            .getString("user_name", "User")
        
        binding.tvWelcome.text = "Welcome, $userName!"
        
        binding.btnRequestRide.setOnClickListener {
            startActivity(Intent(this, RideRequestActivity::class.java))
        }
        
        binding.btnLogout.setOnClickListener {
            getSharedPreferences("rides_prefs", MODE_PRIVATE)
                .edit()
                .clear()
                .apply()
            
            startActivity(Intent(this, LoginActivity::class.java))
            finish()
        }
    }
}
