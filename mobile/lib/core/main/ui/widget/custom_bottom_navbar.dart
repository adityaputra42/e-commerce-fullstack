import 'package:flutter/material.dart';

import 'package:iconify_flutter_plus/iconify_flutter_plus.dart';
import 'package:iconify_flutter_plus/icons/ph.dart';
import 'package:iconify_flutter_plus/icons/material_symbols.dart';
import 'package:iconify_flutter_plus/icons/mdi.dart';
import 'package:iconify_flutter_plus/icons/uil.dart';
import 'package:mobile/core/utils/size_extension.dart';

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
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 8),
      decoration: BoxDecoration(
        color: AppColor.secondaryColor,
        borderRadius: BorderRadius.circular(32),
        boxShadow: [
          BoxShadow(
            spreadRadius: 0.25,
            blurRadius: 0.5,
            color: Theme.of(
              context,
            ).colorScheme.onSurface.withValues(alpha: 0.1),
          ),
        ],
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          navbarItem(
            context,
            title: "Home",
            icon: Mdi.wallet_outline,
            isSelected: selectedIndex == 0,
            index: 0,
          ),
          navbarItem(
            context,
            title: "Swap",
            icon: Ph.swap_bold,
            isSelected: selectedIndex == 1,
            index: 1,
          ),
          navbarItem(
            context,
            title: "Dapp",
            icon: MaterialSymbols.widgets_outline_rounded,
            isSelected: selectedIndex == 2,
            index: 2,
          ),
          navbarItem(
            context,
            title: "Setting",
            icon: Uil.setting,
            isSelected: selectedIndex == 3,
            index: 3,
          ),
        ],
      ),
    );
  }

  InkWell navbarItem(
    BuildContext context, {
    required int index,
    required bool isSelected,
    required String icon,
    String? title,
  }) {
    return InkWell(
      onTap: () {
        if (onTap != null) {
          onTap!(index);
        }
      },
      borderRadius: BorderRadius.circular(24),
      child: AnimatedSize(
        alignment: isSelected ? Alignment.centerLeft : Alignment.center,
        duration: Duration(milliseconds: 300),
        curve: Curves.easeInOut,
        child: isSelected
            ? activeNavbar(context, icon: icon, title: title ?? "")
            : inactiveNavbar(context, icon: icon),
      ),
    );
  }

  Widget activeNavbar(
    BuildContext context, {
    required String icon,
    required String title,
  }) {
    return Container(
      height: 38,
      width: 84,
      padding: const EdgeInsets.all(4),
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(24),
        color: Theme.of(context).colorScheme.primary,
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            padding: const EdgeInsets.all(5),
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              border: Border.all(width: 1, color: AppColor.primaryColor),
              color: AppColor.darkText1,
            ),
            child: Iconify(icon, size: 18, color: AppColor.lightText1),
          ),
          width(4),
          Expanded(
            child: Text(
              title,
              style: AppFont.medium10.copyWith(color: AppColor.lightText1),
              overflow: TextOverflow.ellipsis,
            ),
          ),
          width(4),
        ],
      ),
    );
  }

  Widget inactiveNavbar(BuildContext context, {required String icon}) {
    return Container(
      width: 38,
      height: 38,
      padding: const EdgeInsets.all(9),
      decoration: BoxDecoration(
        shape: BoxShape.circle,
        color: Theme.of(
          context,
        ).colorScheme.onSurfaceVariant.withValues(alpha: 0.2),
      ),
      child: Iconify(icon, size: 20, color: AppColor.grayColor),
    );
  }
}
