import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

class SitePage extends StatefulWidget{
  final String url, siteName;

  SitePage({Key key, this.url, this.siteName}) : super(key: key);

  @override
  _SitePageState createState() => _SitePageState(this.url, this.siteName);
}

class _SitePageState extends State<SitePage> {
  final String siteName, url;
  bool error = true;
  Map<String, dynamic> info;
  final GlobalKey scaffoldKey = GlobalKey();

  _SitePageState(this.url, this.siteName) {
    // call backend api
    String apiUrl = '$url/info/$siteName';
    http.get(apiUrl)
    .then( (response) {
      if (response.statusCode != 404) {
        this.info = Map<String, dynamic>.from(jsonDecode(response.body));
        this.error = false;
        setState((){});
      }
    });
  }
  Widget _renderBookCount(BuildContext context) {
    String bookCount = (this.error) ? '' : this.info['bookCount'].toString();
    String errorCount = (this.error) ? '' : this.info['errorCount'].toString();
    String totalCount = (this.error) ? '' : (this.info['bookCount'] + this.info['errorCount']).toString();
    String maxId = (this.error) ? '' : this.info['maxid'].toString();
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
    String bookRecordCount = (this.error) ? '' : this.info['bookRecordCount'].toString();
    String errorRecordCount = (this.error) ? '' : this.info['errorRecordCount'].toString();
    String totalRecordCount = (this.error) ? '' : (this.info['bookRecordCount'] + this.info['errorRecordCount']).toString();
    return Text('TotalCount : $totalRecordCount ($bookRecordCount + $errorRecordCount)');
  }
  Widget _renderEndCount() {
    String endCount = (this.error) ? '' : this.info['endCount'].toString();
    String endRecordCount = (this.error) ? '' : this.info['endRecordCount'].toString();
    return Text('EndCount : $endCount ($endRecordCount)');
  }
  Widget _renderDownloadCount() {
    String downloadCount = (this.error) ? '' : this.info['downloadCount'].toString();
    String downloadRecordCount = (this.error) ? '' : this.info['downloadRecordCount'].toString();
    return Text('DownloadCount : $downloadCount ($downloadRecordCount)');
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
                  '/$siteName/search/?title=$title&writer=$writer'
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
          '/$siteName/random/'
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
        child: ListView(
          children: [
            this._renderBookCount(context),
            this._renderRecordCount(),
            this._renderEndCount(),
            this._renderDownloadCount(),
            this._renderSearchPanel(),
            Divider(),
            this._renderRandomButton()
          ],
        ),
        margin: EdgeInsets.symmetric(horizontal: 5.0),
      )
    );
  }
}