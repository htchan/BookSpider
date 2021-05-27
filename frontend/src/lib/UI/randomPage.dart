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
  int n = 20;
  Widget _booksPanel;
  final GlobalKey scaffoldKey = GlobalKey();
  final ScrollController scrollController;

  _RandomPageState(this.url, this.siteName)
  : this.scrollController = ScrollController() {
    this._loadPage();
  }

  void _loadPage() {
    String apiUrl = '$url/random/$siteName?num=$n';
    _booksPanel = Center(child: Text("Loading books..."));
    http.get(apiUrl)
    .then( (response) {
      if (response.statusCode >= 200 && response.statusCode < 300) {
        setState((){
          _booksPanel = _renderBooksPanel(List<Map<String, dynamic>>.from(
            jsonDecode(response.body)['books'] ?? []
          ));
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

  Widget _renderRandomButton() {
    return ListTile(
      title: Center(child: Text(
        'Reload',
        style: TextStyle(color: Colors.blue)
      )),
      onTap: () {
        setState(() {
          this._loadPage();
        });
        scrollController.animateTo(0,
          duration: Duration(milliseconds: 500),
          curve: Curves.fastOutSlowIn
        );
      },
      hoverColor: Colors.blue.shade50,
    );
  }

  Widget _renderBooksPanel(List<Map<String, dynamic>> books) {
    if (books.length == 0) { return Center(child: Text('no books found')); }
    List<Widget> list = List<Widget>.from(books.map( (book) => ListTile(
      title: Text('${book['title']} - ${book['writer']}'),
      subtitle: Text('${book['update']} - ${book['chapter']}'),
      onTap: () {
        Navigator.pushNamed(
          this.scaffoldKey.currentContext,
          '/novel/books/$siteName/${book['id']}'
        );
    })));
    if (books.length == 20) { list.add(_renderRandomButton()); }
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