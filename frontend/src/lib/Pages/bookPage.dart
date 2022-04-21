import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:js' as js;
import 'package:url_launcher/url_launcher.dart';

class BookPage extends StatefulWidget{
  final String url, siteName, bookId;

  BookPage({Key key, this.url, this.siteName, this.bookId}) : super(key: key);

  @override
  _BookPageState createState() => _BookPageState(this.url, this.siteName, this.bookId);
}

class _BookPageState extends State<BookPage> {
  final String siteName, url, bookId;
  Widget _body;
  final GlobalKey scaffoldKey = GlobalKey();

  _BookPageState(this.url, this.siteName, this.bookId) {
    // call backend api
    String apiUrl = '$url/books/$siteName/$bookId';
    _body = Center(child: Text('Loading'));
    http.get(Uri.parse(apiUrl))
    .then( (response) {
      if (200 <= response.statusCode && response.statusCode < 300) {
        setState(() {
          _body = _renderBookContent(jsonDecode(response.body));
        });
      } else {
        _body = Center(
          child: Column(
            children: [
              Text(response.statusCode.toString()),
              Text(response.body)
            ],
          )
        );
      }
    });
  }
  
  Widget _renderBookContent(Map<String, dynamic> info) {
    List<Widget> rows = [
      SelectableText('ID: ${info['id']} - ${info['hash']}'),
      Row(
        // TODO: extract this to an external widget
        children: [
          SelectableText('Title: ${info['title']}'),
          ElevatedButton(
            child: Text("Search Google"),
            onPressed: () {
              // TODO: open new page to search google
              String searchUrl = "https://www.google.com/search?q=${info['title']}";
              print("${searchUrl}    ${canLaunch(searchUrl??"")}");
              canLaunch(searchUrl??"").then( (result) {
                if (result) launch(searchUrl??"");
              });
            },
          )
        ],
      ),
      Row(
        // TODO: extract this to external widget
        children: [
          SelectableText('Writer: ${info['writer']}'),
          ElevatedButton(
            child: Text("Search"),
            onPressed: () {
              // TODO: search the writer name internally
                String writer = info['writer'];
                Navigator.pushNamed(
                  this.scaffoldKey.currentContext,
                  '/search/$siteName?writer=$writer'
                );
            },
          )
        ],
      ),
      SelectableText('Type: ${info['type']}'),
      SelectableText('Last Update: ${info['updateDate']}'),
      SelectableText('Last Chapter: ${info['updateChapter']}'),
      SelectableText('Status: ${info['status']}')
    ];
    if (info['status'] == 'download') {
      rows.add(RaisedButton(
        child: Text('Download'),
        onPressed: () => js.context.callMethod('open', ['$url/download/$siteName/$bookId']),
      ));
    } else if (info['end']) {
      rows.add(Text('End'));
    }
    return ListView.separated(
      separatorBuilder: (context, index) => Divider(height: 10,),
      itemCount: rows.length,
      itemBuilder: (context, index) => rows[index],
    );
  }

  @override
  Widget build(BuildContext context) {
    // show the content
    return Scaffold(
      appBar: AppBar(title: Text(this.siteName)),
      key: this.scaffoldKey,
      body: Container(
        child: _body,
        margin: EdgeInsets.symmetric(horizontal: 5.0),
      )
    );
  }
}
