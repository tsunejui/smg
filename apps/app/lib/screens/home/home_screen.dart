import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

class HomeScreen extends StatefulWidget {
  final Widget child;
  
  const HomeScreen({super.key, required this.child});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  int _currentIndex = 0;

  final List<NavigationDestination> _destinations = [
    const NavigationDestination(
      icon: Icon(Icons.home),
      label: '首頁',
    ),
    const NavigationDestination(
      icon: Icon(Icons.topic),
      label: '主題',
    ),
    const NavigationDestination(
      icon: Icon(Icons.article),
      label: '文章',
    ),
    const NavigationDestination(
      icon: Icon(Icons.person),
      label: '個人',
    ),
  ];

  final List<String> _routes = [
    '/home',
    '/topics',
    '/articles',
    '/profile',
  ];

  void _onDestinationSelected(int index) {
    setState(() {
      _currentIndex = index;
    });
    context.go(_routes[index]);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: widget.child,
      bottomNavigationBar: NavigationBar(
        selectedIndex: _currentIndex,
        onDestinationSelected: _onDestinationSelected,
        destinations: _destinations,
      ),
    );
  }
}