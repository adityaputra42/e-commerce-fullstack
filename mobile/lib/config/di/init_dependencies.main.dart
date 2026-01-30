part of 'init_dependencies.dart';

final serviceLocator = GetIt.instance;

Future<void> initDependencies() async {
  PrefHelper.instance.init();
  serviceLocator.registerFactory(() => InternetConnection());

  // core
  serviceLocator.registerLazySingleton(() => MainCubit());
  serviceLocator.registerFactory<ConnectionChecker>(
    () => ConnectionCheckerImpl(serviceLocator()),
  );
  serviceLocator.registerFactory<SplashCubit>(
    () => SplashCubit(),
  );

  serviceLocator.registerFactory<ThemeCubit>(
    () => ThemeCubit(),
  );
  serviceLocator.registerFactory<OnboardingCubit>(() => OnboardingCubit());
}
