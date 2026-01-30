import 'package:go_router/go_router.dart';
import 'package:mobile/core/main/ui/screen/main_screen.dart';
import 'package:mobile/features/auth/presentation/screen/sign_in_screen.dart';

import '../../core/constants/constant.dart';
import '../../features/splash/ui/splash_screen.dart';
import '../../features/auth/presentation/screen/sign_up_screen.dart';
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
        path: '/${RouteNames.onboarding}',
        name: RouteNames.onboarding,
        builder: (context, state) => const OnboardingScreen(),
      ),
      GoRoute(
        path: '/${RouteNames.main}',
        name: RouteNames.main,
        builder: (context, state) => const MainScreen(),
      ),
      GoRoute(
        path: '/${RouteNames.signin}',
        name: RouteNames.signin,
        builder: (context, state) => const SignInScreen(),
        routes: [
          GoRoute(
            path: '/${RouteNames.signup}',
            name: RouteNames.signup,
            builder: (context, state) => const SignUpScreen(),
          ),
        ],
      ),
    ],
  );
}
