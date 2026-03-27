import 'tag.dart';

class Note {
  final String id;
  final String title;
  final String body;
  final String? encryptedTitle;
  final String? encryptedBody;
  final String? noteKeyWrapped;
  final int? keyVersion;
  final bool isEncrypted;
  final String? notebookId;
  final List<Tag> tags;
  final DateTime createdAt;
  final DateTime updatedAt;

  Note({
    required this.id,
    required this.title,
    required this.body,
    this.encryptedTitle,
    this.encryptedBody,
    this.noteKeyWrapped,
    this.keyVersion,
    this.isEncrypted = false,
    this.notebookId,
    required this.tags,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Note.fromJson(Map<String, dynamic> json) {
    return Note(
      id: json['id'] as String,
      title: json['title'] as String? ?? '',
      body: json['body'] as String? ?? '',
      encryptedTitle: json['encryptedTitle'] as String?,
      encryptedBody: json['encryptedBody'] as String?,
      noteKeyWrapped: json['noteKeyWrapped'] as String?,
      keyVersion: json['keyVersion'] as int?,
      isEncrypted: json['isEncrypted'] as bool? ?? false,
      notebookId: json['notebookId'] as String?,
      tags: (json['tags'] as List<dynamic>?)
              ?.map((t) => Tag.fromJson(t as Map<String, dynamic>))
              .toList() ??
          [],
      createdAt: DateTime.parse(json['createdAt'] as String),
      updatedAt: DateTime.parse(json['updatedAt'] as String),
    );
  }

  Note copyWith({String? title, String? body}) {
    return Note(
      id: id,
      title: title ?? this.title,
      body: body ?? this.body,
      encryptedTitle: encryptedTitle,
      encryptedBody: encryptedBody,
      noteKeyWrapped: noteKeyWrapped,
      keyVersion: keyVersion,
      isEncrypted: isEncrypted,
      notebookId: notebookId,
      tags: tags,
      createdAt: createdAt,
      updatedAt: updatedAt,
    );
  }
}
