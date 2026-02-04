import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:iconify_flutter_plus/iconify_flutter_plus.dart';
import 'package:iconify_flutter_plus/icons/mdi.dart';
import 'package:mobile/core/common/widget/input_text.dart';
import 'package:mobile/core/constants/constant.dart';

import '../../../../core/common/widget/custom_tab_bar.dart';
import '../../../../core/theme/theme.dart';
import '../../../../core/utils/size_extension.dart';
import '../../../../core/utils/widget_helper.dart';
import '../../../home/presentation/screen/home_screen.dart';

class ProductScreen extends StatelessWidget {
  const ProductScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: WidgetHelper.appBar(
        context: context,
        title: "Product",
        isCanBack: false,

        bottomWidet: Padding(
          padding: const EdgeInsets.only(top: 8),
          child: Row(
            children: [
              Expanded(
                child: InputText(
                  hintText: "Search Product",
                  filledColor: Theme.of(context).cardColor,
                ),
              ),
              width(8),
              InkWell(
                onTap: () {
                  context.pushNamed(RouteNames.cart);
                },
                child: Container(
                  width: 42,
                  height: 42,
                  padding: EdgeInsets.all(8),
                  decoration: BoxDecoration(
                    color: AppColor.primaryColor.withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(4),
                  ),
                  child: Center(
                    child: Iconify(Mdi.cart_outline, color: AppColor.primaryColor, size: 24),
                  ),
                ),
              ),
            ],
          ),
        ),
        height: 100,
      ),
      body: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 16),
        child: Column(
          children: [
            height(16),
            CustomTabBar(titles: ["All", "Electronics", "Fashion", "Home"], selectedIndex: 0),
            height(8),
            Expanded(
              child: GridView.builder(
                padding: EdgeInsets.only(bottom: 120),
                gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                  crossAxisCount: 2,
                  childAspectRatio: 0.7,
                  mainAxisSpacing: 8,
                  crossAxisSpacing: 8,
                ),
                itemBuilder: (context, index) => CardNewArrival(),
                itemCount: 12,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
