part of '../../../../features/home/presentation/screen/home_screen.dart';

class BannerHome extends StatefulWidget {
  const BannerHome({super.key});

  @override
  State<BannerHome> createState() => _BannerHomeState();
}

class _BannerHomeState extends State<BannerHome> {
  List<String> banner = [
    "https://bunny-wp-pullzone-skqqsftbkd.b-cdn.net/wp-content/uploads/2026/01/1.-Banner-Blog-copy-645x355.jpg",
    "https://bunny-wp-pullzone-skqqsftbkd.b-cdn.net/wp-content/uploads/2026/01/1.-Banner-Blog-738-copy-1.webp",
    "https://bunny-wp-pullzone-skqqsftbkd.b-cdn.net/wp-content/uploads/2024/10/Desain-blog-1-4.jpg",
  ];

  int currentIndex = 0;
  @override
  Widget build(BuildContext context) {
    return SliverToBoxAdapter(
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        child: Column(
          children: [
            CarouselSlider.builder(
              options: CarouselOptions(
                autoPlay: true,
                disableCenter: true,
                viewportFraction: 1,
                aspectRatio: 18 / 9,
                autoPlayInterval: Duration(seconds: 10),
                autoPlayAnimationDuration: Duration(milliseconds: 1500),
                onPageChanged: (index, reason) {
                  setState(() {
                    currentIndex = index;
                  });
                },
              ),
              itemCount: banner.length,
              itemBuilder: (context, index, _) {
                if (banner.isEmpty) {
                  return const ShimmerLoading(radius: 12);
                }
                return CachedNetworkImage(
                  imageUrl: banner[index],
                  imageBuilder: (context, imageProvider) => Container(
                    decoration: BoxDecoration(
                      borderRadius: BorderRadius.circular(12),
                      image: DecorationImage(image: imageProvider, fit: BoxFit.cover),
                    ),
                  ),
                  placeholder: (context, url) => ShimmerLoading(radius: 12),
                  errorWidget: (context, url, error) => Icon(Icons.error),
                );
              },
            ),

            widget.height(12),
            _CarouselIndicator(
              pageLength: banner.length,
              pageIndex: currentIndex,
              carouselDuration: Duration(seconds: 10),
            ),
          ],
        ),
      ),
    );
  }
}

class _CarouselIndicator extends StatelessWidget {
  final Duration carouselDuration;
  final int pageIndex;
  final int pageLength;

  const _CarouselIndicator({
    required this.pageIndex,
    required this.pageLength,
    required this.carouselDuration,
  });

  final double activeLength = 36;
  final double inactiveLength = 10;

  final double borderRadius = 999;

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: List.generate(pageLength, (index) {
        final bool isActive = pageIndex == index;
        return AnimatedContainer(
          duration: const Duration(milliseconds: 400),
          width: isActive ? activeLength : inactiveLength,
          height: inactiveLength,
          margin: const EdgeInsets.symmetric(horizontal: 2),
          decoration: BoxDecoration(
            color: AppColor.secondaryColor.withValues(alpha: isActive ? 0.2 : 0.3),
            borderRadius: BorderRadius.circular(borderRadius),
          ),
          child: Align(
            alignment: Alignment.centerLeft,
            child: AnimatedContainer(
              width: isActive ? activeLength : inactiveLength,
              height: inactiveLength,
              duration: carouselDuration,
              decoration: BoxDecoration(borderRadius: BorderRadius.circular(borderRadius)),
              child: Container(
                decoration: BoxDecoration(
                  color: isActive ? AppColor.primaryColor : Colors.transparent,
                  borderRadius: BorderRadius.circular(borderRadius),
                ),
              ),
            ),
          ),
        );
      }),
    );
  }
}
