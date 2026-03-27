import 'package:dio/dio.dart';
import '../models/notebook.dart';

class NotebooksApi {
  final Dio dio;

  NotebooksApi(this.dio);

  Future<List<Notebook>> fetchAll() async {
    final response = await dio.get('/api/notebooks/');
    return (response.data as List)
        .map((e) => Notebook.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<Notebook> create(String name) async {
    final response = await dio.post('/api/notebooks/', data: {'name': name});
    return Notebook.fromJson(response.data as Map<String, dynamic>);
  }

  Future<Notebook> update(String id, String name) async {
    final response = await dio.put('/api/notebooks/$id', data: {'name': name});
    return Notebook.fromJson(response.data as Map<String, dynamic>);
  }

  Future<void> delete(String id) async {
    await dio.delete('/api/notebooks/$id');
  }
}
