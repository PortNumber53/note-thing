import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../providers/providers.dart';

class NotebookListScreen extends ConsumerWidget {
  const NotebookListScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final notebooksAsync = ref.watch(notebooksProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('Notebooks')),
      body: notebooksAsync.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => Center(child: Text('Error: $e')),
        data: (notebooks) {
          if (notebooks.isEmpty) {
            return const Center(child: Text('No notebooks'));
          }
          return ListView.builder(
            itemCount: notebooks.length,
            itemBuilder: (context, index) {
              final nb = notebooks[index];
              return ListTile(
                leading: Icon(
                  nb.isDefault ? Icons.book : Icons.book_outlined,
                ),
                title: Text(nb.name),
                trailing: Text(
                  '${nb.noteCount}',
                  style: Theme.of(context).textTheme.bodySmall,
                ),
                onTap: () => context.go('/notes'),
              );
            },
          );
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _showCreateDialog(context, ref),
        child: const Icon(Icons.add),
      ),
    );
  }

  void _showCreateDialog(BuildContext context, WidgetRef ref) {
    final controller = TextEditingController();
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('New Notebook'),
        content: TextField(
          controller: controller,
          decoration: const InputDecoration(hintText: 'Notebook name'),
          autofocus: true,
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('Cancel')),
          FilledButton(
            onPressed: () async {
              if (controller.text.trim().isNotEmpty) {
                await ref.read(notebooksApiProvider).create(controller.text.trim());
                ref.invalidate(notebooksProvider);
                if (ctx.mounted) Navigator.pop(ctx);
              }
            },
            child: const Text('Create'),
          ),
        ],
      ),
    );
  }
}
