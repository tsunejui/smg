import 'package:flutter/foundation.dart';
import 'package:smg_app/models/article.dart';
import 'package:smg_app/services/api_service.dart';

class ArticleProvider with ChangeNotifier {
  final ApiService _apiService;
  
  List<Article> _articles = [];
  Map<String, List<Article>> _articlesByTopic = {};
  bool _isLoading = false;
  int _currentPage = 1;
  int _totalPages = 1;
  bool _hasMore = true;

  ArticleProvider(this._apiService);

  List<Article> get articles => _articles;
  Map<String, List<Article>> get articlesByTopic => _articlesByTopic;
  bool get isLoading => _isLoading;
  bool get hasMore => _hasMore;

  Future<void> loadArticles({bool refresh = false}) async {
    if (_isLoading) return;
    
    if (refresh) {
      _currentPage = 1;
      _articles.clear();
      _hasMore = true;
    }

    _isLoading = true;
    notifyListeners();

    try {
      final response = await _apiService.getArticles(page: _currentPage);
      
      final List<dynamic> articleData = response['data'];
      final List<Article> newArticles = articleData.map((json) => Article.fromJson(json)).toList();
      
      if (refresh) {
        _articles = newArticles;
      } else {
        _articles.addAll(newArticles);
      }
      
      _totalPages = response['total_pages'];
      _hasMore = _currentPage < _totalPages;
      
      if (_hasMore) {
        _currentPage++;
      }
      
      _groupArticlesByTopic();
      
    } catch (e) {
      print('Error loading articles: $e');
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  void _groupArticlesByTopic() {
    _articlesByTopic.clear();
    
    for (final article in _articles) {
      if (!_articlesByTopic.containsKey(article.topicId)) {
        _articlesByTopic[article.topicId] = [];
      }
      _articlesByTopic[article.topicId]!.add(article);
    }
  }

  List<Article> getArticlesByTopicId(String topicId) {
    return _articlesByTopic[topicId] ?? [];
  }

  Future<Article?> getArticle(String articleId) async {
    try {
      return await _apiService.getArticle(articleId);
    } catch (e) {
      print('Error getting article: $e');
      return null;
    }
  }

  Future<bool> repostArticle(String articleId, String mediaAccountId, String? customCaption) async {
    try {
      await _apiService.repostArticle(articleId, mediaAccountId, customCaption);
      return true;
    } catch (e) {
      print('Error reposting article: $e');
      return false;
    }
  }

  Future<void> loadReposts({bool refresh = false}) async {
    if (_isLoading) return;
    
    if (refresh) {
      _currentPage = 1;
      _hasMore = true;
    }

    _isLoading = true;
    notifyListeners();

    try {
      final response = await _apiService.getReposts(page: _currentPage);
      
      // Handle reposts data here if needed
      
      _totalPages = response['total_pages'];
      _hasMore = _currentPage < _totalPages;
      
      if (_hasMore) {
        _currentPage++;
      }
      
    } catch (e) {
      print('Error loading reposts: $e');
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  void clearArticles() {
    _articles.clear();
    _articlesByTopic.clear();
    _currentPage = 1;
    _totalPages = 1;
    _hasMore = true;
    notifyListeners();
  }
}