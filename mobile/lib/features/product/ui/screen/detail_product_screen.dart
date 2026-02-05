import 'package:cached_network_image/cached_network_image.dart';
import 'package:carousel_slider/carousel_slider.dart';
import 'package:flutter/material.dart';
import 'package:mobile/core/utils/size_extension.dart';

import '../../../../core/common/widget/shimmer_loading.dart';

class DetailProductScreen extends StatelessWidget {
  const DetailProductScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Stack(
        children: [
          CarouselSlider.builder(
            options: CarouselOptions(
              autoPlay: false,
              disableCenter: true,
              viewportFraction: 1,
              aspectRatio: 8 / 8,
              autoPlayInterval: Duration(seconds: 10),
              autoPlayAnimationDuration: Duration(milliseconds: 1500),
              onPageChanged: (index, reason) {},
            ),
            itemCount: 3,
            itemBuilder: (context, index, _) {
              // if (banner.isEmpty) {
              //   return const ShimmerLoading(radius: 12);
              // }
              return CachedNetworkImage(
                imageUrl:
                    "https://www.resellerdropship.com/public/media/uploads/products/21879.1.jpg",
                imageBuilder: (context, imageProvider) => Container(
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(0),
                    image: DecorationImage(image: imageProvider, fit: BoxFit.cover),
                  ),
                ),
                placeholder: (context, url) => ShimmerLoading(radius: 0),
                errorWidget: (context, url, error) => Icon(Icons.error),
              );
            },
          ),
          SafeArea(
            child: SingleChildScrollView(
              child: Container(
                margin: EdgeInsets.only(top: context.h(0.38)),
                height: context.h(0.8),
                width: double.infinity,
                padding: EdgeInsets.all(16),
                decoration: BoxDecoration(
                  borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
                  color: Theme.of(context).cardColor,
                ),
                child: Column(),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
