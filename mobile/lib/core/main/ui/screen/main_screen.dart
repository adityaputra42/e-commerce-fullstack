import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:mobile/core/main/cubit/main_cubit.dart';
import 'package:mobile/core/main/ui/widget/custom_bottom_navbar.dart';

import '../../../../features/home/presentation/screen/home_screen.dart';

class MainScreen extends StatelessWidget {
  const MainScreen({super.key});

  @override
  Widget build(BuildContext context) {
    body(int selectedIndex) {
      switch (selectedIndex) {
        case 0:
          return const HomeScreen();
        case 1:
          return Center(child: Text("Cart"));
        case 2:
          return Center(child: Text("History"));

        default:
          return Center(child: Text("profile"));
      }
    }

    return Scaffold(
      body: SafeArea(
        child: BlocBuilder<MainCubit, int>(
          builder: (context, state) {
            return Stack(
              children: [
                body(state),
                Align(
                  alignment: Alignment.bottomCenter,
                  child: Padding(
                    padding: EdgeInsets.fromLTRB(24, 0, 24, 16),
                    child: CustomBottomNavbar(
                      selectedIndex: state,
                      onTap: (index) {
                        context.read<MainCubit>().setTab(index);
                      },
                    ),
                  ),
                ),
              ],
            );
          },
        ),
      ),
    );
  }
}
