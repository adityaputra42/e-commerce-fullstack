import 'package:internet_connection_checker_plus/internet_connection_checker_plus.dart';
import 'package:get_it/get_it.dart';
import 'package:mobile/core/common/cubit/theme_cubit.dart';
import 'package:mobile/features/main/cubit/main_cubit.dart';
import 'package:mobile/core/utils/pref_helper.dart';
import 'package:mobile/features/onboarding/cubit/onboarding_cubit.dart';

import 'features/splash/cubit/splash_cubit.dart';
import 'core/utils/connection_checker.dart';

final serviceLocator = GetIt.instance;

Future<void> initDependencies() async {
  PrefHelper.instance.init();
  serviceLocator.registerFactory(() => InternetConnection());
  serviceLocator.registerLazySingleton(() => MainCubit());
  serviceLocator.registerFactory<ConnectionChecker>(() => ConnectionCheckerImpl(serviceLocator()));
  serviceLocator.registerFactory<SplashCubit>(() => SplashCubit());
  serviceLocator.registerFactory<ThemeCubit>(() => ThemeCubit());
  serviceLocator.registerFactory<OnboardingCubit>(() => OnboardingCubit());
}
