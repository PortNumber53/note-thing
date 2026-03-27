import 'package:dio/dio.dart';
import '../models/note.dart';

class NotesApi {
  final Dio dio;

  NotesApi(this.dio);

  Future<List<Note>> fetchAll({String? notebookId, String? tagId}) async {
    final params = <String, String>{};
    if (notebookId != null) params['notebook_id'] = notebookId;
    if (tagId != null) params['tag_id'] = tagId;
    final response = await dio.get('/api/notes/', queryParameters: params);
    return (response.data as List)
        .map((e) => Note.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<Note> fetchById(String id) async {
    final response = await dio.get('/api/notes/$id');
    return Note.fromJson(response.data as Map<String, dynamic>);
  }

  Future<Note> create({
    required String title,
    required String body,
    String? notebookId,
    List<String>? tagIds,
  }) async {
    final response = await dio.post('/api/notes/', data: {
      'title': title,
      'body': body,
      if (notebookId != null) 'notebookId': notebookId,
      if (tagIds != null) 'tagIds': tagIds,
    });
    return Note.fromJson(response.data as Map<String, dynamic>);
  }

  Future<Note> update(String id, {String? title, String? body, String? notebookId}) async {
    final response = await dio.put('/api/notes/$id', data: {
      if (title != null) 'title': title,
      if (body != null) 'body': body,
      if (notebookId != null) 'notebookId': notebookId,
    });
    return Note.fromJson(response.data as Map<String, dynamic>);
  }

  Future<void> delete(String id) async {
    await dio.delete('/api/notes/$id');
  }

  Future<void> restore(String id) async {
    await dio.post('/api/notes/$id/restore');
  }

  Future<void> permanentDelete(String id) async {
    await dio.delete('/api/notes/$id/permanent');
  }

  Future<List<Note>> fetchTrashed() async {
    final response = await dio.get('/api/notes/trash');
    return (response.data as List)
        .map((e) => Note.fromJson(e as Map<String, dynamic>))
        .toList();
  }
}
