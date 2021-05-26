import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:js' as js;

class BookPage extends StatefulWidget{
  final String url, siteName, bookId;

  BookPage({Key key, this.url, this.siteName, this.bookId}) : super(key: key);

  @override
  _BookPageState createState() => _BookPageState(this.url, this.siteName, this.bookId);
}

class _BookPageState extends State<BookPage> {
  final String siteName, url, bookId;
  bool error = true;
  Map<String, dynamic> info;
  final GlobalKey scaffoldKey = GlobalKey();

  _BookPageState(this.url, this.siteName, this.bookId) {
    // call backend api
    String apiUrl = '$url/info/$siteName/$bookId';
    http.get(apiUrl)
    .then( (response) {
      if (response.statusCode != 404) {
        this.info = Map<String, dynamic>.from(jsonDecode(response.body));
        this.error = false;
        setState((){});
      }
    });
  }
  List<Widget> _renderBookContent() {
    if (this.error) {
      return [Text('Loading')];
    }
    String version = (this.error) ? '' : this.info['version'].toString();
    String title = (this.error) ? '' : this.info['title'];
    String writer = (this.error) ? '' : this.info['writer'];
    String type = (this.error) ? '' : this.info['type'];
    String lastUpdate = (this.error) ? '' : this.info['update'];
    String lastChapter = (this.error) ? '' : this.info['chapter'];
    List<Widget> rows = [];
    rows.addAll([
      Text('ID: $bookId - $version'),
      Text('Title: $title'),
      Text('Writer: $writer'),
      Text('Type: $type'),
      Text('Last Update: $lastUpdate'),
      Text('Last Chapter: $lastChapter')
    ]);
    if (this.info['download']) {
      rows.add(RaisedButton(
        child: Text('Download'),
        onPressed: () => js.context.callMethod('open', ['$url/download/$siteName/$bookId']),
      ));
    } else if (this.info['end']) {
      rows.add(Text('End'));
    }
    return rows;
  }

  @override
  Widget build(BuildContext context) {
    // show the content
    return Scaffold(
      appBar: AppBar(title: Text(this.siteName)),
      key: this.scaffoldKey,
      body: Container(
        child: ListView(
          children: this._renderBookContent()
        ),
        margin: EdgeInsets.symmetric(horizontal: 5.0),
      )
    );
  }
}