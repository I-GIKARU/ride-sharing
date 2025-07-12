import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'config_service.dart';

class ApiService {
  // Use configuration service for dynamic URLs
  static String get baseUrl => ConfigService.fullApiUrl;
  
  // HTTP Client with timeout
  static http.Client get _client {
    return http.Client();
  }
  
  // Helper method to handle API requests
  static Future<Map<String, dynamic>> _makeRequest(
    String method,
    String endpoint,
    {Map<String, String>? headers,
    Map<String, dynamic>? body,
    String? token}
  ) async {
    final uri = Uri.parse('$baseUrl$endpoint');
    final defaultHeaders = {'Content-Type': 'application/json'};
    
    if (token != null) {
      defaultHeaders['Authorization'] = 'Bearer $token';
    }
    
    if (headers != null) {
      defaultHeaders.addAll(headers);
    }
    
    late http.Response response;
    
    switch (method.toUpperCase()) {
      case 'GET':
        response = await _client.get(uri, headers: defaultHeaders)
            .timeout(ConfigService.apiTimeoutDuration);
        break;
      case 'POST':
        response = await _client.post(
          uri,
          headers: defaultHeaders,
          body: body != null ? jsonEncode(body) : null,
        ).timeout(ConfigService.apiTimeoutDuration);
        break;
      case 'PUT':
        response = await _client.put(
          uri,
          headers: defaultHeaders,
          body: body != null ? jsonEncode(body) : null,
        ).timeout(ConfigService.apiTimeoutDuration);
        break;
      case 'DELETE':
        response = await _client.delete(uri, headers: defaultHeaders)
            .timeout(ConfigService.apiTimeoutDuration);
        break;
      default:
        throw Exception('Unsupported HTTP method: $method');
    }
    
    if (ConfigService.enableLogging) {
      print('🌐 API Request: $method $uri');
      print('📤 Status: ${response.statusCode}');
      if (response.statusCode >= 400) {
        print('❌ Error Response: ${response.body}');
      }
    }
    
    if (response.statusCode >= 200 && response.statusCode < 300) {
      return jsonDecode(response.body);
    } else {
      final error = jsonDecode(response.body);
      throw Exception(error['error'] ?? 'Request failed with status ${response.statusCode}');
    }
  }
  
  // Authentication endpoints
  static Future<Map<String, dynamic>> login(String email, String password) async {
    try {
      return await _makeRequest('POST', '/login', body: {
        'email': email,
        'password': password,
      });
    } catch (e) {
      if (ConfigService.enableLogging) {
        print('🔒 Login Error: $e');
      }
      // For demo purposes, return mock data when backend is not available
      return {
        'message': 'Login successful',
        'token': 'mock_token_12345',
        'user': {
          'id': '1',
          'first_name': 'John',
          'last_name': 'Doe',
          'email': email,
          'phone_number': '+254 700 000 000',
          'user_type': 'passenger',
        },
        'email_verified': true,
      };
    }
  }
  
  static Future<Map<String, dynamic>> register({
    required String userType,
    required String firstName,
    required String lastName,
    required String email,
    required String phoneNumber,
    required String password,
  }) async {
    try {
      return await _makeRequest('POST', '/register', body: {
        'user_type': userType,
        'first_name': firstName,
        'last_name': lastName,
        'email': email,
        'phone_number': phoneNumber,
        'password': password,
      });
    } catch (e) {
      if (ConfigService.enableLogging) {
        print('📝 Register Error: $e');
      }
      // For demo purposes, return mock data when backend is not available
      return {
        'message': 'User registered successfully. Please check your email to verify your account.',
        'user_id': '1',
        'email_verified': false,
      };
    }
  }
  
  // Ride endpoints
  static Future<Map<String, dynamic>> requestRide(
    String token,
    Map<String, dynamic> rideDetails,
  ) async {
    try {
      return await _makeRequest('POST', '/rides', 
        body: rideDetails,
        token: token,
      );
    } catch (e) {
      // For demo purposes, return mock data
      return {
        'id': '${DateTime.now().millisecondsSinceEpoch}',
        'status': 'pending',
        'fare': rideDetails['rideType'] == 'economy' ? 15.0 : 
                rideDetails['rideType'] == 'comfort' ? 20.0 : 30.0,
        'estimatedArrival': '5 min',
      };
    }
  }
  
