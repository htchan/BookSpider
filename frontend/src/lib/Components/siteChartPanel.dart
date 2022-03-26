import 'package:flutter/material.dart';
import 'package:charts_flutter/flutter.dart' as charts;

class Data {
  final String name;
  final int value;
  Data(this.name, this.value);
}

class SiteChartPanel extends StatelessWidget {
  final GlobalKey scaffoldKey;
  List<Data> data;
  int maxId;


  SiteChartPanel(this.scaffoldKey, info) {
    var downloadCount = info['statusCount']['3'];
    var errorCount = info['errorCount'];
    var bookCount = info['bookCount'] - errorCount - downloadCount;
    data = [
      Data('Download', downloadCount),
      Data('Book', bookCount),
      Data('error', errorCount)
    ];
    maxId = info['latestSuccessBookId'];
  }

  List<charts.Series<Data, String>> _formatData() {
    return [
      charts.Series<Data, String>(
        id: 'DownloadData',
        domainFn: (Data data, _) => data.name,
        measureFn: (Data data, _) => data.value,
        data: data,
        // Set a label accessor to control the text of the arc label.
        labelAccessorFn: (Data row, _) => '${row.name}: ${row.value}',
      )
    ];
  }
  
  @override
  Widget build(BuildContext context) {
    return Stack(
      children: <Widget>[
        charts.PieChart(
          this._formatData(),
          animate: true,
          defaultRenderer: charts.ArcRendererConfig(
            arcWidth: (MediaQuery.of(scaffoldKey.currentContext).size.height / 8).round(),
            arcRendererDecorators: [new charts.ArcLabelDecorator()]),
        ),
        Center(child: Text(
          maxId.toString(),
          style: TextStyle(
            fontSize: 30.0,
            color: Colors.blue,
            fontWeight: FontWeight.bold
          )
        ))
      ],
    );
  }
}