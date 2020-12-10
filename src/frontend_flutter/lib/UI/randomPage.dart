import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

class RandomPage extends StatefulWidget{
  final String url, siteName;

  RandomPage({Key key, this.url, this.siteName}) : super(key: key);

  @override
  _RandomPageState createState() => _RandomPageState(this.url, this.siteName);
}

class _RandomPageState extends State<RandomPage> {
  final String url, siteName;
  bool error = true;
  int n = 20;
  List<Map<String, dynamic>> books;
  final GlobalKey scaffoldKey = GlobalKey();
  final ScrollController scrollController;

  _RandomPageState(this.url, this.siteName)
  : this.scrollController = ScrollController() {
    this._loadPage();
  }
  void _loadPage() {
    String apiUrl = '$url/random/$siteName?num=$n';
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
      return [Text('loading...')];
    }
    for (Map<String, dynamic> book in this.books) {
      buttons.add(ListTile(
        title: Text(book['title'] + ' - ' + book['writer']),
        subtitle: Text(book['update'] + ' - ' + book['chapter']),
        onTap: () {
          // go to book page
          Navigator.pushNamed(
            this.scaffoldKey.currentContext,
            '/$siteName/${book['num']}/'
          );
        }
      ));
    }
    if (this.books.length == 20) {
      buttons.add(TextButton(
        child: Text('Reload'),
        onPressed: () {
          this._loadPage();
          scrollController.animateTo(
            0,
            duration: Duration(milliseconds: 500),
            curve: Curves.fastOutSlowIn
          );
        },
      ));
    }
    if (this.books.length == 0) {
      buttons.add(Text('No books found'));
    }
    return buttons;
  }

  @override
  Widget build(BuildContext context) {
    // show the content
    List<Widget> buttons = this._renderBooks();
    return Scaffold(
      appBar: AppBar(title: Text('Book')),
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