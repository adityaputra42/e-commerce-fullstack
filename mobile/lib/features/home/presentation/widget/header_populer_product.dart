part of '../../../../features/home/presentation/screen/home_screen.dart';

class HeaderPoppulerProduct extends StatelessWidget {
  const HeaderPoppulerProduct({super.key});

  @override
  Widget build(BuildContext context) {
    return SliverToBoxAdapter(
      child: Padding(
        padding: EdgeInsets.symmetric(horizontal: 16),
        child: Column(
          children: [
            height(16),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text("Populer Product", style: AppFont.medium14),
                InkWell(
                  onTap: () {
                    context.pushNamed(RouteNames.populerProduct);
                  },
                  child: Text(
                    "See All",
                    style: AppFont.reguler14.copyWith(color: AppColor.primaryColor),
                  ),
                ),
              ],
            ),
            height(12),
          ],
        ),
      ),
    );
  }
}
