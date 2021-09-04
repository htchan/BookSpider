import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import '../Components/bookList.dart';

class SearchPage extends StatefulWidget{
  final String url, siteName, title, writer;

  SearchPage({Key key, this.url, this.siteName, this.title, this.writer}) : super(key: key);

  @override
  _SearchPageState createState() => _SearchPageState(this.url, this.siteName, this.title, this.writer);
}

class _SearchPageState extends State<SearchPage> {
  final String url, siteName, title, writer;
  int _page = 0;
  Widget _booksPanel;
  final GlobalKey scaffoldKey = GlobalKey();
  final ScrollController scrollController;

  _SearchPageState(this.url, this.siteName, this.title, this.writer)
  : this.scrollController = ScrollController() {
    this._loadPage(_page);
  }
  void _loadPage(int page) {
    String apiUrl = '$url/search/$siteName?title=$title&writer=$writer&page=$page';
    _booksPanel = Center(child: Text('Loading Books...'));
    http.get(Uri.parse(apiUrl))
    .then( (response) {
      if (response.statusCode != 404) {
        setState((){
          _booksPanel = BookList(
            scaffoldKey,
            siteName,
            List<Map<String, dynamic>>.from(
              jsonDecode(response.body)['books'] ?? []
            ),
            _page > 0 ? loadAboveButton : null,
            loadMoreButton
          );
        });
      } else {
        _booksPanel = Center(
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

  Widget loadAboveButton(ScrollController controller) {
    return TextButton(
      child: Text('Load Above'),
      onPressed: () {
        setState(() { this._loadPage(--_page); });
        controller.animateTo(
          0,
          duration: Duration(milliseconds: 500),
          curve: Curves.fastOutSlowIn
        );
      },
    );
  }

  Widget loadMoreButton(ScrollController controller) {
    return TextButton(
      child: Text('Load More'),
      onPressed: () {
        setState(() { this._loadPage(++_page); });
        controller.animateTo(
          controller.position.maxScrollExtent,
          duration: Duration(milliseconds: 500),
          curve: Curves.fastOutSlowIn
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
        child: _booksPanel,
        margin: EdgeInsets.symmetric(horizontal: 5.0),
      ),
    );
  }
}