import 'package:flutter/material.dart';

import 'package:iconify_flutter_plus/iconify_flutter_plus.dart';
import 'package:iconify_flutter_plus/icons/ph.dart';
import 'package:iconify_flutter_plus/icons/material_symbols.dart';
import 'package:iconify_flutter_plus/icons/mdi.dart';

import '../../../../config/theme/app_color.dart';
import '../../../../config/theme/app_font.dart';

class CustomBottomNavbar extends StatelessWidget {
  const CustomBottomNavbar({super.key, this.selectedIndex, this.onTap});
  final int? selectedIndex;
  final Function(int index)? onTap;
  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      margin: EdgeInsets.all(16),
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(8),
        color: Theme.of(context).cardColor,
      ),
      child: BottomNavigationBar(
        type: BottomNavigationBarType.fixed,
        items: <BottomNavigationBarItem>[
          BottomNavigationBarItem(
            icon: Padding(
              padding: EdgeInsets.all(2),
              child: Iconify(
                Mdi.wallet_outline,
                color: Theme.of(context).hintColor,
                size: 22,
              ),
            ),
            activeIcon: Padding(
              padding: EdgeInsets.all(2),
              child: Iconify(
                Mdi.wallet_outline,
                color: AppColor.primaryColor,
                size: 22,
              ),
            ),
            label: 'Home',
          ),
          BottomNavigationBarItem(
            icon: Padding(
              padding: EdgeInsets.all(2),
              child: Iconify(
                Ph.swap_bold,
                color: Theme.of(context).hintColor,
                size: 22,
              ),
            ),
            activeIcon: Padding(
              padding: EdgeInsets.all(2),
              child: Iconify(
                Ph.swap_bold,
                color: AppColor.primaryColor,
                size: 22,
              ),
            ),
            label: 'Swap',
          ),
          BottomNavigationBarItem(
            icon: Padding(
              padding: EdgeInsets.all(2),
              child: Iconify(
                MaterialSymbols.widgets_outline_rounded,
                color: Theme.of(context).hintColor,
                size: 22,
              ),
            ),
            activeIcon: Padding(
              padding: EdgeInsets.all(2),
              child: Iconify(
                MaterialSymbols.widgets_outline_rounded,
                color: AppColor.primaryColor,
                size: 22,
              ),
            ),
            label: 'DApp',
          ),
          BottomNavigationBarItem(
            icon: Padding(
              padding: EdgeInsets.all(2),
              child: Icon(
                Icons.settings_outlined,
                color: Theme.of(context).hintColor,
                size: 22,
              ),
            ),
            activeIcon: Padding(
              padding: EdgeInsets.all(2),
              child: Icon(
                Icons.settings_outlined,
                color: AppColor.primaryColor,
                size: 22,
              ),
            ),
            label: 'Settings',
          ),
        ],
        elevation: 0,
        currentIndex: selectedIndex ?? 0,
        onTap: onTap,
        backgroundColor: Colors.transparent,
        selectedItemColor: AppColor.primaryColor,
        selectedLabelStyle: AppFont.medium12,
        unselectedItemColor: Theme.of(context).hintColor,
        showUnselectedLabels: true,
        unselectedLabelStyle: AppFont.reguler12,
      ),
    );
  }
}
