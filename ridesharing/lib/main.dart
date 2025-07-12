import 'package:flutter/material.dart';
import 'package:ridesharing/screens/splash_screen.dart';
import 'package:ridesharing/screens/login_screen.dart';
import 'package:ridesharing/screens/register_screen.dart';
import 'package:ridesharing/screens/home_screen.dart';
import 'package:ridesharing/screens/profile_screen.dart';
import 'package:ridesharing/screens/ride_booking_screen.dart';
import 'package:ridesharing/screens/ride_tracking_screen.dart';
import 'package:ridesharing/screens/payment_screen.dart';
import 'package:ridesharing/screens/ride_history_screen.dart';
import 'package:ridesharing/utils/theme.dart';
import 'package:ridesharing/services/config_service.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  // Initialize configuration service
  await ConfigService.initialize();
  
  // Print configuration for debugging
  ConfigService.printConfig();
  
  // Validate configuration
  if (!ConfigService.isConfigValid()) {
    print('⚠️ Warning: Invalid configuration detected!');
  }
  
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: ConfigService.appName,
      theme: AppTheme.lightTheme,
      darkTheme: AppTheme.darkTheme,
      themeMode: ThemeMode.system,
      debugShowCheckedModeBanner: false,
      initialRoute: '/',
      routes: {
        '/': (context) => const SplashScreen(),
        '/login': (context) => const LoginScreen(),
        '/register': (context) => const RegisterScreen(),
        '/home': (context) => const HomeScreen(),
        '/profile': (context) => const ProfileScreen(),
        '/ride_booking': (context) => const RideBookingScreen(),
        '/ride_tracking': (context) => const RideTrackingScreen(),
        '/payment': (context) => const PaymentScreen(),
        '/ride_history': (context) => const RideHistoryScreen(),
      },
    );
  }
}
