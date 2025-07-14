import 'package:flutter/foundation.dart';
import 'package:smg_app/models/topic.dart';
import 'package:smg_app/services/api_service.dart';

class TopicProvider with ChangeNotifier {
  final ApiService _apiService;
  
  List<Topic> _topics = [];
  bool _isLoading = false;
  int _currentPage = 1;
  int _totalPages = 1;
  bool _hasMore = true;

  TopicProvider(this._apiService);

  List<Topic> get topics => _topics;
  bool get isLoading => _isLoading;
  bool get hasMore => _hasMore;

  Future<void> loadTopics({bool refresh = false}) async {
    if (_isLoading) return;
    
    if (refresh) {
      _currentPage = 1;
      _topics.clear();
      _hasMore = true;
    }

    _isLoading = true;
    notifyListeners();

    try {
      final response = await _apiService.getTopics(page: _currentPage);
      
      final List<dynamic> topicData = response['data'];
      final List<Topic> newTopics = topicData.map((json) => Topic.fromJson(json)).toList();
      
      if (refresh) {
        _topics = newTopics;
      } else {
        _topics.addAll(newTopics);
      }
      
      _totalPages = response['total_pages'];
      _hasMore = _currentPage < _totalPages;
      
      if (_hasMore) {
        _currentPage++;
      }
      
    } catch (e) {
      print('Error loading topics: $e');
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<bool> createTopic(String name, String? description, List<String> keywords, List<String> platforms) async {
    try {
      final topic = await _apiService.createTopic(name, description, keywords, platforms);
      _topics.insert(0, topic);
      notifyListeners();
      return true;
    } catch (e) {
      print('Error creating topic: $e');
      return false;
    }
  }

  Future<bool> updateTopic(String topicId, String name, String? description, List<String> keywords, List<String> platforms) async {
    try {
      final updatedTopic = await _apiService.updateTopic(topicId, name, description, keywords, platforms);
      
      final index = _topics.indexWhere((topic) => topic.id == topicId);
      if (index != -1) {
        _topics[index] = updatedTopic;
        notifyListeners();
      }
      
      return true;
    } catch (e) {
      print('Error updating topic: $e');
      return false;
    }
  }

  Future<bool> deleteTopic(String topicId) async {
    try {
      await _apiService.deleteTopic(topicId);
      _topics.removeWhere((topic) => topic.id == topicId);
      notifyListeners();
      return true;
    } catch (e) {
      print('Error deleting topic: $e');
      return false;
    }
  }

  Topic? getTopicById(String topicId) {
    try {
      return _topics.firstWhere((topic) => topic.id == topicId);
    } catch (e) {
      return null;
    }
  }

  void clearTopics() {
    _topics.clear();
    _currentPage = 1;
    _totalPages = 1;
    _hasMore = true;
    notifyListeners();
  }
}