import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:mobile/config/theme/theme.dart';
import 'package:mobile/core/common/cubit/theme_cubit.dart';
import 'package:mobile/core/splash/cubit/splash_cubit.dart';

import 'config/di/init_dependencies.dart';
import 'config/routes/app_routes.dart';
import 'core/main/cubit/main_cubit.dart';
import 'features/onboarding/cubit/onboarding_cubit.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await initDependencies();
  runApp(
    MultiBlocProvider(
      providers: [
        BlocProvider(create: (context) => serviceLocator<MainCubit>()),
        BlocProvider(
          create: (context) => serviceLocator<SplashCubit>()..initApp(),
        ),
        BlocProvider(create: (context) => serviceLocator<OnboardingCubit>()),
        BlocProvider(create: (context) => serviceLocator<OnboardingCubit>()),
        BlocProvider(
          create: (context) => serviceLocator<ThemeCubit>()..loadTheme(),
        ),
      ],
      child: const MyApp(),
    ),
  );
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});
  @override
  Widget build(BuildContext context) {
    return BlocBuilder<ThemeCubit, bool>(
      builder: (context, isDarkMode) {
        return MaterialApp.router(
          routerConfig: AppRouter.router,
          debugShowCheckedModeBanner: false,
          title: 'E-commerce App',
          theme: Styles.themeData(isDarkMode, context),
        );
      },
    );
  }
}
