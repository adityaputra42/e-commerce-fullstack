import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:mobile/core/theme/app_font.dart';
import 'package:mobile/core/constants/constant.dart';

import '../cubit/splash_cubit.dart';

class SplashScreen extends StatelessWidget {
  const SplashScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocListener<SplashCubit, SplashState>(
      listener: (context, state) {
        if (state is SplashToOnboarding) {
          context.go('/${RouteNames.main}');
        } else if (state is SplashToHome) {
          context.go('/${RouteNames.main}');
        }
      },
      child: Scaffold(
        body: Center(child: Text('Splash Screen', style: AppFont.semibold16)),
      ),
    );
  }
}
