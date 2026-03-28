class ApiConfig {
  static const String baseUrl = String.fromEnvironment('BACKEND_URL');
  static const bool enableGoogleOAuth = bool.fromEnvironment('ENABLE_GOOGLE_OAUTH', defaultValue: false);

  static const Duration connectTimeout = Duration(seconds: 10);
  static const Duration receiveTimeout = Duration(seconds: 15);
}
