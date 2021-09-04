import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import '../Components/bookList.dart';

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
    http.get(Uri.parse(apiUrl))
    .then( (response) {
      if (response.statusCode >= 200 && response.statusCode < 300) {
        setState((){
          _booksPanel = BookList(
            scaffoldKey, 
            siteName,
            List<Map<String, dynamic>>.from(
              jsonDecode(response.body)['books'] ?? []
            ),
            null,
            randomButton
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

  Widget randomButton(ScrollController controller) {
    return ListTile(
      title: Center(child: Text(
        'Reload',
        style: TextStyle(color: Colors.blue)
      )),
      onTap: () {
        setState(() {
          this._loadPage();
        });
        controller.animateTo(0,
          duration: Duration(milliseconds: 500),
          curve: Curves.fastOutSlowIn
        );
      },
      hoverColor: Colors.blue.shade50,
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