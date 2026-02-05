import 'package:flutter/material.dart';
import 'package:shimmer/shimmer.dart';

class ShimmerLoading extends StatelessWidget {
  const ShimmerLoading({
    super.key,
    this.margin,
    this.height,
    this.width = double.infinity,
    this.radius = 8,
  });
  final EdgeInsetsGeometry? margin;
  final double? height;
  final double? width;
  final double radius;

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: margin,
      child: Shimmer.fromColors(
        baseColor: Theme.of(context).canvasColor,
        highlightColor: Theme.of(context).highlightColor,
        child: Container(
          width: width,
          height: height,
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(radius),
            color: Theme.of(context).canvasColor,
          ),
        ),
      ),
    );
  }
}

class ShimmerLoadingCircle extends StatelessWidget {
  const ShimmerLoadingCircle({super.key, this.margin, this.height, this.width = double.infinity});
  final EdgeInsetsGeometry? margin;
  final double? height;
  final double? width;

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: margin,
      child: Shimmer.fromColors(
        baseColor: Theme.of(context).canvasColor,
        highlightColor: Theme.of(context).highlightColor,
        child: Container(
          width: width,
          height: height,
          decoration: BoxDecoration(shape: BoxShape.circle, color: Theme.of(context).canvasColor),
        ),
      ),
    );
  }
}
