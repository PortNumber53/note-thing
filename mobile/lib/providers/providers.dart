import 'dart:convert';
import 'dart:typed_data';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../api/api_client.dart';
import '../api/auth_api.dart';
import '../api/notes_api.dart';
import '../api/notebooks_api.dart';
import '../api/tags_api.dart';
import '../api/search_api.dart';
import '../crypto/crypto_provider.dart';
import '../crypto/crypto_service.dart';
import '../models/user.dart';
import '../models/note.dart';
import '../models/notebook.dart';
import '../models/tag.dart';

// Core
final secureStorageProvider = Provider((_) => const FlutterSecureStorage());

final apiClientProvider = Provider((ref) {
  return ApiClient(storage: ref.read(secureStorageProvider));
});

// APIs
final authApiProvider = Provider((ref) => AuthApi(ref.read(apiClientProvider).dio));
final notesApiProvider = Provider((ref) => NotesApi(ref.read(apiClientProvider).dio));
final notebooksApiProvider = Provider((ref) => NotebooksApi(ref.read(apiClientProvider).dio));
final tagsApiProvider = Provider((ref) => TagsApi(ref.read(apiClientProvider).dio));
final searchApiProvider = Provider((ref) => SearchApi(ref.read(apiClientProvider).dio));

// Auth state
final authStateProvider = StateNotifierProvider<AuthNotifier, AsyncValue<User?>>(
  (ref) => AuthNotifier(ref),
);

class AuthNotifier extends StateNotifier<AsyncValue<User?>> {
  final Ref ref;

  AuthNotifier(this.ref) : super(const AsyncValue.data(null));

  Future<void> checkStoredToken() async {
    final storage = ref.read(secureStorageProvider);
    final token = await storage.read(key: 'jwt_token');
    if (token != null) {
      try {
        final user = await ref.read(authApiProvider).getMe();
        state = AsyncValue.data(user);
      } catch (_) {
        await storage.delete(key: 'jwt_token');
        state = const AsyncValue.data(null);
      }
    }
  }

  Future<void> loginWithGoogle(String idToken) async {
    state = const AsyncValue.loading();
    try {
      final result = await ref.read(authApiProvider).loginWithGoogle(idToken);
      await ref.read(secureStorageProvider).write(key: 'jwt_token', value: result.token);
      state = AsyncValue.data(result.user);
    } catch (e, st) {
      state = AsyncValue.error(e, st);
    }
  }

  Future<void> signupWithEmail(String email, String password, String name) async {
    state = const AsyncValue.loading();
    try {
      final result = await ref.read(authApiProvider).signup(email, password, name);
      await ref.read(secureStorageProvider).write(key: 'jwt_token', value: result.token);
      state = AsyncValue.data(result.user);
      await ref.read(cryptoStateProvider.notifier).fetchEncryptionStatus();
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }

  Future<void> loginWithEmail(String email, String password) async {
    state = const AsyncValue.loading();
    try {
      final result = await ref.read(authApiProvider).login(email, password);
      await ref.read(secureStorageProvider).write(key: 'jwt_token', value: result.token);
      state = AsyncValue.data(result.user);
      await ref.read(cryptoStateProvider.notifier).fetchEncryptionStatus();
    } catch (e, st) {
      state = AsyncValue.error(e, st);
      rethrow;
    }
  }

  Future<void> logout() async {
    await ref.read(secureStorageProvider).delete(key: 'jwt_token');
    state = const AsyncValue.data(null);
  }
}

// Notes - use a record as family key for value equality
final notesProvider = FutureProvider.family<List<Note>, ({String? notebookId, String? tagId})>((ref, filters) async {
  final notes = await ref.read(notesApiProvider).fetchAll(
    notebookId: filters.notebookId,
    tagId: filters.tagId,
  );
  return _decryptNotes(ref, notes);
});

final trashedNotesProvider = FutureProvider<List<Note>>((ref) async {
  final notes = await ref.read(notesApiProvider).fetchTrashed();
  return _decryptNotes(ref, notes);
});

Future<List<Note>> _decryptNotes(Ref ref, List<Note> notes) async {
  final cryptoState = ref.read(cryptoStateProvider);
  if (!cryptoState.isUnlocked || cryptoState.kek == null) return notes;

  final crypto = ref.read(cryptoServiceProvider);
  final result = <Note>[];
  for (final note in notes) {
    if (note.isEncrypted &&
        note.encryptedTitle != null &&
        note.encryptedBody != null &&
        note.noteKeyWrapped != null) {
      try {
        final dek = await crypto.unwrapKey(
          base64Decode(note.noteKeyWrapped!),
          cryptoState.kek!,
        );
        final title = await crypto.decryptField(dek, base64Decode(note.encryptedTitle!));
        final body = await crypto.decryptField(dek, base64Decode(note.encryptedBody!));
        result.add(note.copyWith(title: title, body: body));
      } catch (_) {
        result.add(note.copyWith(title: '[Decryption failed]'));
      }
    } else {
      result.add(note);
    }
  }
  return result;
}

// Notebooks
final notebooksProvider = FutureProvider<List<Notebook>>((ref) {
  return ref.read(notebooksApiProvider).fetchAll();
});

// Tags
final tagsProvider = FutureProvider<List<Tag>>((ref) {
  return ref.read(tagsApiProvider).fetchAll();
});

// Search
final searchQueryProvider = StateProvider<String>((ref) => '');

final searchResultsProvider = FutureProvider<List<Note>>((ref) {
  final query = ref.watch(searchQueryProvider);
  if (query.isEmpty) return Future.value([]);
  return ref.read(searchApiProvider).search(query);
});
