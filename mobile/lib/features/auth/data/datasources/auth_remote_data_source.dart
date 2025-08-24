import 'package:mobile/features/auth/data/models/user_model.dart';

import '../../../../core/error/exceptions.dart';

abstract interface class AuthRemoteDataSource {
  // Session? get currentUserSession;
  Future<UserModel> signUpWithEmailPassword({
    required String name,
    required String email,
    required String password,
  });
  Future<UserModel> loginWithEmailPassword({
    required String email,
    required String password,
  });
  Future<UserModel?> getCurrentUserData();
}

class AuthRemoteDataSourceImpl implements AuthRemoteDataSource {
  // final SupabaseClient supabaseClient;
  AuthRemoteDataSourceImpl();

  // @override
  // Session? get currentUserSession => supabaseClient.auth.currentSession;

  @override
  Future<UserModel> loginWithEmailPassword({
    required String email,
    required String password,
  }) async {
    try {
      return UserModel(email: email, id: "id", name: "name");
    } catch (e) {
      throw ServerException(e.toString());
    }
  }

  @override
  Future<UserModel> signUpWithEmailPassword({
    required String name,
    required String email,
    required String password,
  }) async {
    try {
      return UserModel(email: email, id: "id", name: name);
    } catch (e) {
      throw ServerException(e.toString());
    }
  }

  @override
  Future<UserModel?> getCurrentUserData() async {
    try {
      return null;
    } catch (e) {
      throw ServerException(e.toString());
    }
  }
}
