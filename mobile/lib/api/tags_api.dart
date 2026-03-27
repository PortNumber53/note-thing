import 'package:dio/dio.dart';
import '../models/tag.dart';

class TagsApi {
  final Dio dio;

  TagsApi(this.dio);

  Future<List<Tag>> fetchAll() async {
    final response = await dio.get('/api/tags/');
    return (response.data as List)
        .map((e) => Tag.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<Tag> create(String name) async {
    final response = await dio.post('/api/tags/', data: {'name': name});
    return Tag.fromJson(response.data as Map<String, dynamic>);
  }

  Future<void> delete(String id) async {
    await dio.delete('/api/tags/$id');
  }

  Future<void> setNoteTags(String noteId, List<String> tagIds) async {
    await dio.put('/api/notes/$noteId/tags', data: {'tagIds': tagIds});
  }
}
