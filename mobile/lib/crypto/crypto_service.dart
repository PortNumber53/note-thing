import 'dart:convert';
import 'dart:typed_data';
import 'package:cryptography/cryptography.dart';
import 'package:hashlib/hashlib.dart';

const _verifyPlaintext = 'note-thing-verify';

class CryptoService {
  final _aesGcm = AesGcm.with256bits();

  /// Derive KEK from password + salt using Argon2id
  Future<SecretKey> deriveKEK(String password, Uint8List salt) async {
    final argon2 = Argon2(
      iterations: 3,
      memorySizeKB: 65536,
      parallelism: 4,
      hashLength: 32,
      salt: salt,
    );
    final hash = argon2.convert(password.codeUnits);
    return SecretKey(hash.bytes);
  }

  /// Generate a random 256-bit note key
  Future<SecretKey> generateNoteKey() async {
    return _aesGcm.newSecretKey();
  }

  /// Encrypt plaintext field, returns iv + ciphertext bytes
  Future<Uint8List> encryptField(SecretKey key, String plaintext) async {
    final secretBox = await _aesGcm.encryptString(plaintext, secretKey: key);
    // Combine nonce + ciphertext + mac
    final nonce = secretBox.nonce;
    final ct = secretBox.cipherText;
    final mac = secretBox.mac.bytes;
    final result = Uint8List(nonce.length + ct.length + mac.length);
    result.setAll(0, nonce);
    result.setAll(nonce.length, ct);
    result.setAll(nonce.length + ct.length, mac);
    return result;
  }

  /// Decrypt field from iv + ciphertext bytes
  Future<String> decryptField(SecretKey key, Uint8List data) async {
    final nonce = data.sublist(0, 12);
    final macStart = data.length - 16;
    final ct = data.sublist(12, macStart);
    final mac = Mac(data.sublist(macStart));
    final secretBox = SecretBox(ct, nonce: nonce, mac: mac);
    final plainBytes = await _aesGcm.decrypt(secretBox, secretKey: key);
    return utf8.decode(plainBytes);
  }

  /// Wrap DEK with KEK
  Future<Uint8List> wrapKey(SecretKey dek, SecretKey kek) async {
    final dekBytes = await dek.extractBytes();
    final secretBox = await _aesGcm.encrypt(dekBytes, secretKey: kek);
    final nonce = secretBox.nonce;
    final ct = secretBox.cipherText;
    final mac = secretBox.mac.bytes;
    final result = Uint8List(nonce.length + ct.length + mac.length);
    result.setAll(0, nonce);
    result.setAll(nonce.length, ct);
    result.setAll(nonce.length + ct.length, mac);
    return result;
  }

  /// Unwrap DEK from wrapped bytes
  Future<SecretKey> unwrapKey(Uint8List data, SecretKey kek) async {
    final nonce = data.sublist(0, 12);
    final macStart = data.length - 16;
    final ct = data.sublist(12, macStart);
    final mac = Mac(data.sublist(macStart));
    final secretBox = SecretBox(ct, nonce: nonce, mac: mac);
    final dekBytes = await _aesGcm.decrypt(secretBox, secretKey: kek);
    return SecretKey(dekBytes);
  }

  /// Create verify token for password verification
  Future<String> createVerifyToken(SecretKey kek) async {
    final encrypted = await encryptField(kek, _verifyPlaintext);
    return base64Encode(encrypted);
  }

  /// Verify password by decrypting verify token
  Future<bool> verifyPassword(SecretKey kek, String token) async {
    try {
      final data = base64Decode(token);
      final decrypted = await decryptField(kek, data);
      return decrypted == _verifyPlaintext;
    } catch (_) {
      return false;
    }
  }
}
