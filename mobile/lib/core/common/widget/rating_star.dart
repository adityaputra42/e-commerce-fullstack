import 'package:flutter/material.dart';
import 'package:mobile/core/theme/theme.dart';

class RatingStar extends StatelessWidget {
  final double voteAverage;
  final double starSize;
  final double fontSize;
  final int star;

  final MainAxisAlignment alignment;

  const RatingStar({
    super.key,
    this.star = 5,
    this.voteAverage = 0,
    this.starSize = 16,
    this.fontSize = 12,
    this.alignment = MainAxisAlignment.start,
  });

  @override
  Widget build(BuildContext context) {
    int n = (voteAverage / 2).round();
    List<Widget> widgets = [];
    widgets.add(
      Text(
        "${(voteAverage / 2).roundToDouble()}",
        style: AppFont.medium12.copyWith(fontSize: fontSize),
      ),
    );
    widgets.add(const SizedBox(width: 4));
    widgets.addAll(
      List.generate(
        star,
        (index) => Icon(
          index < n ? Icons.star_rounded : Icons.star_outline_rounded,
          color: AppColor.yellowColor,
          size: starSize,
        ),
      ),
    );
    return Row(mainAxisAlignment: alignment, children: widgets);
  }
}
