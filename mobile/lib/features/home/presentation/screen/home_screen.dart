import 'package:flutter/material.dart';
import 'package:iconify_flutter_plus/iconify_flutter_plus.dart';
import 'package:iconify_flutter_plus/icons/ic.dart';
import 'package:iconify_flutter_plus/icons/mdi.dart';
import 'package:mobile/config/theme/theme.dart';

import '../../../../core/utils/size_extension.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: CustomScrollView(
        slivers: [
          SliverAppBar(
            pinned: true,
            floating: false,
            snap: false,
            toolbarHeight: 60,
            backgroundColor: AppColor.secondaryColor,
            title: Row(
              children: [
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text("Hi, Anisa", style: AppFont.reguler12),
                      height(2),
                      Text("Discover your style", style: AppFont.semibold20),
                    ],
                  ),
                ),
                Container(
                  padding: EdgeInsets.all(6),
                  decoration: BoxDecoration(
                    color: AppColor.primaryColor.withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(4),
                  ),
                  child: Iconify(Mdi.notifications_none, color: AppColor.primaryColor, size: 20),
                ),
                width(8),
                Container(
                  padding: EdgeInsets.all(6),
                  decoration: BoxDecoration(
                    color: AppColor.primaryColor.withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(4),
                  ),
                  child: Iconify(Ic.outline_message, color: AppColor.primaryColor, size: 20),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
