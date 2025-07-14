class Article {
  final String id;
  final String title;
  final String content;
  final String originalUrl;
  final String platform;
  final String? authorName;
  final String? authorId;
  final DateTime publishedAt;
  final String topicId;
  final String userId;
  final DateTime createdAt;
  final DateTime updatedAt;

  Article({
    required this.id,
    required this.title,
    required this.content,
    required this.originalUrl,
    required this.platform,
    this.authorName,
    this.authorId,
    required this.publishedAt,
    required this.topicId,
    required this.userId,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Article.fromJson(Map<String, dynamic> json) {
    return Article(
      id: json['id'],
      title: json['title'],
      content: json['content'],
      originalUrl: json['original_url'],
      platform: json['platform'],
      authorName: json['author_name'],
      authorId: json['author_id'],
      publishedAt: DateTime.parse(json['published_at']),
      topicId: json['topic_id'],
      userId: json['user_id'],
      createdAt: DateTime.parse(json['created_at']),
      updatedAt: DateTime.parse(json['updated_at']),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'content': content,
      'original_url': originalUrl,
      'platform': platform,
      'author_name': authorName,
      'author_id': authorId,
      'published_at': publishedAt.toIso8601String(),
      'topic_id': topicId,
      'user_id': userId,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }
}