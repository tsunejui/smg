import 'dart:convert';
import 'package:dio/dio.dart';
import 'package:smg_app/models/user.dart';
import 'package:smg_app/models/topic.dart';
import 'package:smg_app/models/article.dart';

class ApiService {
  late Dio _dio;
  static const String baseUrl = 'http://localhost:8080/api';

  ApiService() {
    _dio = Dio(BaseOptions(
      baseUrl: baseUrl,
      headers: {
        'Content-Type': 'application/json',
      },
    ));

    _dio.interceptors.add(LogInterceptor(
      requestBody: true,
      responseBody: true,
    ));
  }

  void setToken(String token) {
    _dio.options.headers['Authorization'] = 'Bearer $token';
  }

  void removeToken() {
    _dio.options.headers.remove('Authorization');
  }

  // Auth endpoints
  Future<Map<String, dynamic>> login(String email, String password) async {
    final response = await _dio.post('/auth/login', data: {
      'email': email,
      'password': password,
    });
    return response.data;
  }

  Future<Map<String, dynamic>> register(String name, String email, String password) async {
    final response = await _dio.post('/auth/register', data: {
      'name': name,
      'email': email,
      'password': password,
    });
    return response.data;
  }

  Future<Map<String, dynamic>> generateQRCode(String userId) async {
    final response = await _dio.post('/auth/qr-generate', data: {
      'user_id': userId,
    });
    return response.data;
  }

  Future<Map<String, dynamic>> verifyQRCode(String token) async {
    final response = await _dio.post('/auth/qr-verify', data: {
      'token': token,
    });
    return response.data;
  }

  // User endpoints
  Future<User> getProfile() async {
    final response = await _dio.get('/users/profile');
    return User.fromJson(response.data);
  }

  Future<User> updateProfile(String? name, String? image) async {
    final response = await _dio.put('/users/profile', data: {
      'name': name,
      'image': image,
    });
    return User.fromJson(response.data);
  }

  // Topic endpoints
  Future<Map<String, dynamic>> getTopics({int page = 1, int pageSize = 20}) async {
    final response = await _dio.get('/topics', queryParameters: {
      'page': page,
      'page_size': pageSize,
    });
    return response.data;
  }

  Future<Topic> createTopic(String name, String? description, List<String> keywords, List<String> platforms) async {
    final response = await _dio.post('/topics', data: {
      'name': name,
      'description': description,
      'keywords': keywords,
      'platforms': platforms,
    });
    return Topic.fromJson(response.data);
  }

  Future<Topic> getTopic(String topicId) async {
    final response = await _dio.get('/topics/$topicId');
    return Topic.fromJson(response.data);
  }

  Future<Topic> updateTopic(String topicId, String name, String? description, List<String> keywords, List<String> platforms) async {
    final response = await _dio.put('/topics/$topicId', data: {
      'name': name,
      'description': description,
      'keywords': keywords,
      'platforms': platforms,
    });
    return Topic.fromJson(response.data);
  }

  Future<void> deleteTopic(String topicId) async {
    await _dio.delete('/topics/$topicId');
  }

  // Article endpoints
  Future<Map<String, dynamic>> getArticles({int page = 1, int pageSize = 20}) async {
    final response = await _dio.get('/articles', queryParameters: {
      'page': page,
      'page_size': pageSize,
    });
    return response.data;
  }

  Future<Article> getArticle(String articleId) async {
    final response = await _dio.get('/articles/$articleId');
    return Article.fromJson(response.data);
  }

  Future<Map<String, dynamic>> repostArticle(String articleId, String mediaAccountId, String? customCaption) async {
    final response = await _dio.post('/articles/$articleId/repost', data: {
      'media_account_id': mediaAccountId,
      'custom_caption': customCaption,
    });
    return response.data;
  }

  Future<Map<String, dynamic>> getReposts({int page = 1, int pageSize = 20}) async {
    final response = await _dio.get('/articles/reposts', queryParameters: {
      'page': page,
      'page_size': pageSize,
    });
    return response.data;
  }

  // Media endpoints
  Future<List<dynamic>> getMediaAccounts() async {
    final response = await _dio.get('/media/accounts');
    return response.data;
  }

  Future<Map<String, dynamic>> connectPlatform(String platform, String code, String accountName) async {
    final response = await _dio.post('/media/connect/$platform', data: {
      'code': code,
      'account_name': accountName,
    });
    return response.data;
  }

  Future<void> disconnectAccount(String accountId) async {
    await _dio.post('/media/disconnect/$accountId');
  }
}