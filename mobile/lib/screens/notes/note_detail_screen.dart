import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_markdown/flutter_markdown.dart';
import 'package:go_router/go_router.dart';
import '../../providers/providers.dart';

class NoteDetailScreen extends ConsumerWidget {
  final String noteId;

  const NoteDetailScreen({super.key, required this.noteId});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final notesAsync = ref.watch(notesProvider((notebookId: null, tagId: null)));

    return notesAsync.when(
      loading: () => const Scaffold(body: Center(child: CircularProgressIndicator())),
      error: (e, _) => Scaffold(body: Center(child: Text('Error: $e'))),
      data: (notes) {
        final note = notes.where((n) => n.id == noteId).firstOrNull;
        if (note == null) {
          return Scaffold(
            appBar: AppBar(),
            body: const Center(child: Text('Note not found')),
          );
        }

        return Scaffold(
          appBar: AppBar(
            title: Text(note.isEncrypted && note.title.isEmpty
                ? 'Encrypted note'
                : note.title.isEmpty
                    ? 'Untitled'
                    : note.title),
            actions: [
              if (!note.isEncrypted || note.title.isNotEmpty)
                IconButton(
                  icon: const Icon(Icons.edit),
                  onPressed: () => context.go('/notes/${note.id}/edit'),
                ),
              IconButton(
                icon: const Icon(Icons.delete_outline),
                onPressed: () async {
                  await ref.read(notesApiProvider).delete(note.id);
                  ref.invalidate(notesProvider((notebookId: null, tagId: null)));
                  if (context.mounted) context.go('/notes');
                },
              ),
            ],
          ),
          body: SelectionArea(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  if (note.tags.isNotEmpty)
                    Wrap(
                      spacing: 4,
                      children: note.tags.map((tag) => Chip(
                        label: Text(tag.name),
                        visualDensity: VisualDensity.compact,
                      )).toList(),
                    ),
                  if (note.tags.isNotEmpty) const SizedBox(height: 12),
                Expanded(
                  child: note.isEncrypted && note.body.isEmpty
                      ? Center(
                          child: Column(
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              Icon(Icons.lock_outline, size: 48, color: Theme.of(context).colorScheme.outline),
                              const SizedBox(height: 12),
                              Text(
                                'This note is encrypted',
                                style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                                      color: Theme.of(context).colorScheme.onSurfaceVariant,
                                    ),
                              ),
                              const SizedBox(height: 4),
                              Text(
                                'Unlock from the web app to view',
                                style: Theme.of(context).textTheme.bodySmall?.copyWith(
                                      color: Theme.of(context).colorScheme.outline,
                                    ),
                              ),
                            ],
                          ),
                        )
                      : Markdown(data: note.body),
                ),
              ],
              ),
            ),
          ),
        );
      },
    );
  }
}
