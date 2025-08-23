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
        color: Theme.of(context).cardColor,
        borderRadius: BorderRadius.circular(8),
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
        children: [
          Expanded(
            child: InkWell(
              onTap: () {
                if (onTap != null) {
                  onTap!(0);
                }
              },
              child: (selectedIndex == 0)
                  ? activeNavbar(
                      context,
                      icon: Mdi.wallet_outline,
                      title: "Home",
                    )
                  : inactiveNavbar(context, icon: Mdi.wallet_outline),
            ),
          ),
          width(2),
          Expanded(
            child: InkWell(
              onTap: () {
                if (onTap != null) {
                  onTap!(1);
                }
              },
              child: (selectedIndex == 1)
                  ? activeNavbar(context, icon: Ph.swap_bold, title: "Swap")
                  : inactiveNavbar(context, icon: Ph.swap_bold),
            ),
          ),
          width(2),
          Expanded(
            child: InkWell(
              onTap: () {
                if (onTap != null) {
                  onTap!(2);
                }
              },
              child: (selectedIndex == 2)
                  ? activeNavbar(
                      context,
                      icon: MaterialSymbols.widgets_outline_rounded,
                      title: "DApp",
                    )
                  : inactiveNavbar(
                      context,
                      icon: MaterialSymbols.widgets_outline_rounded,
                    ),
            ),
          ),
          width(2),
          Expanded(
            child: InkWell(
              onTap: () {
                if (onTap != null) {
                  onTap!(3);
                }
              },
              child: (selectedIndex == 3)
                  ? activeNavbar(context, icon: Uil.setting, title: "Setting")
                  : inactiveNavbar(context, icon: Uil.setting),
            ),
          ),
        ],
      ),
    );
  }

  Widget activeNavbar(
    BuildContext context, {
    required String icon,
    required String title,
  }) {
    return Container(
      padding: const EdgeInsets.all(4),
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(24),
        color: Theme.of(context).colorScheme.surface,
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            padding: const EdgeInsets.all(5),
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              border: Border.all(width: 1, color: AppColor.primaryColor),
              color: Theme.of(context).colorScheme.surface,
            ),
            child: Iconify(
              icon,
              size: 18,
              color: Theme.of(context).colorScheme.onSurface,
            ),
          ),
          width(4),
          Expanded(
            child: Text(
              title,
              style: AppFont.medium10.copyWith(
                color: Theme.of(context).colorScheme.onSurface,
              ),
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
      padding: const EdgeInsets.all(8),
      decoration: BoxDecoration(
        shape: BoxShape.circle,
        color: Theme.of(context).colorScheme.surface,
      ),
      child: Iconify(icon, size: 20, color: AppColor.grayColor),
    );
  }
}
