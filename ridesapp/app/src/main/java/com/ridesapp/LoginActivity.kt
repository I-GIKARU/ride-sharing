package com.ridesapp

import android.content.Intent
import android.os.Bundle
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.lifecycleScope
import com.ridesapp.api.ApiClient
import com.ridesapp.databinding.ActivityLoginBinding
import com.ridesapp.models.LoginRequest
import kotlinx.coroutines.launch

class LoginActivity : AppCompatActivity() {
    
    private lateinit var binding: ActivityLoginBinding
    
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityLoginBinding.inflate(layoutInflater)
        setContentView(binding.root)
        
        binding.btnLogin.setOnClickListener {
            val email = binding.etEmail.text.toString()
            val password = binding.etPassword.text.toString()
            
            if (email.isNotEmpty() && password.isNotEmpty()) {
                login(email, password)
            } else {
                Toast.makeText(this, "Please fill all fields", Toast.LENGTH_SHORT).show()
            }
        }
    }
    
    private fun login(email: String, password: String) {
        lifecycleScope.launch {
            try {
                val response = ApiClient.apiService.login(LoginRequest(email, password))
                if (response.isSuccessful) {
                    response.body()?.let { loginResponse ->
                        // Save login state
                        getSharedPreferences("rides_prefs", MODE_PRIVATE)
                            .edit()
                            .putBoolean("is_logged_in", true)
                            .putString("user_name", loginResponse.user.firstName)
                            .apply()
                        
                        startActivity(Intent(this@LoginActivity, HomeActivity::class.java))
                        finish()
                    }
                } else {
                    Toast.makeText(this@LoginActivity, "Login failed", Toast.LENGTH_SHORT).show()
                }
            } catch (e: Exception) {
                Toast.makeText(this@LoginActivity, "Network error: ${e.message}", Toast.LENGTH_SHORT).show()
            }
        }
    }
}
