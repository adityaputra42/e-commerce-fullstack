part of '../../../../features/home/presentation/screen/home_screen.dart';

class ListPopulerProduct extends StatelessWidget {
  const ListPopulerProduct({super.key});

  @override
  Widget build(BuildContext context) {
    return SliverList.builder(
      itemBuilder: (context, index) {
        return CardPopulerProduct();
      },
      itemCount: 5,
    );
  }
}

class CardPopulerProduct extends StatelessWidget {
  const CardPopulerProduct({super.key});

  @override
  Widget build(BuildContext context) {
    return CardGeneral(
      height: 76,
      margin: EdgeInsets.only(bottom: 12, left: 16, right: 16),
      padding: EdgeInsets.all(4),
      child: Row(
        children: [
          CachedNetworkImage(
            width: 68,
            height: 68,
            imageUrl: "https://www.resellerdropship.com/public/media/uploads/products/21879.1.jpg",
            imageBuilder: (context, imageProvider) => Container(
              decoration: BoxDecoration(
                borderRadius: BorderRadius.circular(8),
                image: DecorationImage(image: imageProvider, fit: BoxFit.cover),
              ),
            ),
            placeholder: (context, url) => ShimmerLoading(radius: 8),
            errorWidget: (context, url, error) => Icon(Icons.error),
          ),
          width(12),
          Expanded(
            child: Padding(
              padding: const EdgeInsets.symmetric(vertical: 4),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text("Shiren Abaya", style: AppFont.medium14),
                            height(2),
                            Text(
                              "Abaya",
                              style: AppFont.reguler10.copyWith(color: Theme.of(context).hintColor),
                            ),
                          ],
                        ),
                      ),
                      RatingStar(voteAverage: 8.4),
                    ],
                  ),

                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    crossAxisAlignment: CrossAxisAlignment.end,
                    children: [
                      Text("\$120", style: AppFont.medium12),
                      Container(
                        padding: EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                        decoration: BoxDecoration(
                          borderRadius: BorderRadius.circular(6),
                          color: AppColor.primaryColor,
                        ),
                        child: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Text(
                              "Add to Cart",
                              style: AppFont.reguler10.copyWith(color: AppColor.darkText1),
                            ),
                            width(4),
                            Iconify(Mdi.plus, color: Colors.white, size: 12),
                          ],
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ),
          width(8),
        ],
      ),
    );
  }
}
