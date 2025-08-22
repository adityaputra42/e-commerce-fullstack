import 'package:shared_preferences/shared_preferences.dart';

class PrefHelper {
  final SharedPreferences _pref;

  PrefHelper(this._pref);

  // Token
  Future<bool> saveToken(String token) async {
    return await _pref.setString('token', token);
  }

  String get token => _pref.getString('token') ?? '';

  Future<bool> removeToken() async {
    return await _pref.remove('token');
  }

  // First install
  Future<bool> setFirstInstall() async {
    return await _pref.setBool("firstInstall", false);
  }

  bool get isFirstInstall => _pref.getBool("firstInstall") ?? true;
}
