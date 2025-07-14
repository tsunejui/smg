import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:shared_preferences/shared_preferences.dart';

class StorageService {
  static const _secureStorage = FlutterSecureStorage();
  SharedPreferences? _prefs;

  Future<void> init() async {
    _prefs = await SharedPreferences.getInstance();
  }

  // Secure storage for sensitive data
  Future<void> saveToken(String token) async {
    await _secureStorage.write(key: 'access_token', value: token);
  }

  Future<String?> getToken() async {
    return await _secureStorage.read(key: 'access_token');
  }

  Future<void> saveRefreshToken(String token) async {
    await _secureStorage.write(key: 'refresh_token', value: token);
  }

  Future<String?> getRefreshToken() async {
    return await _secureStorage.read(key: 'refresh_token');
  }

  Future<void> deleteTokens() async {
    await _secureStorage.delete(key: 'access_token');
    await _secureStorage.delete(key: 'refresh_token');
  }

  // Regular storage for non-sensitive data
  Future<void> saveUserData(String userData) async {
    await _prefs?.setString('user_data', userData);
  }

  Future<String?> getUserData() async {
    return _prefs?.getString('user_data');
  }

  Future<void> deleteUserData() async {
    await _prefs?.remove('user_data');
  }

  Future<void> saveLastSync(DateTime dateTime) async {
    await _prefs?.setString('last_sync', dateTime.toIso8601String());
  }

  Future<DateTime?> getLastSync() async {
    final dateString = _prefs?.getString('last_sync');
    if (dateString != null) {
      return DateTime.parse(dateString);
    }
    return null;
  }

  Future<void> saveSetting(String key, String value) async {
    await _prefs?.setString(key, value);
  }

  Future<String?> getSetting(String key) async {
    return _prefs?.getString(key);
  }

  Future<void> clearAll() async {
    await _secureStorage.deleteAll();
    await _prefs?.clear();
  }
}