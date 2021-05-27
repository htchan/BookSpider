import 'package:flutter/material.dart';
import 'package:url_strategy/url_strategy.dart';

import './UI/mainPage.dart';
import './UI/sitePage.dart';
import './UI/searchPage.dart';
import './UI/randomPage.dart';
import './UI/bookPage.dart';
import './UI/stagePage.dart';

void main() {
  setPathUrlStrategy();
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  // This widget is the root of your application.
  String url = 'http://192.168.128.146:9427/api/novel';
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Book',
      theme: ThemeData(
        textTheme: Theme.of(context).textTheme.apply(
          fontSizeFactor: 1.25,
        ),
        primarySwatch: Colors.blue,
        visualDensity: VisualDensity.adaptivePlatformDensity,
      ),
      initialRoute: '/novel',
      onGenerateRoute: (settings) {
        if (settings.name == '/novel') {
          return MaterialPageRoute(builder: (context) => MainPage(url: url,),
            settings: settings);
        }
        var uri = Uri.parse(settings.name);
        print(uri.pathSegments);
        if (uri.pathSegments.indexOf('stage') == 1) {
          return MaterialPageRoute(
            builder: (context) => StagePage(
              url: url,
            ), 
            settings: settings);
        } else if (uri.pathSegments.length >= 2 && uri.pathSegments.indexOf('sites') == 1) {
          return MaterialPageRoute(
            builder: (context) => SitePage(
              url: url,
              siteName: uri.pathSegments[2]
            ),
            settings: settings);
        } else if (uri.pathSegments.length >= 2 && uri.pathSegments.indexOf('search') == 1) {
          return MaterialPageRoute(
            builder: (context) => SearchPage(
              url: url, 
              siteName: uri.pathSegments[2],
              title: uri.queryParameters['title'],
              writer: uri.queryParameters['writer']
            ),
            settings: settings);
        } else if (uri.pathSegments.length >= 2 && uri.pathSegments.indexOf('random') == 1) {
          return MaterialPageRoute(
            builder: (context) => RandomPage(
              url: url, 
              siteName: uri.pathSegments[2]
            ),
            settings: settings);
        } else if (uri.pathSegments.length >= 3 && uri.pathSegments.indexOf('books') == 1) {
          return MaterialPageRoute(
            builder: (context) => BookPage(
              url: url,
              siteName: uri.pathSegments[2],
              bookId: uri.pathSegments[3]
            ),
            settings: settings);
        }
      }
    );
  }
}

/*
http://host/                                            => main page
http://host/<site>                                      => site page
http://host/<site>/search?title=<title>,writer=<writer> => search page
http://host/<site>/random                               => random page
http://host/<site>/<num>                                => book page
http://host/<site>/<num>/<version>                      => book page
*/