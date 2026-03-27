import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'config/theme.dart';
import 'router/app_router.dart';
import 'providers/providers.dart';

void main() {
  runApp(const ProviderScope(child: NoteThingApp()));
}

class NoteThingApp extends ConsumerStatefulWidget {
  const NoteThingApp({super.key});

  @override
  ConsumerState<NoteThingApp> createState() => _NoteThingAppState();
}

class _NoteThingAppState extends ConsumerState<NoteThingApp> {
  @override
  void initState() {
    super.initState();
    ref.read(authStateProvider.notifier).checkStoredToken();
  }

  @override
  Widget build(BuildContext context) {
    final router = ref.watch(routerProvider);

    return MaterialApp.router(
      title: 'Note Thing',
      theme: appTheme(),
      darkTheme: appDarkTheme(),
      routerConfig: router,
      debugShowCheckedModeBanner: false,
    );
  }
}
