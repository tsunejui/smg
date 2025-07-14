import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:go_router/go_router.dart';
import 'package:smg_app/providers/auth_provider.dart';
import 'package:smg_app/providers/topic_provider.dart';
import 'package:smg_app/providers/article_provider.dart';
import 'package:smg_app/services/api_service.dart';
import 'package:smg_app/services/storage_service.dart';
import 'package:smg_app/router/app_router.dart';
import 'package:smg_app/theme/app_theme.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  final storageService = StorageService();
  await storageService.init();
  
  final apiService = ApiService();
  
  runApp(
    MultiProvider(
      providers: [
        ChangeNotifierProvider(
          create: (_) => AuthProvider(apiService, storageService),
        ),
        ChangeNotifierProvider(
          create: (_) => TopicProvider(apiService),
        ),
        ChangeNotifierProvider(
          create: (_) => ArticleProvider(apiService),
        ),
      ],
      child: const SMGApp(),
    ),
  );
}

class SMGApp extends StatelessWidget {
  const SMGApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp.router(
      title: 'SMG App',
      theme: AppTheme.lightTheme,
      darkTheme: AppTheme.darkTheme,
      routerConfig: AppRouter.router,
      debugShowCheckedModeBanner: false,
    );
  }
}