import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:intl/intl.dart';
import '../models/note.dart';

class NoteCard extends StatelessWidget {
  final Note note;
  final VoidCallback onTap;

  const NoteCard({super.key, required this.note, required this.onTap});

  @override
  Widget build(BuildContext context) {
    final preview = note.body
        .replaceAllMapped(RegExp(r'!\[([^\]]*)\]\([^)]*\)'), (m) => m[1]!)  // images
        .replaceAllMapped(RegExp(r'\[([^\]]*)\]\([^)]*\)'), (m) => m[1]!)   // links
        .replaceAll(RegExp(r'^#{1,6}\s+', multiLine: true), '') // headings
        .replaceAll(RegExp(r'^>\s?', multiLine: true), '')      // blockquotes
        .replaceAll(RegExp(r'^[-*+]\s+', multiLine: true), '')  // unordered lists
        .replaceAll(RegExp(r'^\d+\\?\.\s+', multiLine: true), '') // ordered lists
        .replaceAll(RegExp(r'^---+$', multiLine: true), '')     // horizontal rules
        .replaceAllMapped(RegExp(r'\\([^\s])'), (m) => m[1]!)   // escaped chars
        .replaceAll(RegExp(r'[*_`~\[\]]'), '')                  // remaining inline formatting
        .trim();
    final dateStr = DateFormat.MMMd().format(note.updatedAt);

    return InkWell(
      onTap: onTap,
      onLongPress: () {
        if (note.body.isNotEmpty) {
          Clipboard.setData(ClipboardData(text: note.body));
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(content: Text('Note copied to clipboard'), duration: Duration(seconds: 2)),
          );
        }
      },
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Expanded(
                  child: Row(
                    children: [
                      if (note.isEncrypted && note.title.isEmpty)
                        Padding(
                          padding: const EdgeInsets.only(right: 4),
                          child: Icon(Icons.lock, size: 14, color: Theme.of(context).colorScheme.primary),
                        ),
                      Expanded(
                        child: Text(
                          note.isEncrypted && note.title.isEmpty
                              ? 'Encrypted note'
                              : note.title.isEmpty
                                  ? 'Untitled'
                                  : note.title,
                          style: Theme.of(context).textTheme.titleSmall?.copyWith(
                                fontStyle: note.isEncrypted && note.title.isEmpty ? FontStyle.italic : null,
                              ),
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ),
                    ],
                  ),
                ),
                Text(
                  dateStr,
                  style: Theme.of(context).textTheme.bodySmall?.copyWith(
                        color: Theme.of(context).colorScheme.onSurfaceVariant,
                      ),
                ),
              ],
            ),
            if (preview.isNotEmpty) ...[
              const SizedBox(height: 4),
              Text(
                preview,
                maxLines: 2,
                overflow: TextOverflow.ellipsis,
                style: Theme.of(context).textTheme.bodySmall?.copyWith(
                      color: Theme.of(context).colorScheme.onSurfaceVariant,
                    ),
              ),
            ],
            if (note.tags.isNotEmpty) ...[
              const SizedBox(height: 6),
              Wrap(
                spacing: 4,
                children: note.tags
                    .map((tag) => Chip(
                          label: Text(tag.name),
                          visualDensity: VisualDensity.compact,
                          padding: EdgeInsets.zero,
                          labelStyle: Theme.of(context).textTheme.labelSmall,
                        ))
                    .toList(),
              ),
            ],
            const Divider(height: 1),
          ],
        ),
      ),
    );
  }
}
