import 'package:dio/dio.dart';
import '../models/note.dart';

class SearchApi {
  final Dio dio;

  SearchApi(this.dio);

  Future<List<Note>> search(String query) async {
    final response = await dio.get('/api/search', queryParameters: {'q': query});
    return (response.data as List)
        .map((e) => Note.fromJson(e as Map<String, dynamic>))
        .toList();
  }
}
