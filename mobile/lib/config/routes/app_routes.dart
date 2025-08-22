import 'package:go_router/go_router.dart';
import 'package:mobile/core/main/ui/screen/main_screen.dart';

import '../../core/constants/constant.dart';
import '../../core/splash/ui/splash_screen.dart';
import '../../features/onboarding/ui/screen/onboarding_screen.dart';

class AppRouter {
  static final router = GoRouter(
    initialLocation: '/',
    routes: [
      GoRoute(
        path: '/',
        name: RouteNames.splash,
        builder: (context, state) => const SplashScreen(),
      ),

      GoRoute(
        path: '/onboarding',
        name: RouteNames.onboarding,
        builder: (context, state) => const OnboardingScreen(),
      ),
      GoRoute(
        path: '/main',
        name: RouteNames.main,
        builder: (context, state) => const MainScreen(),
      ),
    ],
  );
}
