import 'package:fl_chart/fl_chart.dart';
import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';

const ERROR_KEY = "0", IN_PROGRESS_KEY = "1", END_KEY = "2";

class Site {
  final String name;
  final int bookCount,
      writerCount,
      bookUniqueCount,
      bookMaxID,
      bookLatestSuccessID,
      bookDownloadCount,
      statusErrorCount,
      statusInProgressCount,
      statusEndCount;

  Site.from(this.name, Map<String, dynamic> map)
      : this.bookCount = map['BookCount'] ?? 0,
        this.writerCount = map['WriterCount'] ?? 0,
        this.bookUniqueCount = map['UniqueBookCount'] ?? 0,
        this.bookMaxID = map['MaxBookID'] ?? 0,
        this.bookLatestSuccessID = map['LatestSuccessID'] ?? 0,
        this.bookDownloadCount = map['DownloadCount'] ?? 0,
        this.statusErrorCount = map['StatusCount'][ERROR_KEY] ?? 0,
        this.statusInProgressCount = map['StatusCount'][IN_PROGRESS_KEY] ?? 0,
        this.statusEndCount = map['StatusCount'][END_KEY] ?? 0;

  List<PieChartSectionData> get sections {
    return [
      PieChartSectionData(
        color: Colors.red,
        value: this.statusErrorCount.toDouble(),
        title: 'Error',
        radius: 50.0,
        titleStyle: TextStyle(fontWeight: FontWeight.bold, color: Colors.white),
      ),
      PieChartSectionData(
        color: Colors.yellow,
        value: this.statusInProgressCount.toDouble(),
        title: 'In Progress',
        radius: 50.0,
        titleStyle: TextStyle(fontWeight: FontWeight.bold, color: Colors.black),
      ),
      PieChartSectionData(
        color: Colors.blue,
        value: (this.statusEndCount - this.bookDownloadCount).toDouble(),
        title: 'End',
        radius: 50.0,
        titleStyle: TextStyle(fontWeight: FontWeight.bold, color: Colors.white),
      ),
      PieChartSectionData(
        color: Colors.green,
        value: this.bookDownloadCount.toDouble(),
        title: 'Download',
        radius: 50.0,
        titleStyle: TextStyle(fontWeight: FontWeight.bold, color: Colors.white),
      ),
    ];
  }
}
