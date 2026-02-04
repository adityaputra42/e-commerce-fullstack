import 'package:flutter/material.dart';
import 'package:mobile/core/common/widget/custom_tab_bar.dart';

import '../../../../core/utils/widget_helper.dart';
import 'home_screen.dart';

class PopulerProductScreen extends StatelessWidget {
  const PopulerProductScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: WidgetHelper.appBar(context: context, title: "Popular Products"),
      body: NestedScrollView(
        headerSliverBuilder: (context, inerBox) {
          return [
            SliverToBoxAdapter(
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                child: CustomTabBar(
                  titles: ["All", "Electronics", "Fashion", "Home"],
                  selectedIndex: 0,
                ),
              ),
            ),
          ];
        },
        body: ListView.builder(
          itemBuilder: (context, index) => CardPopulerProduct(),
          itemCount: 10,
        ),
      ),
    );
  }
}
