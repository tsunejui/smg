import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:smg_app/models/user.dart';
import 'package:smg_app/services/api_service.dart';
import 'package:smg_app/services/storage_service.dart';

class AuthProvider with ChangeNotifier {
  final ApiService _apiService;
  final StorageService _storageService;
  
  User? _user;
  bool _isLoading = false;
  bool _isAuthenticated = false;

  AuthProvider(this._apiService, this._storageService) {
    _loadUserFromStorage();
  }

  User? get user => _user;
  bool get isLoading => _isLoading;
  bool get isAuthenticated => _isAuthenticated;

  Future<void> _loadUserFromStorage() async {
    _isLoading = true;
    notifyListeners();

    try {
      final token = await _storageService.getToken();
      final userData = await _storageService.getUserData();

      if (token != null && userData != null) {
        _apiService.setToken(token);
        _user = User.fromJson(jsonDecode(userData));
        _isAuthenticated = true;
      }
    } catch (e) {
      print('Error loading user from storage: $e');
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<bool> login(String email, String password) async {
    _isLoading = true;
    notifyListeners();

    try {
      final response = await _apiService.login(email, password);
      
      _user = User.fromJson(response['user']);
      _isAuthenticated = true;
      
      await _storageService.saveToken(response['access_token']);
      await _storageService.saveRefreshToken(response['refresh_token']);
      await _storageService.saveUserData(jsonEncode(_user!.toJson()));
      
      _apiService.setToken(response['access_token']);
      
      return true;
    } catch (e) {
      print('Login error: $e');
      return false;
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<bool> register(String name, String email, String password) async {
    _isLoading = true;
    notifyListeners();

    try {
      final response = await _apiService.register(name, email, password);
      
      _user = User.fromJson(response['user']);
      _isAuthenticated = true;
      
      await _storageService.saveToken(response['access_token']);
      await _storageService.saveRefreshToken(response['refresh_token']);
      await _storageService.saveUserData(jsonEncode(_user!.toJson()));
      
      _apiService.setToken(response['access_token']);
      
      return true;
    } catch (e) {
      print('Register error: $e');
      return false;
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<Map<String, dynamic>?> generateQRCode() async {
    if (_user == null) return null;

    try {
      final response = await _apiService.generateQRCode(_user!.id);
      return response;
    } catch (e) {
      print('QR code generation error: $e');
      return null;
    }
  }

  Future<bool> verifyQRCode(String token) async {
    _isLoading = true;
    notifyListeners();

    try {
      final response = await _apiService.verifyQRCode(token);
      
      _user = User.fromJson(response['user']);
      _isAuthenticated = true;
      
      await _storageService.saveToken(response['access_token']);
      await _storageService.saveRefreshToken(response['refresh_token']);
      await _storageService.saveUserData(jsonEncode(_user!.toJson()));
      
      _apiService.setToken(response['access_token']);
      
      return true;
    } catch (e) {
      print('QR code verification error: $e');
      return false;
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<void> logout() async {
    _user = null;
    _isAuthenticated = false;
    
    await _storageService.deleteTokens();
    await _storageService.deleteUserData();
    
    _apiService.removeToken();
    
    notifyListeners();
  }

  Future<bool> updateProfile({String? name, String? image}) async {
    if (_user == null) return false;

    try {
      final updatedUser = await _apiService.updateProfile(name, image);
      _user = updatedUser;
      
      await _storageService.saveUserData(jsonEncode(_user!.toJson()));
      
      notifyListeners();
      return true;
    } catch (e) {
      print('Profile update error: $e');
      return false;
    }
  }
}