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
  bool error = true;
  int page = 0;
  List<Map<String, dynamic>> books;
  final GlobalKey scaffoldKey = GlobalKey();
  final ScrollController scrollController;

  _SearchPageState(this.url, this.siteName, this.title, this.writer)
  : this.scrollController = ScrollController() {
    this._loadPage(this.page);
  }
  void _loadPage(int page) {
    String apiUrl = '$url/search/$siteName?title=$title&writer=$writer&page=$page';
    http.get(apiUrl)
    .then( (response) {
      if (response.statusCode != 404) {
        this.books = List<Map<String, dynamic>>.from(jsonDecode(response.body)['books']);
        this.error = false;
        setState((){});
      }
    });
  }
  List<Widget> _renderBooks() {
    List<Widget> buttons = [];
    if (this.error) {
      return [Center(child: Text('Loading Books...'))];
    }
    if (this.page > 0) {
      buttons.add(TextButton(
        child: Text('Load Above'),
        onPressed: () {
          this.error = true;
          setState(() {});
          this._loadPage(--this.page);
        },
      ));
    }
    for (Map<String, dynamic> book in this.books) {
      buttons.add(ListTile(
        title: Text(book['title'] + ' - ' + book['writer']),
        subtitle: Text(book['update'] + ' - ' + book['chapter']),
        onTap: () {
          // go to book page
          Navigator.pushNamed(
            this.scaffoldKey.currentContext,
            '/books/$siteName/${book['id']}'
          );
        }
      ));
    }
    if (this.books.length == 20) {
      buttons.add(TextButton(
        child: Text('Load More'),
        onPressed: () {
          this.error = true;
          setState(() {});
          this._loadPage(++this.page);
          scrollController.animateTo(
            0,
            duration: Duration(milliseconds: 500),
            curve: Curves.fastOutSlowIn
          );
        },
      ));
    }
    if (this.books.length == 0) {
      buttons.add(Center(child: Text('No books found')));
    }
    return buttons;
  }

  @override
  Widget build(BuildContext context) {
    // show the content
    List<Widget> buttons = this._renderBooks();
    return Scaffold(
      appBar: AppBar(title: Text(this.siteName)),
      key: this.scaffoldKey,
      body: Container(
        child: ListView.separated(
          controller: scrollController,
          separatorBuilder: (context, index) => Divider(height: 10,),
          itemCount: buttons.length,
          itemBuilder: (context, index) => buttons[index],
          
        ),
        margin: EdgeInsets.symmetric(horizontal: 5.0),
      ),
    );
  }
}