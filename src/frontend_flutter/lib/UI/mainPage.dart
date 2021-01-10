import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

class MainPage extends StatefulWidget{
  final String url;

  MainPage({Key key, this.url}) : super(key: key);

  @override
  _MainPageState createState() => _MainPageState(this.url);
}

class _MainPageState extends State<MainPage> {
  final String url;
  bool error = true;
  Map<String, dynamic> info;
  final GlobalKey scaffoldKey = GlobalKey();

  _MainPageState(this.url) {
    // call backend api
    String apiUrl = '$url/info';
    http.get(apiUrl)
    .then( (response) {
      if (response.statusCode != 404) {
        this.info = Map<String, dynamic>.from(jsonDecode(response.body));
        this.error = false;
        setState((){});
      }
    });
  }
  List<Widget> _renderSiteButton() {
    List<String> siteNames = (this.error) ? [] : List<String>.from(this.info['siteNames']);
    List<Widget> buttons = [RaisedButton(
      child: Text('Stage'),
      color: Colors.lightBlueAccent,
      onPressed: () {
        Navigator.pushNamed(
          this.scaffoldKey.currentContext, 
          '/stage/'
        );
      }
    )];
    for (String name in siteNames) {
      buttons.add(RaisedButton(
        child: Text(name),
        onPressed: () {
          // redirect to site page with its name
          Navigator.pushNamed(
            this.scaffoldKey.currentContext,
            '/$name/'
          );
        },
      ));
    }
    return buttons;
  }

  @override
  Widget build(BuildContext context) {
    // show the content
    List<Widget> buttons = this._renderSiteButton();
    return Scaffold(
      appBar: AppBar(title: Text('Book')),
      key: this.scaffoldKey,
      body: Container(
        child: ListView.separated(
          separatorBuilder: (context, index) => Divider(height: 10,),
          itemCount: buttons.length,
          itemBuilder: (context, index) => buttons[index],
          
        ),
        margin: EdgeInsets.symmetric(horizontal: 5.0),
      ),
    );
  }
}