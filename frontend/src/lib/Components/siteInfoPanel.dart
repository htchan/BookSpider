import 'package:bookspider/models/all_model.dart';
import 'package:flutter/material.dart';

class SiteInfoPanel extends StatelessWidget {
  final GlobalKey scaffoldKey;
  final Site site;

  SiteInfoPanel(this.scaffoldKey, this.site);

  Widget _renderBookCount() {
    var totalCount = site.statusErrorCount +
        site.statusInProgressCount +
        site.statusEndCount;
    Color totalCountColor =
        (totalCount == site.bookCount) ? Colors.black : Colors.red;
    return RichText(
        textScaleFactor: Theme.of(scaffoldKey.currentContext!)
                .textTheme
                .bodyText1!
                .fontSize! /
            14,
        text: TextSpan(style: TextStyle(color: Colors.black), children: [
          TextSpan(text: 'BookCount : '),
          TextSpan(
              text: '$totalCount, ${site.bookCount}',
              style: TextStyle(color: totalCountColor)),
          TextSpan(
              text:
                  '(${site.statusErrorCount} + ${site.statusInProgressCount} + ${site.statusEndCount})')
        ]));
  }

  Widget _renderRecordCount() {
    return Text(
        'TotalCount : ${site.bookCount + site.statusErrorCount} (${site.bookCount} + ${site.statusErrorCount})');
  }

  Widget _renderEndCount() {
    return Text('EndCount : ${site.statusEndCount}');
  }

  Widget _renderDownloadCount() {
    return Text('DownloadCount : ${site.bookDownloadCount}');
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        _renderBookCount(),
        _renderRecordCount(),
        _renderEndCount(),
        _renderDownloadCount(),
      ],
      crossAxisAlignment: CrossAxisAlignment.start,
    );
  }
}
