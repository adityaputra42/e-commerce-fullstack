import 'package:cached_network_image/cached_network_image.dart';
import 'package:carousel_slider/carousel_slider.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:iconify_flutter_plus/iconify_flutter_plus.dart';
import 'package:iconify_flutter_plus/icons/mdi.dart';
import 'package:mobile/core/theme/theme.dart';
import 'package:mobile/core/common/widget/card_general.dart';
import 'package:mobile/core/common/widget/rating_star.dart';
import 'package:mobile/core/common/widget/shimmer_loading.dart';

import '../../../../core/constants/constant.dart';
import '../../../../core/utils/size_extension.dart';

part '../../../../features/home/presentation/widget/app_bar_home.dart';
part '../../../../features/home/presentation/widget/banner_home.dart';
part '../../../../features/home/presentation/widget/new_arrival.dart';
part '../../../../features/home/presentation/widget/list_populer_product.dart';
part '../widget/header_populer_product.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: CustomScrollView(
        slivers: [
          AppBarHome(),
          BannerHome(),
          NewArrivalWidget(),
          HeaderPoppulerProduct(),
          ListPopulerProduct(),
          SliverPadding(padding: EdgeInsets.only(bottom: 120)),
        ],
      ),
    );
  }
}
