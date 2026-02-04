import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';
import 'package:iconify_flutter_plus/iconify_flutter_plus.dart';
import 'package:iconify_flutter_plus/icons/mdi.dart';
import 'package:mobile/core/common/widget/card_general.dart';
import 'package:mobile/core/common/widget/primary_button.dart';
import 'package:mobile/core/utils/widget_helper.dart';

import '../../../../core/common/widget/shimmer_loading.dart';
import '../../../../core/theme/theme.dart';
import '../../../../core/utils/size_extension.dart';

class CartScreen extends StatelessWidget {
  const CartScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: WidgetHelper.appBar(context: context, title: "Cart"),
      body: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 16),
        child: Column(
          children: [
            height(16),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text("Select All", style: AppFont.medium14),
                SizedBox(
                  width: 16,
                  height: 16,
                  child: Checkbox(
                    value: true,
                    onChanged: (value) {},
                    activeColor: AppColor.primaryColor,
                    checkColor: Colors.white,
                    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(4)),
                  ),
                ),
              ],
            ),
            height(16),
            Expanded(
              child: ListView.builder(
                itemBuilder: (context, index) {
                  return Padding(
                    padding: const EdgeInsets.only(bottom: 12.0),
                    child: CardCartProduct(),
                  );
                },
                itemCount: 5,
              ),
            ),
          ],
        ),
      ),
      bottomNavigationBar: SafeArea(
        child: CardGeneral(
          margin: EdgeInsets.zero,
          radius: 0,
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisSize: MainAxisSize.min,
                children: [
                  Text("Total Price", style: AppFont.reguler12),
                  height(4),
                  Text("\$600", style: AppFont.semibold16),
                ],
              ),
              PrimaryButton(width: context.w(0.4), title: "Beli Sekarang", onPressed: () {}),
            ],
          ),
        ),
      ),
    );
  }
}

class CardCartProduct extends StatefulWidget {
  const CardCartProduct({super.key});

  @override
  State<CardCartProduct> createState() => _CardCartProductState();
}

class _CardCartProductState extends State<CardCartProduct> {
  int indexCount = 0;
  void decreament() {
    setState(() {
      if (indexCount > 0) {
        indexCount -= 1;
      }
    });
  }

  void increment() {
    setState(() {
      indexCount += 1;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        SizedBox(
          width: 16,
          height: 16,
          child: Checkbox(
            value: true,
            onChanged: (value) {},
            activeColor: AppColor.primaryColor,
            checkColor: Colors.white,
            shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(4)),
          ),
        ),
        widget.width(12),
        Expanded(
          child: CardGeneral(
            height: 80,
            margin: EdgeInsets.zero,
            padding: EdgeInsets.all(4),
            child: Row(
              children: [
                CachedNetworkImage(
                  width: 72,
                  height: 72,
                  imageUrl:
                      "https://www.resellerdropship.com/public/media/uploads/products/21879.1.jpg",
                  imageBuilder: (context, imageProvider) => Container(
                    decoration: BoxDecoration(
                      borderRadius: BorderRadius.circular(8),
                      image: DecorationImage(image: imageProvider, fit: BoxFit.cover),
                    ),
                  ),
                  placeholder: (context, url) => ShimmerLoading(radius: 8),
                  errorWidget: (context, url, error) => Icon(Icons.error),
                ),
                widget.width(12),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text("Shiren Abaya", style: AppFont.medium14),
                      widget.height(2),
                      Text(
                        "Abaya",
                        style: AppFont.reguler10.copyWith(color: Theme.of(context).hintColor),
                      ),
                      widget.height(8),
                      Text("\$120", style: AppFont.medium12),
                    ],
                  ),
                ),
                widget.width(8),
                Row(
                  children: [
                    InkWell(
                      onTap: () {
                        decreament();
                      },
                      child: Container(
                        height: 24,
                        width: 24,
                        padding: EdgeInsets.all(4),
                        decoration: BoxDecoration(
                          border: BoxBorder.all(
                            color: indexCount > 0
                                ? AppColor.primaryColor
                                : Theme.of(context).hintColor,
                            width: 1,
                          ),
                          borderRadius: BorderRadius.circular(4),
                        ),
                        child: Iconify(
                          Mdi.minus,
                          color: indexCount > 0
                              ? AppColor.primaryColor
                              : Theme.of(context).hintColor,
                          size: 16,
                        ),
                      ),
                    ),
                    widget.width(4),
                    SizedBox(
                      width: 28,
                      child: Center(child: Text(indexCount.toString(), style: AppFont.medium14)),
                    ),
                    widget.width(4),
                    InkWell(
                      onTap: () {
                        increment();
                      },
                      child: Container(
                        height: 24,
                        width: 24,
                        padding: EdgeInsets.all(4),
                        decoration: BoxDecoration(
                          border: BoxBorder.all(color: AppColor.primaryColor, width: 1),
                          borderRadius: BorderRadius.circular(4),
                        ),
                        child: Iconify(Mdi.plus, color: AppColor.primaryColor, size: 16),
                      ),
                    ),
                  ],
                ),
                widget.width(8),
              ],
            ),
          ),
        ),
      ],
    );
  }
}
