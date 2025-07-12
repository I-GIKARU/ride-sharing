import 'package:flutter/material.dart';

class User {
  final String id;
  final String name;
  final String email;
  final String phone;
  final String profilePicture;
  final double rating;
  
  User({
    required this.id,
    required this.name,
    required this.email,
    required this.phone,
    this.profilePicture = '',
    this.rating = 0.0,
  });
}

class Driver {
  final String id;
  final String name;
  final String phone;
  final String carModel;
  final String licensePlate;
  final double rating;
  final String profilePicture;
  
  Driver({
    required this.id,
    required this.name,
    required this.phone,
    required this.carModel,
    required this.licensePlate,
    this.rating = 0.0,
    this.profilePicture = '',
  });
}

class Location {
  final String address;
  final double latitude;
  final double longitude;
  final String name;
  
  Location({
    required this.address,
    required this.latitude,
    required this.longitude,
    this.name = '',
  });
}

class Ride {
  final String id;
  final Location pickup;
  final Location destination;
  final DateTime requestTime;
  final DateTime? pickupTime;
  final DateTime? dropoffTime;
  final String status; // pending, accepted, in_progress, completed, cancelled
  final double fare;
  final String paymentMethod;
  final bool isPaid;
  final User user;
  final Driver? driver;
  final String rideType; // economy, comfort, premium
  final double? rating;
  
  Ride({
    required this.id,
    required this.pickup,
    required this.destination,
    required this.requestTime,
    this.pickupTime,
    this.dropoffTime,
    required this.status,
    required this.fare,
    required this.paymentMethod,
    this.isPaid = false,
    required this.user,
    this.driver,
    required this.rideType,
    this.rating,
  });
}

class PaymentMethod {
  final String id;
  final String type; // credit_card, paypal, etc.
  final String name; // Visa ending in 1234
  final bool isDefault;
  
  PaymentMethod({
    required this.id,
    required this.type,
    required this.name,
    this.isDefault = false,
  });
}
