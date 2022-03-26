import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

import '../Components/bookSearchBar.dart';
import '../Components/siteChartPanel.dart';
import '../Components/siteInfoPanel.dart';


class SitePage extends StatefulWidget{
  final String url, siteName;

  SitePage({Key key, this.url, this.siteName}) : super(key: key);

  @override
  _SitePageState createState() => _SitePageState(this.url, this.siteName);
}

class _SitePageState extends State<SitePage> with SingleTickerProviderStateMixin {
  final String siteName, url;
  Widget _chartPanel, _dataPanel;
  final GlobalKey scaffoldKey = GlobalKey();

  _SitePageState(this.url, this.siteName) {
    // call backend api
    String apiUrl = '$url/sites/$siteName';
    _chartPanel = Center(child: Text("Loading Chart"));
    _dataPanel = Center(child: Text("Loading Data"));
    http.get(Uri.parse(apiUrl))
    .then( (response) {
      if (response.statusCode != 404) {
        Map<String, dynamic> info = Map<String, dynamic>.from(jsonDecode(response.body));
        print("from response ${response.body}\n$info");
        setState((){
          _chartPanel = SiteChartPanel(scaffoldKey, info);
          _dataPanel = SiteInfoPanel(scaffoldKey, info);
        });
      } else {
        _chartPanel = _dataPanel = Center(
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

  @override
  Widget build(BuildContext context) {
    // show the content
    final PageController pageController = PageController( initialPage: 0 );
    return Scaffold(
      appBar: AppBar(title: Text(siteName)),
      key: scaffoldKey,
      body: Container(
        child: Column(
          children:[
            Container(
              height: MediaQuery.of(context).size.height * 0.5,
              child: PageView(
                children: [
                  _chartPanel, 
                  _dataPanel,
                ],
                controller: pageController,
              )
            ),
            BookSearchBar(scaffoldKey: scaffoldKey, siteName: siteName),
            // _renderRandomButton(),
          ],
        ),
        margin: EdgeInsets.symmetric(horizontal: 5.0),
      )
    );
  }
}