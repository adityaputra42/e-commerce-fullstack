part of '../../../../features/home/presentation/screen/home_screen.dart';

class NewArrivalWidget extends StatelessWidget {
  const NewArrivalWidget({super.key});

  @override
  Widget build(BuildContext context) {
    return SliverToBoxAdapter(
      child: Column(
        children: [
          Padding(
            padding: EdgeInsets.symmetric(horizontal: 16),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text("New Arrival", style: AppFont.medium14),
                InkWell(
                  onTap: () {
                    context.pushNamed(RouteNames.newArrival);
                  },
                  child: Text(
                    "See All",
                    style: AppFont.reguler14.copyWith(color: AppColor.primaryColor),
                  ),
                ),
              ],
            ),
          ),
          height(16),
          SizedBox(
            height: 244,
            child: ListView.builder(
              scrollDirection: Axis.horizontal,
              itemBuilder: (context, index) {
                return Padding(
                  padding: EdgeInsets.only(left: index == 0 ? 16 : 0, right: 16),
                  child: CardNewArrival(),
                );
              },
              itemCount: 10,
            ),
          ),
        ],
      ),
    );
  }
}

class CardNewArrival extends StatelessWidget {
  const CardNewArrival({super.key});

  @override
  Widget build(BuildContext context) {
    return CardGeneral(
      radius: 12,
      width: 176,
      margin: EdgeInsets.symmetric(vertical: 1),
      padding: EdgeInsets.all(4),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          CachedNetworkImage(
            width: 168,
            height: 160,
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

          Padding(
            padding: const EdgeInsets.all(8),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text("Shiren Abaya", style: AppFont.medium14),
                height(2),
                Text(
                  "Abaya",
                  style: AppFont.reguler10.copyWith(color: Theme.of(context).hintColor),
                ),
                height(2),
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    Text("\$120", style: AppFont.medium12),
                    Container(
                      padding: EdgeInsets.all(2),
                      decoration: BoxDecoration(
                        color: AppColor.primaryColor,
                        borderRadius: BorderRadius.circular(4),
                      ),
                      child: Center(child: Iconify(Mdi.plus, color: Colors.white, size: 16)),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
