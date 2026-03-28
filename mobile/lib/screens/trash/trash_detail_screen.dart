import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_markdown/flutter_markdown.dart';
import 'package:go_router/go_router.dart';
import '../../providers/providers.dart';

class TrashDetailScreen extends ConsumerWidget {
  final String noteId;

  const TrashDetailScreen({super.key, required this.noteId});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final trashedAsync = ref.watch(trashedNotesProvider);

    return trashedAsync.when(
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
            title: Text(note.title.isEmpty ? 'Untitled' : note.title),
            actions: [
              IconButton(
                icon: const Icon(Icons.restore),
                tooltip: 'Restore',
                onPressed: () async {
                  await ref.read(notesApiProvider).restore(note.id);
                  ref.invalidate(trashedNotesProvider);
                  if (context.mounted) context.go('/trash');
                },
              ),
              IconButton(
                icon: Icon(Icons.delete_forever, color: Theme.of(context).colorScheme.error),
                tooltip: 'Delete permanently',
                onPressed: () async {
                  final confirm = await showDialog<bool>(
                    context: context,
                    builder: (ctx) => AlertDialog(
                      title: const Text('Delete permanently?'),
                      content: const Text('This cannot be undone.'),
                      actions: [
                        TextButton(onPressed: () => Navigator.pop(ctx, false), child: const Text('Cancel')),
                        TextButton(
                          onPressed: () => Navigator.pop(ctx, true),
                          child: Text('Delete', style: TextStyle(color: Theme.of(context).colorScheme.error)),
                        ),
                      ],
                    ),
                  );
                  if (confirm == true) {
                    await ref.read(notesApiProvider).permanentDelete(note.id);
                    ref.invalidate(trashedNotesProvider);
                    if (context.mounted) context.go('/trash');
                  }
                },
              ),
            ],
          ),
          body: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Container(
                width: double.infinity,
                padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                color: Theme.of(context).colorScheme.errorContainer,
                child: Row(
                  children: [
                    Icon(Icons.delete_outline, size: 16, color: Theme.of(context).colorScheme.onErrorContainer),
                    const SizedBox(width: 8),
                    Text(
                      'This note is in the trash',
                      style: Theme.of(context).textTheme.bodySmall?.copyWith(
                            color: Theme.of(context).colorScheme.onErrorContainer,
                          ),
                    ),
                  ],
                ),
              ),
              if (note.tags.isNotEmpty)
                Padding(
                  padding: const EdgeInsets.fromLTRB(16, 12, 16, 0),
                  child: Wrap(
                    spacing: 4,
                    children: note.tags.map((tag) => Chip(
                      label: Text(tag.name),
                      visualDensity: VisualDensity.compact,
                    )).toList(),
                  ),
                ),
              Expanded(
                child: SelectionArea(
                  child: Padding(
                    padding: const EdgeInsets.all(16),
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
                              ],
                            ),
                          )
                        : Markdown(data: note.body),
                  ),
                ),
              ),
            ],
          ),
        );
      },
    );
  }
}
