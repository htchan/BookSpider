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
  List<Widget> _buttons;
  final GlobalKey scaffoldKey = GlobalKey();

  _MainPageState(this.url) {
    String apiUrl = '$url/info';
    _buttons = _renderStageButton();
    http.get(Uri.parse(apiUrl))
    .then( (response) {
      if (response.statusCode >= 200 && response.statusCode < 300) {
        setState(() {
          _buttons.addAll(_renderSiteButtons(
            Map<String, dynamic>.from(jsonDecode(response.body))
          ));
        });
      } else {
        setState(() {
          _buttons.addAll([ 
            Text(response.statusCode.toString()),
            Text(response.body)
          ]);
        });
      }
    });
  }

  Iterable<Widget> _renderSiteButtons(Map<String, dynamic> info) {
    List<String> siteNames = List<String>.from(info['siteNames'] ?? []);
    return siteNames.map( (name) => RaisedButton(
      child: Text(name),
      onPressed: () {
        Navigator.pushNamed(
          scaffoldKey.currentContext,
          '/sites/$name'
        );
      },
    ));
  }

  List<Widget> _renderStageButton() {
    List<Widget> buttons = [RaisedButton(
      child: Text('Stage'),
      color: Colors.lightBlueAccent,
      onPressed: () {
        Navigator.pushNamed(
          this.scaffoldKey.currentContext, 
          '/stage'
        );
      }
    )];
    return buttons;
  }

  @override
  Widget build(BuildContext context) {
    // show the content
    return Scaffold(
      appBar: AppBar(title: Text('Book')),
      key: this.scaffoldKey,
      body: Container(
        child: ListView.separated(
          separatorBuilder: (context, index) => Divider(height: 10,),
          itemCount: _buttons.length,
          itemBuilder: (context, index) => _buttons[index],
        ),
        margin: EdgeInsets.symmetric(horizontal: 5.0),
      ),
    );
  }
}