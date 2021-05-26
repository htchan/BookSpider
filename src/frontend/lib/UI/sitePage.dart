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
  bool load = false;
  Map<String, dynamic> info;
  final GlobalKey scaffoldKey = GlobalKey();
  TabController _tabController;

  _SitePageState(this.url, this.siteName) {
    // call backend api
    String apiUrl = '$url/info/$siteName';
    http.get(apiUrl)
    .then( (response) {
      if (response.statusCode != 404) {
        this.info = Map<String, dynamic>.from(jsonDecode(response.body));
        this.load = true;
        setState((){});
      }
    });
    _tabController = TabController(length: 2, vsync: this);
  }
  Widget _renderBookCount() {
    String bookCount = (this.load) ? this.info['bookCount'].toString() : '';
    String errorCount = (this.load) ? this.info['errorCount'].toString() : '';
    String totalCount = (this.load) ? (int.parse(this.info['bookCount']) + int.parse(this.info['errorCount'])).toString() : '';
    String maxId = (this.load) ? this.info['maxid'].toString() : '';
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
  Widget _renderRecordCount() {
    String bookRecordCount = (this.load) ? this.info['bookRecordCount'].toString() : '';
    String errorRecordCount = (this.load) ? this.info['errorRecordCount'].toString() : '';
    String totalRecordCount = (this.load) ? (int.parse(this.info['bookRecordCount']) + int.parse(this.info['errorRecordCount'])).toString() : '';
    return Text('TotalCount : $totalRecordCount ($bookRecordCount + $errorRecordCount)');
  }
  Widget _renderEndCount() {
    String endCount = (this.load) ? this.info['endCount'].toString() : '';
    String endRecordCount = (this.load) ? this.info['endRecordCount'].toString() : '';
    return Text('EndCount : $endCount ($endRecordCount)');
  }
  Widget _renderDownloadCount() {
    String downloadCount = (this.load) ? this.info['downloadCount'].toString() : '';
    String downloadRecordCount = (this.load) ? this.info['downloadRecordCount'].toString() : '';
    return Text('DownloadCount : $downloadCount ($downloadRecordCount)');
  }
  List<charts.Series<Data, String>> _formatData() {
    if (!this.load) return [];

    List<Data> data = [
      Data('Download', int.parse(this.info['downloadCount'])),
      Data('Book', int.parse(this.info['bookCount']) - int.parse(this.info['downloadCount'])),
      Data('error', int.parse(this.info['errorCount']))
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
  Widget _renderData() {
    if (!this.load) return Center(child: Text("Loading Data"));
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
  Widget _renderChart() {
    if (!this.load) return Center(child: Text("Loading Chart"),);

    return Stack(
      children: <Widget>[
        charts.PieChart(
          this._formatData(),
          animate: true,
          defaultRenderer: charts.ArcRendererConfig(
            arcWidth: 100,
            arcRendererDecorators: [new charts.ArcLabelDecorator()]),
        ),
        Center(child: Text(
          (this.load) ? this.info["maxid"].toString() : '',
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
                  Tab(text: 'Chart',), Tab(text: 'Data',)
                ],
                controller: this._tabController,)),
            Container(
              height: 500,
              child: TabBarView(
                children: [
                  this._renderChart(),
                  this._renderData()
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