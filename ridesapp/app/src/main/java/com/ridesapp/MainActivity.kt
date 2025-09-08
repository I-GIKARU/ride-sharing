package com.ridesapp

import android.content.Intent
import android.os.Bundle
import androidx.appcompat.app.AppCompatActivity
import com.ridesapp.databinding.ActivityMainBinding

class MainActivity : AppCompatActivity() {
    
    private lateinit var binding: ActivityMainBinding
    
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)
        
        // Check if user is logged in (simplified)
        val isLoggedIn = getSharedPreferences("rides_prefs", MODE_PRIVATE)
            .getBoolean("is_logged_in", false)
        
        if (isLoggedIn) {
            startActivity(Intent(this, HomeActivity::class.java))
        } else {
            startActivity(Intent(this, LoginActivity::class.java))
        }
        finish()
    }
}
