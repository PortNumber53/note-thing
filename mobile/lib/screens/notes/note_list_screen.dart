import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../providers/providers.dart';
import '../../widgets/note_card.dart';

class NoteListScreen extends ConsumerWidget {
  final String? notebookId;
  final String? tagId;

  const NoteListScreen({super.key, this.notebookId, this.tagId});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final notesAsync = ref.watch(notesProvider((
      notebookId: notebookId,
      tagId: tagId,
    )));

    return Scaffold(
      appBar: AppBar(title: const Text('Notes')),
      body: notesAsync.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => Center(child: Text('Error: $e')),
        data: (notes) {
          if (notes.isEmpty) {
            return const Center(child: Text('No notes yet'));
          }
          return RefreshIndicator(
            onRefresh: () => ref.refresh(notesProvider((
              notebookId: notebookId,
              tagId: tagId,
            )).future),
            child: ListView.builder(
              itemCount: notes.length,
              itemBuilder: (context, index) {
                final note = notes[index];
                return NoteCard(
                  note: note,
                  onTap: () => context.go('/notes/${note.id}'),
                );
              },
            ),
          );
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.go('/notes/new'),
        child: const Icon(Icons.add),
      ),
    );
  }
}
