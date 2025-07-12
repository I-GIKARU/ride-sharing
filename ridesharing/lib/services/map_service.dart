import 'package:flutter/material.dart';
import 'dart:async';

class MapService {
  // This is a mock map service that would be replaced with actual map integration
  // like Google Maps or Mapbox in a production app
  
  static Future<Map<String, dynamic>> getRouteDetails(
    double startLat, 
    double startLng, 
    double endLat, 
    double endLng
  ) async {
    // Simulate network delay
    await Future.delayed(const Duration(seconds: 1));
    
    // Return mock route data
    return {
      'distance': 5.2, // in kilometers
      'duration': 15, // in minutes
      'polyline': 'mock_polyline_data',
      'fare': {
        'economy': 12.50,
        'comfort': 18.75,
        'premium': 25.30,
      },
    };
  }
  
  static Future<Map<String, dynamic>> getNearbyDrivers(double lat, double lng) async {
    // Simulate network delay
    await Future.delayed(const Duration(seconds: 1));
    
    // Return mock nearby drivers
    return {
      'drivers': [
        {
          'id': 'd1',
          'lat': lat + 0.002,
          'lng': lng - 0.003,
          'type': 'economy',
          'eta': 3, // minutes
        },
        {
          'id': 'd2',
          'lat': lat - 0.001,
          'lng': lng + 0.002,
          'type': 'comfort',
          'eta': 5, // minutes
        },
        {
          'id': 'd3',
          'lat': lat + 0.003,
          'lng': lng + 0.001,
          'type': 'premium',
          'eta': 8, // minutes
        },
      ],
    };
  }
  
  static Future<Map<String, dynamic>> geocodeAddress(String address) async {
    // Simulate network delay
    await Future.delayed(const Duration(seconds: 1));
    
    // Return mock geocoded location
    // In a real app, this would convert an address to coordinates
    return {
      'lat': 37.7749,
      'lng': -122.4194,
      'formatted_address': address,
    };
  }
  
  static Future<Map<String, dynamic>> reverseGeocode(double lat, double lng) async {
    // Simulate network delay
    await Future.delayed(const Duration(seconds: 1));
    
    // Return mock address from coordinates
    // In a real app, this would convert coordinates to an address
    return {
      'address': '123 Main St, San Francisco, CA 94105',
      'place_name': 'Downtown',
    };
  }
  
  static Widget buildMapWidget({
    required double centerLat,
    required double centerLng,
    double zoom = 15.0,
    List<Map<String, dynamic>> markers = const [],
    String? routePolyline,
  }) {
    // This is a placeholder widget that would be replaced with an actual map widget
    // like GoogleMap or MapboxMap in a production app
    return Container(
      decoration: BoxDecoration(
        color: Colors.grey.shade300,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Stack(
        children: [
          Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                const Icon(
                  Icons.map,
                  size: 48,
                  color: Colors.blue,
                ),
                const SizedBox(height: 16),
                const Text(
                  'Map View',
                  style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 8),
                Text(
                  'Center: $centerLat, $centerLng',
                  style: TextStyle(
                    color: Colors.grey.shade700,
                  ),
                ),
                if (markers.isNotEmpty)
                  Text(
                    '${markers.length} markers on map',
                    style: TextStyle(
                      color: Colors.grey.shade700,
                    ),
                  ),
                if (routePolyline != null)
                  const Text(
                    'Route displayed',
                    style: TextStyle(
                      color: Colors.blue,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
              ],
            ),
          ),
          // Markers would be added here in a real implementation
          ...markers.map((marker) => Positioned(
            left: 0,
            right: 0,
            top: 0,
            bottom: 0,
            child: Center(
              child: Icon(
                marker['type'] == 'pickup' ? Icons.my_location : 
                marker['type'] == 'destination' ? Icons.location_on :
                Icons.directions_car,
                color: marker['type'] == 'pickup' ? Colors.green : 
                       marker['type'] == 'destination' ? Colors.red :
                       Colors.blue,
              ),
            ),
          )).toList(),
        ],
      ),
    );
  }
}
