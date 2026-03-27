class Notebook {
  final String id;
  final String name;
  final bool isDefault;
  final int noteCount;
  final DateTime createdAt;
  final DateTime updatedAt;

  Notebook({
    required this.id,
    required this.name,
    required this.isDefault,
    required this.noteCount,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Notebook.fromJson(Map<String, dynamic> json) {
    return Notebook(
      id: json['id'] as String,
      name: json['name'] as String,
      isDefault: json['isDefault'] as bool,
      noteCount: json['noteCount'] as int,
      createdAt: DateTime.parse(json['createdAt'] as String),
      updatedAt: DateTime.parse(json['updatedAt'] as String),
    );
  }
}
