import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/providers.dart';

class TrashScreen extends ConsumerWidget {
  const TrashScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final trashedAsync = ref.watch(trashedNotesProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('Trash')),
      body: trashedAsync.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => Center(child: Text('Error: $e')),
        data: (notes) {
          if (notes.isEmpty) {
            return const Center(child: Text('Trash is empty'));
          }
          return ListView.builder(
            itemCount: notes.length,
            itemBuilder: (context, index) {
              final note = notes[index];
              return Dismissible(
                key: Key(note.id),
                background: Container(
                  color: Colors.green,
                  alignment: Alignment.centerLeft,
                  padding: const EdgeInsets.only(left: 16),
                  child: const Icon(Icons.restore, color: Colors.white),
                ),
                secondaryBackground: Container(
                  color: Colors.red,
                  alignment: Alignment.centerRight,
                  padding: const EdgeInsets.only(right: 16),
                  child: const Icon(Icons.delete_forever, color: Colors.white),
                ),
                confirmDismiss: (direction) async {
                  if (direction == DismissDirection.startToEnd) {
                    await ref.read(notesApiProvider).restore(note.id);
                    ref.invalidate(trashedNotesProvider);
                    return true;
                  } else {
                    await ref.read(notesApiProvider).permanentDelete(note.id);
                    ref.invalidate(trashedNotesProvider);
                    return true;
                  }
                },
                child: ListTile(
                  title: Text(note.title.isEmpty ? 'Untitled' : note.title),
                  subtitle: Text(
                    note.body.length > 80 ? '${note.body.substring(0, 80)}...' : note.body,
                    maxLines: 1,
                  ),
                ),
              );
            },
          );
        },
      ),
    );
  }
}
