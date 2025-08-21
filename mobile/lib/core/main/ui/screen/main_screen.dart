import 'package:flutter/material.dart';
import 'package:mobile/core/main/ui/widget/custom_bottom_navbar.dart';

class MainScreen extends StatelessWidget {
  const MainScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      bottomNavigationBar: CustomBottomNavbar(selectedIndex: 0, onTap: null),
    );
  }
}
