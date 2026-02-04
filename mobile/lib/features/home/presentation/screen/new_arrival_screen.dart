import 'package:flutter/material.dart';
import '../../../../core/common/widget/input_text.dart';
import '../../../../core/utils/size_extension.dart';
import '../../../../core/utils/widget_helper.dart';
import 'home_screen.dart';

class NewArrivalScreen extends StatelessWidget {
  const NewArrivalScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: WidgetHelper.appBar(context: context, title: "New Arrival"),
      body: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 16),
        child: Column(
          children: [
            height(8),
            InputText(
              hintText: "Search",
              filled: true,
              icon: Icon(Icons.search, size: 20),
              filledColor: Theme.of(context).cardColor,
            ),
            height(12),
            Expanded(
              child: GridView.builder(
                gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                  crossAxisCount: 2,
                  childAspectRatio: 0.7,
                  mainAxisSpacing: 12,
                  crossAxisSpacing: 12,
                ),
                itemBuilder: (context, index) => CardNewArrival(),
                itemCount: 12,
              ),
            ),
            height(16),
          ],
        ),
      ),
    );
  }
}
