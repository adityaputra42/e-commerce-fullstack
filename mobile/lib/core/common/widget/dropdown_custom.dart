import 'package:dropdown_button2/dropdown_button2.dart';
import 'package:flutter/material.dart';

import '../../theme/theme.dart';
import '../../utils/size_extension.dart';

class DropDownCustom extends StatelessWidget {
  const DropDownCustom({
    super.key,
    required this.listData,
    this.title,
    required this.hint,
    this.value,
    this.onChange,
    this.color,
    this.hintStyle,
    this.textStyle,
    this.contentPadding,
    this.borderColor,
    this.borderRadius,
    this.filledColor,
    this.filled = true,
    this.heightWidget,
    this.prefixIcon,
    this.enable = true,
    this.prefix,
    this.suffix,
    this.icon,
  });
  final EdgeInsetsGeometry? contentPadding;
  final BorderRadius? borderRadius;
  final Color? borderColor;
  final Color? filledColor;
  final bool filled;
  final String hint;
  final String? title;
  final String? value;
  final double? heightWidget;
  final Widget? icon;
  final Widget? suffix;
  final Widget? prefix;
  final List<String> listData;
  final Function(String?)? onChange;
  final Color? color;
  final TextStyle? textStyle;
  final TextStyle? hintStyle;
  final Widget? prefixIcon;
  final bool enable;
  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      mainAxisSize: MainAxisSize.min,
      children: [
        title == null
            ? const SizedBox()
            : Column(
                children: [
                  Text(
                    title ?? '',
                    style: AppFont.medium12.copyWith(
                      color: Theme.of(context).colorScheme.onSurface,
                    ),
                  ),
                  height(8),
                ],
              ),
        SizedBox(
          height: heightWidget,
          child: DropdownButtonFormField2<String>(
            isExpanded: true,
            isDense: true,
            decoration: InputDecoration(
              enabled: enable,
              isDense: true,

              contentPadding: contentPadding ?? EdgeInsets.zero,
              suffixIcon: icon,
              suffix: suffix,
              prefixIcon: prefixIcon,
              prefixIconConstraints: BoxConstraints(
                maxHeight: 24,
                minHeight: 12,
                maxWidth: 72,
                minWidth: 36,
              ),
              suffixIconConstraints: BoxConstraints(
                maxHeight: 24,
                minHeight: 12,
                maxWidth: 72,
                minWidth: 36,
              ),
              prefix: prefix,
              hintText: hint,
              filled: filled,
              fillColor: filledColor ?? Theme.of(context).colorScheme.surface,
              hintStyle:
                  hintStyle ?? AppFont.reguler12.copyWith(color: Theme.of(context).hintColor),
              border: OutlineInputBorder(
                borderRadius: borderRadius ?? BorderRadius.circular(10),
                borderSide: BorderSide(
                  color: borderColor ?? Theme.of(context).colorScheme.outline,
                  width: 0.5,
                ),
              ),
              disabledBorder: OutlineInputBorder(
                borderRadius: borderRadius ?? BorderRadius.circular(10),
                borderSide: BorderSide(
                  color: borderColor ?? Theme.of(context).colorScheme.outline,
                  width: 0.5,
                ),
              ),
              enabledBorder: OutlineInputBorder(
                borderRadius: borderRadius ?? BorderRadius.circular(10),
                borderSide: BorderSide(
                  color: borderColor ?? Theme.of(context).colorScheme.outline,
                  width: 0.5,
                ),
              ),
              focusedBorder: OutlineInputBorder(
                borderRadius: borderRadius ?? BorderRadius.circular(10),
                borderSide: const BorderSide(color: AppColor.primaryColor),
              ),
            ),
            hint: Text(
              hint,
              style: AppFont.reguler14.copyWith(
                fontWeight: FontWeight.w300,
                color: AppColor.grayColor,
              ),
            ),
            items: listData
                .map(
                  (item) => DropdownMenuItem<String>(
                    value: item,
                    child: Text(
                      item,
                      style: AppFont.reguler12.copyWith(
                        color: Theme.of(context).colorScheme.onSurface,
                      ),
                    ),
                  ),
                )
                .toList(),
            validator: (value) {
              if (value == null) {
                return 'Pilih satuan.';
              }
              return null;
            },
            value: value,
            onChanged: onChange,
            buttonStyleData: ButtonStyleData(
              height: 42,
              padding: EdgeInsets.only(right: 8),
              decoration: BoxDecoration(color: Colors.transparent),
            ),

            iconStyleData: const IconStyleData(icon: Icon(Icons.expand_more), iconSize: 20),
            dropdownStyleData: DropdownStyleData(
              decoration: BoxDecoration(
                color: color ?? Theme.of(context).cardColor,
                borderRadius: BorderRadius.circular(12),
              ),
            ),
            menuItemStyleData: const MenuItemStyleData(
              height: 32,
              padding: EdgeInsets.only(left: 16, right: 12),
            ),
          ),
        ),
      ],
    );
  }
}