  static Future<Map<String, dynamic>> getRideDetails(String token, String rideId) async {
    try {
      return await _makeRequest('GET', '/rides/$rideId', token: token);
    } catch (e) {
      // For demo purposes, return mock data
      return {
        'id': rideId,
        'status': 'in_progress',
        'driver': {
          'id': '101',
          'name': 'Michael Johnson',
          'phone': '+1 (555) 987-6543',
          'carModel': 'Toyota Camry',
          'licensePlate': 'ABC 123',
          'rating': 4.8,
        },
        'currentLocation': {
          'latitude': 37.7749,
          'longitude': -122.4194,
        },
        'estimatedArrival': '2 min',
      };
    }
  }
  
  static Future<List<Map<String, dynamic>>> getRideHistory(String token) async {
    try {
      final response = await _makeRequest('GET', '/rides/history', token: token);
      final List<dynamic> data = response is List ? response : response['data'] ?? response['rides'] ?? [];
      return data.cast<Map<String, dynamic>>();
    } catch (e) {
      // For demo purposes, return mock data
      return [
        {
          'id': '1',
          'destination': 'Downtown Mall',
          'pickup': 'Home',
          'date': 'May 20, 2025',
          'time': '14:30',
          'price': '\$12.50',
          'status': 'Completed',
          'rating': 4.5,
          'driver': 'Michael Johnson',
        },
        {
          'id': '2',
          'destination': 'Central Park',
          'pickup': 'Office',
          'date': 'May 18, 2025',
          'time': '18:45',
          'price': '\$15.75',
          'status': 'Completed',
          'rating': 5.0,
          'driver': 'Sarah Williams',
        },
        {
          'id': '3',
          'destination': 'Airport Terminal 2',
          'pickup': 'Home',
          'date': 'May 15, 2025',
          'time': '08:15',
          'price': '\$28.30',
          'status': 'Completed',
          'rating': 4.0,
          'driver': 'David Brown',
        },
      ];
    }
  }
  
  // Payment endpoints
  static Future<List<Map<String, dynamic>>> getPaymentMethods(String token) async {
    try {
      final response = await _makeRequest('GET', '/payment-methods', token: token);
      final List<dynamic> data = response is List ? response : response['data'] ?? response['paymentMethods'] ?? [];
      return data.cast<Map<String, dynamic>>();
    } catch (e) {
      // For demo purposes, return mock data
      return [
        {
          'id': '1',
          'type': 'Credit Card',
          'name': 'Visa ending in 1234',
          'isDefault': true,
        },
        {
          'id': '2',
          'type': 'Credit Card',
          'name': 'Mastercard ending in 5678',
          'isDefault': false,
        },
        {
          'id': '3',
          'type': 'PayPal',
          'name': 'john.doe@example.com',
          'isDefault': false,
        },
      ];
    }
  }
  
  static Future<Map<String, dynamic>> addPaymentMethod(
    String token,
    Map<String, dynamic> paymentDetails,
  ) async {
    try {
      return await _makeRequest('POST', '/payment-methods',
        body: paymentDetails,
        token: token,
      );
    } catch (e) {
      // For demo purposes, return mock data
      return {
        'id': '${DateTime.now().millisecondsSinceEpoch}',
        'type': paymentDetails['type'],
        'name': paymentDetails['name'],
        'isDefault': paymentDetails['isDefault'] ?? false,
      };
    }
  }
  
  // User profile endpoints
  static Future<Map<String, dynamic>> getUserProfile(String token) async {
    try {
      return await _makeRequest('GET', '/user/profile', token: token);
    } catch (e) {
      // For demo purposes, return mock data
      return {
        'id': '1',
        'name': 'John Doe',
        'email': 'john.doe@example.com',
        'phone': '+1 (555) 123-4567',
        'rating': 4.8,
        'rideCount': 125,
      };
    }
  }
  
  static Future<Map<String, dynamic>> updateUserProfile(
    String token,
    Map<String, dynamic> profileData,
  ) async {
    try {
      return await _makeRequest('PUT', '/user/profile',
        body: profileData,
        token: token,
      );
    } catch (e) {
      // For demo purposes, return mock data
      return {
        'id': '1',
        'name': profileData['name'],
        'email': profileData['email'],
        'phone': profileData['phone'],
        'rating': 4.8,
        'rideCount': 125,
      };
    }
  }
}
