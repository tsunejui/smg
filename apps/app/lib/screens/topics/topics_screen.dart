import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:go_router/go_router.dart';
import 'package:smg_app/providers/topic_provider.dart';
import 'package:smg_app/models/topic.dart';
import 'package:smg_app/widgets/topic_card.dart';
import 'package:smg_app/widgets/create_topic_dialog.dart';

class TopicsScreen extends StatefulWidget {
  const TopicsScreen({super.key});

  @override
  State<TopicsScreen> createState() => _TopicsScreenState();
}

class _TopicsScreenState extends State<TopicsScreen> {
  final ScrollController _scrollController = ScrollController();

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _loadTopics();
    });
    _scrollController.addListener(_onScroll);
  }

  @override
  void dispose() {
    _scrollController.removeListener(_onScroll);
    _scrollController.dispose();
    super.dispose();
  }

  void _loadTopics() {
    final topicProvider = Provider.of<TopicProvider>(context, listen: false);
    topicProvider.loadTopics(refresh: true);
  }

  void _onScroll() {
    if (_scrollController.position.pixels >= _scrollController.position.maxScrollExtent * 0.8) {
      final topicProvider = Provider.of<TopicProvider>(context, listen: false);
      if (topicProvider.hasMore && !topicProvider.isLoading) {
        topicProvider.loadTopics();
      }
    }
  }

  void _showCreateTopicDialog() {
    showDialog(
      context: context,
      builder: (context) => const CreateTopicDialog(),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('主題管理'),
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: _showCreateTopicDialog,
          ),
        ],
      ),
      body: Consumer<TopicProvider>(
        builder: (context, topicProvider, child) {
          if (topicProvider.topics.isEmpty && topicProvider.isLoading) {
            return const Center(
              child: CircularProgressIndicator(),
            );
          }

          if (topicProvider.topics.isEmpty) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(
                    Icons.topic,
                    size: 64,
                    color: Theme.of(context).colorScheme.outline,
                  ),
                  const SizedBox(height: 16),
                  Text(
                    '尚無主題',
                    style: Theme.of(context).textTheme.headlineSmall,
                  ),
                  const SizedBox(height: 8),
                  Text(
                    '點擊右上角的 + 號來建立第一個主題',
                    style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                      color: Theme.of(context).colorScheme.outline,
                    ),
                  ),
                  const SizedBox(height: 24),
                  ElevatedButton.icon(
                    onPressed: _showCreateTopicDialog,
                    icon: const Icon(Icons.add),
                    label: const Text('建立主題'),
                  ),
                ],
              ),
            );
          }

          return RefreshIndicator(
            onRefresh: () async {
              topicProvider.loadTopics(refresh: true);
            },
            child: ListView.builder(
              controller: _scrollController,
              padding: const EdgeInsets.all(16),
              itemCount: topicProvider.topics.length + (topicProvider.hasMore ? 1 : 0),
              itemBuilder: (context, index) {
                if (index == topicProvider.topics.length) {
                  return const Padding(
                    padding: EdgeInsets.all(16),
                    child: Center(
                      child: CircularProgressIndicator(),
                    ),
                  );
                }

                final topic = topicProvider.topics[index];
                return TopicCard(
                  topic: topic,
                  onTap: () {
                    context.push('/topics/${topic.id}');
                  },
                  onEdit: () {
                    _showEditTopicDialog(topic);
                  },
                  onDelete: () {
                    _showDeleteConfirmation(topic);
                  },
                );
              },
            ),
          );
        },
      ),
    );
  }

  void _showEditTopicDialog(Topic topic) {
    showDialog(
      context: context,
      builder: (context) => CreateTopicDialog(topic: topic),
    );
  }

  void _showDeleteConfirmation(Topic topic) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('刪除主題'),
        content: Text('確定要刪除主題「${topic.name}」嗎？此操作無法復原。'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('取消'),
          ),
          TextButton(
            onPressed: () {
              Navigator.of(context).pop();
              _deleteTopic(topic);
            },
            child: const Text('刪除'),
          ),
        ],
      ),
    );
  }

  void _deleteTopic(Topic topic) async {
    final topicProvider = Provider.of<TopicProvider>(context, listen: false);
    final success = await topicProvider.deleteTopic(topic.id);
    
    if (success && mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('主題「${topic.name}」已刪除'),
          backgroundColor: Colors.green,
        ),
      );
    } else if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('刪除失敗，請稍後再試'),
          backgroundColor: Colors.red,
        ),
      );
    }
  }
}