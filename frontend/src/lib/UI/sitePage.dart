import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:charts_flutter/flutter.dart' as charts;

class SitePage extends StatefulWidget{
  final String url, siteName;

  SitePage({Key key, this.url, this.siteName}) : super(key: key);

  @override
  _SitePageState createState() => _SitePageState(this.url, this.siteName);
}

class Data {
  final String name;
  final int value;
  Data(this.name, this.value);
}

class _SitePageState extends State<SitePage> with SingleTickerProviderStateMixin {
  final String siteName, url;
  Widget _chartPanel, _dataPanel;
  final GlobalKey scaffoldKey = GlobalKey();
  TabController _tabController;

  _SitePageState(this.url, this.siteName) {
    // call backend api
    String apiUrl = '$url/sites/$siteName';
    _chartPanel = Center(child: Text("Loading Chart"));
    _dataPanel = Center(child: Text("Loading Data"));
    http.get(apiUrl)
    .then( (response) {
      if (response.statusCode != 404) {
        Map<String, dynamic> info = Map<String, dynamic>.from(jsonDecode(response.body));
        setState((){
          _chartPanel = _renderChartPanel(info);
          _dataPanel = _renderDataPanel(info);
        });
      } else {
        _chartPanel = _dataPanel = Center(
          child: Column(
            children: [
              Text(response.statusCode.toString()),
              Text(response.body)
            ],
          )
        );
      }
    });
    _tabController = TabController(length: 2, vsync: this);
  }

  Widget _renderBookCount(Map<String, dynamic> info) {
    var bookCount = info['bookCount'];
    var errorCount = info['errorCount'];
    var totalCount = info['bookCount'] + info['errorCount'];
    var maxId = info['maxid'];
    Color totalCountColor = (totalCount == maxId) ? Colors.black : Colors.red;
    return RichText(
      textScaleFactor: Theme.of(context).textTheme.bodyText1.fontSize / 14,
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

  Widget _renderDataPanel(Map<String, dynamic> info) {
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
  
  List<charts.Series<Data, String>> _formatData(Map<String, dynamic> info) {
    List<Data> data = [
      Data('Download', info['downloadCount']),
      Data('Book', info['bookCount'] - info['downloadCount']),
      Data('error', info['errorCount'])
    ];
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

  Widget _renderChartPanel(Map<String, dynamic> info) {
    return Stack(
      children: <Widget>[
        charts.PieChart(
          this._formatData(info),
          animate: true,
          defaultRenderer: charts.ArcRendererConfig(
            arcWidth: 100,
            arcRendererDecorators: [new charts.ArcLabelDecorator()]),
        ),
        Center(child: Text(
          info["maxid"].toString(),
          style: TextStyle(
            fontSize: 30.0,
            color: Colors.blue,
            fontWeight: FontWeight.bold
          )
        ))
      ],);
  }

  Widget _renderSearchPanel() {
    TextEditingController titleControler, writerController;
    titleControler = TextEditingController();
    writerController = TextEditingController();
    return Center(
      child: Card(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: <Widget>[
            TextField(
              decoration: InputDecoration(labelText: 'Book Title'),
              controller: titleControler
            ),
            TextField(
              decoration: InputDecoration(labelText: 'Book Writer'),
              controller: writerController,
            ),
            TextButton(
              child: Text('Submit'),
              onPressed: () {
                String title = titleControler.text;
                String writer = writerController.text;
                Navigator.pushNamed(
                  this.scaffoldKey.currentContext, 
                  '/search/$siteName?title=$title&writer=$writer'
                );
              },
            )
          ],
        ),
      ),
    );
  }
  Widget _renderRandomButton() {
    return RaisedButton(
      child: Text('Random'),
      onPressed: () {
        Navigator.pushNamed(
          this.scaffoldKey.currentContext,
          '/random/$siteName'
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    // show the content
    return Scaffold(
      appBar: AppBar(title: Text(this.siteName)),
      key: this.scaffoldKey,
      body: Container(
        child: Column(
          children:[
            Container(
              decoration: new BoxDecoration(color: Theme.of(context).primaryColor),
              child: TabBar(
                tabs: <Widget>[
                  Tab(text: 'Chart',), 
                  Tab(text: 'Data',),
                ],
                controller: this._tabController,)),
            Container(
              height: 500,
              child: TabBarView(
                children: [
                  _chartPanel, 
                  _dataPanel,
                ],
                controller: this._tabController,
              )
            ),
            this._renderSearchPanel(),
            Divider(),
            this._renderRandomButton(),
          ],
        ),
        margin: EdgeInsets.symmetric(horizontal: 5.0),
      )
    );
  }
}