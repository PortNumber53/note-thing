import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../models/note.dart';
import '../../providers/providers.dart';
import '../../widgets/note_card.dart';

class SearchScreen extends ConsumerStatefulWidget {
  const SearchScreen({super.key});

  @override
  ConsumerState<SearchScreen> createState() => _SearchScreenState();
}

class _SearchScreenState extends ConsumerState<SearchScreen> {
  final _controller = TextEditingController();
  String _query = '';

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    // Load all decrypted notes for client-side search
    final notesAsync = ref.watch(notesProvider((notebookId: null, tagId: null)));

    return Scaffold(
      appBar: AppBar(
        title: TextField(
          controller: _controller,
          decoration: const InputDecoration(
            hintText: 'Search notes...',
            border: InputBorder.none,
          ),
          onChanged: (value) {
            setState(() => _query = value.trim().toLowerCase());
          },
          autofocus: true,
        ),
      ),
      body: notesAsync.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => Center(child: Text('Error: $e')),
        data: (allNotes) {
          final notes = _query.isEmpty
              ? <Note>[]
              : allNotes.where((n) {
                  return n.title.toLowerCase().contains(_query) ||
                      n.body.toLowerCase().contains(_query);
                }).toList();

          if (_query.isEmpty) {
            return Center(
              child: Text(
                'Type to search',
                style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                      color: Theme.of(context).colorScheme.onSurfaceVariant,
                    ),
              ),
            );
          }

          if (notes.isEmpty) {
            return Center(
              child: Text(
                'No results',
                style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                      color: Theme.of(context).colorScheme.onSurfaceVariant,
                    ),
              ),
            );
          }

          return ListView.builder(
            itemCount: notes.length,
            itemBuilder: (context, index) {
              final note = notes[index];
              return NoteCard(
                note: note,
                onTap: () => context.go('/notes/${note.id}'),
              );
            },
          );
        },
      ),
    );
  }
}
