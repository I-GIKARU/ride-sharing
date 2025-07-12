import 'package:flutter_dotenv/flutter_dotenv.dart';

class ConfigService {
  static bool _isInitialized = false;

  // Initialize the configuration service
  static Future<void> initialize() async {
    if (_isInitialized) return;
    
    await dotenv.load(fileName: ".env");
    _isInitialized = true;
  }

  // API Configuration
  static String get apiBaseUrl => dotenv.env['API_BASE_URL'] ?? 'http://localhost:8080';
  static String get apiVersion => dotenv.env['API_VERSION'] ?? '/api/v1';
  static String get fullApiUrl => '$apiBaseUrl$apiVersion';

  // App Information
  static String get appName => dotenv.env['APP_NAME'] ?? 'Kenyan Ride Share';
  static String get appVersion => dotenv.env['APP_VERSION'] ?? '1.0.0';

  // Environment
  static String get environment => dotenv.env['ENVIRONMENT'] ?? 'development';
  static bool get isProduction => environment == 'production';
  static bool get isDevelopment => environment == 'development';

  // Debug Configuration
  static bool get debugMode => dotenv.env['DEBUG_MODE']?.toLowerCase() == 'true';
  static bool get enableLogging => dotenv.env['ENABLE_LOGGING']?.toLowerCase() == 'true';

  // API Timeouts
  static int get apiTimeout => int.tryParse(dotenv.env['API_TIMEOUT'] ?? '30') ?? 30;
  static int get connectionTimeout => int.tryParse(dotenv.env['CONNECTION_TIMEOUT'] ?? '10') ?? 10;

  // Cache Configuration
  static int get cacheDuration => int.tryParse(dotenv.env['CACHE_DURATION'] ?? '300') ?? 300;

  // Map Configuration
  static double get defaultLatitude => double.tryParse(dotenv.env['DEFAULT_LATITUDE'] ?? '-1.286389') ?? -1.286389;
  static double get defaultLongitude => double.tryParse(dotenv.env['DEFAULT_LONGITUDE'] ?? '36.817223') ?? 36.817223;
  static double get defaultZoom => double.tryParse(dotenv.env['DEFAULT_ZOOM'] ?? '15.0') ?? 15.0;

  // Helper methods
  static Duration get apiTimeoutDuration => Duration(seconds: apiTimeout);
  static Duration get connectionTimeoutDuration => Duration(seconds: connectionTimeout);
  static Duration get cacheDurationObject => Duration(seconds: cacheDuration);

  // Print configuration for debugging
  static void printConfig() {
    if (!enableLogging) return;
    
    print('🔧 Configuration Service');
    print('API Base URL: $apiBaseUrl');
    print('API Version: $apiVersion');
    print('Full API URL: $fullApiUrl');
    print('Environment: $environment');
    print('Debug Mode: $debugMode');
    print('App Name: $appName');
    print('App Version: $appVersion');
    print('Default Location: ($defaultLatitude, $defaultLongitude)');
  }

  // Validate configuration
  static bool isConfigValid() {
    try {
      // Check if essential configurations are present
      if (apiBaseUrl.isEmpty || apiVersion.isEmpty) {
        return false;
      }
      
      // Validate URL format
      final uri = Uri.tryParse(fullApiUrl);
      if (uri == null || !uri.hasAbsolutePath) {
        return false;
      }
      
      return true;
    } catch (e) {
      return false;
    }
  }

  // Get environment-specific configurations
  static Map<String, dynamic> getEnvironmentConfig() {
    return {
      'api_base_url': apiBaseUrl,
      'api_version': apiVersion,
      'full_api_url': fullApiUrl,
      'environment': environment,
      'debug_mode': debugMode,
      'app_name': appName,
      'app_version': appVersion,
      'timeouts': {
        'api': apiTimeout,
        'connection': connectionTimeout,
      },
      'cache_duration': cacheDuration,
      'default_location': {
        'latitude': defaultLatitude,
        'longitude': defaultLongitude,
        'zoom': defaultZoom,
      },
    };
  }
}
