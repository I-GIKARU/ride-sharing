package com.ridesapp.api

import com.ridesapp.models.*
import retrofit2.Response
import retrofit2.http.*

interface ApiService {
    
    @POST("login")
    suspend fun login(@Body request: LoginRequest): Response<LoginResponse>
    
    @POST("ride_requests")
    suspend fun createRideRequest(@Body request: RideRequest): Response<Any>
    
    @GET("ride_requests/nearby_drivers")
    suspend fun getNearbyDrivers(
        @Query("latitude") latitude: Double,
        @Query("longitude") longitude: Double,
        @Query("radius") radius: Int = 5
    ): Response<List<Driver>>
}

object ApiClient {
    private const val BASE_URL = "http://10.0.2.2:8080/api/v1/"
    
    val apiService: ApiService by lazy {
        retrofit2.Retrofit.Builder()
            .baseUrl(BASE_URL)
            .addConverterFactory(retrofit2.converter.gson.GsonConverterFactory.create())
            .build()
            .create(ApiService::class.java)
    }
}
