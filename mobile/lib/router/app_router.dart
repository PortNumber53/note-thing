import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../crypto/crypto_provider.dart';
import '../providers/providers.dart';
import '../screens/auth/login_screen.dart';
import '../screens/auth/unlock_screen.dart';
import '../screens/home/home_screen.dart';
import '../screens/notes/note_list_screen.dart';
import '../screens/notes/note_detail_screen.dart';
import '../screens/notes/note_editor_screen.dart';
import '../screens/notebooks/notebook_list_screen.dart';
import '../screens/search/search_screen.dart';
import '../screens/trash/trash_screen.dart';
import '../screens/settings/settings_screen.dart';

final routerProvider = Provider<GoRouter>((ref) {
  final authState = ref.watch(authStateProvider);
  final cryptoState = ref.watch(cryptoStateProvider);

  return GoRouter(
    initialLocation: '/notes',
    redirect: (context, state) {
      final isLoggedIn = authState.valueOrNull != null;
      final isLoginRoute = state.matchedLocation == '/login';
      final isUnlockRoute = state.matchedLocation == '/unlock';

      if (!isLoggedIn && !isLoginRoute) return '/login';
      if (isLoggedIn && isLoginRoute) {
        // After login, check if we need to unlock
        if (cryptoState.isEncryptionEnabled && !cryptoState.isUnlocked) {
          return '/unlock';
        }
        return '/notes';
      }

      // If encryption is enabled but not unlocked, redirect to unlock
      if (isLoggedIn && cryptoState.isEncryptionEnabled && !cryptoState.isUnlocked && !isUnlockRoute) {
        return '/unlock';
      }

      // If unlocked and on unlock page, go to notes
      if (isLoggedIn && cryptoState.isUnlocked && isUnlockRoute) {
        return '/notes';
      }

      return null;
    },
    routes: [
      GoRoute(path: '/login', builder: (_, _) => const LoginScreen()),
      GoRoute(path: '/unlock', builder: (_, _) => const UnlockScreen()),
      ShellRoute(
        builder: (_, _, child) => HomeScreen(child: child),
        routes: [
          GoRoute(
            path: '/notes',
            builder: (_, _) => const NoteListScreen(),
            routes: [
              GoRoute(
                path: 'new',
                builder: (_, _) => const NoteEditorScreen(),
              ),
              GoRoute(
                path: ':id',
                builder: (_, state) => NoteDetailScreen(
                  noteId: state.pathParameters['id']!,
                ),
                routes: [
                  GoRoute(
                    path: 'edit',
                    builder: (_, state) => NoteEditorScreen(
                      noteId: state.pathParameters['id'],
                    ),
                  ),
                ],
              ),
            ],
          ),
          GoRoute(
            path: '/notebooks',
            builder: (_, _) => const NotebookListScreen(),
          ),
          GoRoute(
            path: '/search',
            builder: (_, _) => const SearchScreen(),
          ),
          GoRoute(
            path: '/trash',
            builder: (_, _) => const TrashScreen(),
          ),
          GoRoute(
            path: '/settings',
            builder: (_, _) => const SettingsScreen(),
          ),
        ],
      ),
    ],
  );
});
