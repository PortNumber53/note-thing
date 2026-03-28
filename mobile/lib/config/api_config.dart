class ApiConfig {
  static const String baseUrl = String.fromEnvironment('BACKEND_URL');

  static const Duration connectTimeout = Duration(seconds: 10);
  static const Duration receiveTimeout = Duration(seconds: 15);
}
