class Topic {
  final String id;
  final String name;
  final String? description;
  final List<String> keywords;
  final List<String> platforms;
  final String userId;
  final DateTime createdAt;
  final DateTime updatedAt;

  Topic({
    required this.id,
    required this.name,
    this.description,
    required this.keywords,
    required this.platforms,
    required this.userId,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Topic.fromJson(Map<String, dynamic> json) {
    return Topic(
      id: json['id'],
      name: json['name'],
      description: json['description'],
      keywords: List<String>.from(json['keywords']),
      platforms: List<String>.from(json['platforms']),
      userId: json['user_id'],
      createdAt: DateTime.parse(json['created_at']),
      updatedAt: DateTime.parse(json['updated_at']),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'description': description,
      'keywords': keywords,
      'platforms': platforms,
      'user_id': userId,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }
}