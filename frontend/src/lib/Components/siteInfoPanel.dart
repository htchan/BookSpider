import 'package:flutter/material.dart';

class SiteInfoPanel extends StatelessWidget {
  GlobalKey scaffoldKey;
  Map<String, dynamic> info;

  SiteInfoPanel(this.scaffoldKey, this.info);
  
  Widget _renderBookCount(Map<String, dynamic> info) {
    var bookCount = info['bookCount'];
    var errorCount = info['errorCount'];
    var totalCount = info['bookCount'] + info['errorCount'];
    var maxId = info['maxid'];
    Color totalCountColor = (totalCount == maxId) ? Colors.black : Colors.red;
    return RichText(
      textScaleFactor: Theme.of(scaffoldKey.currentContext).textTheme.bodyText1.fontSize / 14,
      text: TextSpan(
        style: TextStyle(color: Colors.black),
        children:[
          TextSpan(text: 'BookCount : '),
          TextSpan(text: '$totalCount, $maxId', style: TextStyle(color: totalCountColor)),
          TextSpan(text: '($bookCount + $errorCount)')
        ]
      )
    );
  }
  Widget _renderRecordCount(Map<String, dynamic> info) {
    var bookRecordCount = info['bookRecordCount'];
    var errorRecordCount = info['errorRecordCount'];
    var totalRecordCount = info['bookRecordCount'] + info['errorRecordCount'];
    return Text('TotalCount : $totalRecordCount ($bookRecordCount + $errorRecordCount)');
  }
  Widget _renderEndCount(Map<String, dynamic> info) {
    var endCount = info['endCount'];
    var endRecordCount = info['endRecordCount'];
    return Text('EndCount : $endCount ($endRecordCount)');
  }
  Widget _renderDownloadCount(Map<String, dynamic> info) {
    var downloadCount = info['downloadCount'];
    var downloadRecordCount = info['downloadRecordCount'];
    return Text('DownloadCount : $downloadCount ($downloadRecordCount)');
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        _renderBookCount(info),
        _renderRecordCount(info),
        _renderEndCount(info),
        _renderDownloadCount(info),
      ],
      crossAxisAlignment: CrossAxisAlignment.start,
    );
  }
}