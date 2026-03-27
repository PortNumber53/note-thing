import 'dart:convert';
import 'dart:typed_data';
import 'package:cryptography/cryptography.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../api/api_client.dart';
import '../providers/providers.dart';
import 'crypto_service.dart';

final cryptoServiceProvider = Provider((_) => CryptoService());

class CryptoState {
  final bool isEncryptionEnabled;
  final bool isUnlocked;
  final SecretKey? kek;
  final int? keyVersion;
  final Uint8List? kdfSalt;
  final String? kekVerifyToken;

  CryptoState({
    this.isEncryptionEnabled = false,
    this.isUnlocked = false,
    this.kek,
    this.keyVersion,
    this.kdfSalt,
    this.kekVerifyToken,
  });

  CryptoState copyWith({
    bool? isEncryptionEnabled,
    bool? isUnlocked,
    SecretKey? kek,
    int? keyVersion,
    Uint8List? kdfSalt,
    String? kekVerifyToken,
  }) {
    return CryptoState(
      isEncryptionEnabled: isEncryptionEnabled ?? this.isEncryptionEnabled,
      isUnlocked: isUnlocked ?? this.isUnlocked,
      kek: kek ?? this.kek,
      keyVersion: keyVersion ?? this.keyVersion,
      kdfSalt: kdfSalt ?? this.kdfSalt,
      kekVerifyToken: kekVerifyToken ?? this.kekVerifyToken,
    );
  }
}

final cryptoStateProvider = StateNotifierProvider<CryptoNotifier, CryptoState>(
  (ref) => CryptoNotifier(ref),
);

class CryptoNotifier extends StateNotifier<CryptoState> {
  final Ref ref;

  CryptoNotifier(this.ref) : super(CryptoState());

  CryptoService get _crypto => ref.read(cryptoServiceProvider);
  ApiClient get _api => ref.read(apiClientProvider);

  Future<void> fetchEncryptionStatus() async {
    try {
      final response = await _api.dio.get('/api/encryption');
      final data = response.data as Map<String, dynamic>;
      if (data['enabled'] == true) {
        state = state.copyWith(
          isEncryptionEnabled: true,
          kdfSalt: base64Decode(data['kdfSalt'] as String),
          keyVersion: data['keyVersion'] as int,
          kekVerifyToken: data['kekVerify'] as String,
        );
      }
    } catch (_) {}
  }

  Future<bool> unlock(String password) async {
    if (state.kdfSalt == null || state.kekVerifyToken == null) return false;
    try {
      final kek = await _crypto.deriveKEK(password, state.kdfSalt!);
      final valid = await _crypto.verifyPassword(kek, state.kekVerifyToken!);
      if (!valid) return false;
      state = state.copyWith(kek: kek, isUnlocked: true);
      return true;
    } catch (_) {
      return false;
    }
  }

  void lock() {
    state = CryptoState(
      isEncryptionEnabled: state.isEncryptionEnabled,
      keyVersion: state.keyVersion,
      kdfSalt: state.kdfSalt,
      kekVerifyToken: state.kekVerifyToken,
    );
  }
}
