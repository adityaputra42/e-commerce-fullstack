part of '../../../../features/home/presentation/screen/home_screen.dart';

class AppBarHome extends StatelessWidget {
  const AppBarHome({super.key});

  @override
  Widget build(BuildContext context) {
    return SliverAppBar(
      surfaceTintColor: Theme.of(context).colorScheme.surface,
      pinned: true,
      floating: false,
      snap: false,
      toolbarHeight: 60,
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
          InkWell(
            onTap: () {
              context.pushNamed(RouteNames.cart);
            },
            child: Container(
              padding: EdgeInsets.all(6),
              decoration: BoxDecoration(
                color: AppColor.primaryColor.withValues(alpha: 0.1),
                borderRadius: BorderRadius.circular(4),
              ),
              child: Iconify(Mdi.cart_outline, color: AppColor.primaryColor, size: 20),
            ),
          ),
        ],
      ),
    );
  }
}
