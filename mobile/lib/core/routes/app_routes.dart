import 'package:go_router/go_router.dart';
import 'package:mobile/features/main/ui/screen/main_screen.dart';
import 'package:mobile/features/auth/presentation/screen/sign_in_screen.dart';
import 'package:mobile/features/home/presentation/screen/new_arrival_screen.dart';
import 'package:mobile/features/home/presentation/screen/populer_product_screen.dart';

import '../../features/cart/ui/screen/cart_screen.dart';
import '../constants/constant.dart';
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
        routes: [
          GoRoute(
            path: '/${RouteNames.cart}',
            name: RouteNames.cart,
            builder: (context, state) => const CartScreen(),
          ),
          GoRoute(
            path: '/${RouteNames.newArrival}',
            name: RouteNames.newArrival,
            builder: (context, state) => const NewArrivalScreen(),
          ),
          GoRoute(
            path: '/${RouteNames.populerProduct}',
            name: RouteNames.populerProduct,
            builder: (context, state) => const PopulerProductScreen(),
          ),
        ],
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
