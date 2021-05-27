import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

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
    http.get(apiUrl)
    .then( (response) {
      if (response.statusCode != 404) {
        setState((){
          _booksPanel = _renderBooksPenal(List<Map<String, dynamic>>.from(
            jsonDecode(response.body)['books'] ?? []
          ));
        });
      }
    });
  }

  Widget _renderLoadAbove() {
    return TextButton(
      child: Text('Load Above'),
      onPressed: () {
        setState(() { this._loadPage(--_page); });
      },
    );
  }

  Widget _renderLoadMore() {
    return TextButton(
      child: Text('Load More'),
      onPressed: () {
        setState(() { this._loadPage(++_page); });
        scrollController.animateTo(
          0,
          duration: Duration(milliseconds: 500),
          curve: Curves.fastOutSlowIn
        );
      },
    );
  }

  Widget _renderBooksPenal(List<Map<String, dynamic>> books) {
    if (books.length == 0) {
      return Center(child: Text('No books found'));
    }
    List<Widget> list = [];
    if (_page > 0) { list.add(_renderLoadAbove()); }
    list.addAll(books.map( (book) => ListTile(
        title: Text('${book['title']} - ${book['writer']}'),
        subtitle: Text('${book['update']} - ${book['chapter']}'),
        onTap: () {
          // go to book page
          Navigator.pushNamed(
            this.scaffoldKey.currentContext,
            '/books/$siteName/${book['id']}'
          );
        }
      )
    ));
    if (books.length == 20) { list.add(_renderLoadMore()); }
    return ListView.separated(
      controller: scrollController,
      separatorBuilder: (context, index) => Divider(height: 10,),
      itemCount: list.length,
      itemBuilder: (context, index) => list[index],
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