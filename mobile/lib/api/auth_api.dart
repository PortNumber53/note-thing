import 'package:dio/dio.dart';
import '../models/user.dart';

class AuthResult {
  final String token;
  final User user;

  AuthResult({required this.token, required this.user});
}

class AuthApi {
  final Dio dio;

  AuthApi(this.dio);

  Future<AuthResult> loginWithGoogle(String idToken) async {
    final response = await dio.post('/auth/google/token', data: {
      'idToken': idToken,
    });
    return AuthResult(
      token: response.data['token'] as String,
      user: User.fromJson(response.data['user'] as Map<String, dynamic>),
    );
  }

  Future<AuthResult> signup(String email, String password, String name) async {
    final response = await dio.post('/auth/signup', data: {
      'email': email,
      'password': password,
      'name': name,
    });
    return AuthResult(
      token: response.data['token'] as String,
      user: User.fromJson(response.data['user'] as Map<String, dynamic>),
    );
  }

  Future<AuthResult> login(String email, String password) async {
    final response = await dio.post('/auth/login', data: {
      'email': email,
      'password': password,
    });
    return AuthResult(
      token: response.data['token'] as String,
      user: User.fromJson(response.data['user'] as Map<String, dynamic>),
    );
  }

  Future<User> getMe() async {
    final response = await dio.get('/api/me');
    return User.fromJson(response.data as Map<String, dynamic>);
  }
}
