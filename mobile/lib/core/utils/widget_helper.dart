import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../theme/app_font.dart';

class WidgetHelper {
  static appBar({
    required BuildContext context,
    Function()? onTap,
    Function()? onTapTitle,
    required String title,
    Widget? icon,
    Widget? iconTitle,
    Widget? bottomWidet,
    Color? color,
    double? fontSize,
    double height = 60,
    bool isCanBack = true,
  }) {
    return AppBar(
      elevation: 0.1,
      title: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Row(
            children: [
              isCanBack
                  ? InkWell(
                      onTap:
                          onTap ??
                          () {
                            context.pop();
                          },
                      child: Icon(
                        Icons.arrow_back_ios_new_rounded,
                        color: Theme.of(context).colorScheme.onSurface,
                        size: 20,
                      ),
                    )
                  : const SizedBox(),
              SizedBox(width: isCanBack ? 12 : 0),
              Expanded(
                child: InkWell(
                  onTap: onTapTitle ?? () {},
                  child: iconTitle != null
                      ? Row(
                          mainAxisSize: MainAxisSize.min,
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Text(
                              title,
                              style: AppFont.medium16.copyWith(fontSize: fontSize ?? 16),

                              overflow: TextOverflow.ellipsis,
                            ),
                            iconTitle,
                          ],
                        )
                      : Text(
                          title,
                          style: AppFont.medium16.copyWith(fontSize: fontSize ?? 16),

                          overflow: TextOverflow.ellipsis,
                        ),
                ),
              ),
              const SizedBox(width: 16),
              icon ?? const SizedBox(),
            ],
          ),
          bottomWidet ?? const SizedBox(),
        ],
      ),
      automaticallyImplyLeading: false,
      backgroundColor: color ?? Theme.of(context).colorScheme.surface,
      surfaceTintColor: color ?? Theme.of(context).colorScheme.surface,
      shadowColor: Theme.of(context).canvasColor,
      toolbarHeight: height,
    );
  }
}
