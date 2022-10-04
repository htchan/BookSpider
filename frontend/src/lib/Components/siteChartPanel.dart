import 'package:bookspider/models/all_model.dart';
import 'package:flutter/material.dart';
import 'package:fl_chart/fl_chart.dart';

class SiteChartPanel extends StatelessWidget {
  final GlobalKey scaffoldKey;
  final Site site;

  SiteChartPanel(this.scaffoldKey, this.site);

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: <Widget>[
        PieChart(PieChartData(
            borderData: FlBorderData(
              show: false,
            ),
            sectionsSpace: 0,
            centerSpaceRadius: 100,
            sections: site.sections)),
        Center(
            child: Text(site.bookMaxID.toString(),
                style: TextStyle(
                    fontSize: 30.0,
                    color: Colors.blue,
                    fontWeight: FontWeight.bold)))
      ],
    );
  }
}
