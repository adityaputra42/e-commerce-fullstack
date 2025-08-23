import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:mobile/core/main/cubit/main_cubit.dart';
import 'package:mobile/core/main/ui/widget/custom_bottom_navbar.dart';

class MainScreen extends StatelessWidget {
  const MainScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      bottomNavigationBar: BlocBuilder<MainCubit, int>(
        builder: (context, state) {
          return SafeArea(
            child: CustomBottomNavbar(
              selectedIndex: state,
              onTap: (index) {
                context.read<MainCubit>().setTab(index);
              },
            ),
          );
        },
      ),
    );
  }
}
