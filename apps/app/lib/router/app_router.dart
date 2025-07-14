import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:smg_app/screens/splash_screen.dart';
import 'package:smg_app/screens/auth/login_screen.dart';
import 'package:smg_app/screens/auth/qr_login_screen.dart';
import 'package:smg_app/screens/home/home_screen.dart';
import 'package:smg_app/screens/topics/topics_screen.dart';
import 'package:smg_app/screens/topics/topic_detail_screen.dart';
import 'package:smg_app/screens/articles/articles_screen.dart';
import 'package:smg_app/screens/articles/article_detail_screen.dart';
import 'package:smg_app/screens/profile/profile_screen.dart';

class AppRouter {
  static final GoRouter router = GoRouter(
    initialLocation: '/',
    routes: [
      GoRoute(
        path: '/',
        builder: (context, state) => const SplashScreen(),
      ),
      GoRoute(
        path: '/login',
        builder: (context, state) => const LoginScreen(),
      ),
      GoRoute(
        path: '/qr-login',
        builder: (context, state) => const QRLoginScreen(),
      ),
      ShellRoute(
        builder: (context, state, child) => HomeScreen(child: child),
        routes: [
          GoRoute(
            path: '/home',
            builder: (context, state) => const HomeDashboard(),
          ),
          GoRoute(
            path: '/topics',
            builder: (context, state) => const TopicsScreen(),
          ),
          GoRoute(
            path: '/topics/:id',
            builder: (context, state) => TopicDetailScreen(
              topicId: state.pathParameters['id']!,
            ),
          ),
          GoRoute(
            path: '/articles',
            builder: (context, state) => const ArticlesScreen(),
          ),
          GoRoute(
            path: '/articles/:id',
            builder: (context, state) => ArticleDetailScreen(
              articleId: state.pathParameters['id']!,
            ),
          ),
          GoRoute(
            path: '/profile',
            builder: (context, state) => const ProfileScreen(),
          ),
        ],
      ),
    ],
  );
}

class HomeDashboard extends StatelessWidget {
  const HomeDashboard({super.key});

  @override
  Widget build(BuildContext context) {
    return const Center(
      child: Text('Home Dashboard'),
    );
  }
}