import 'package:flutter/material.dart';
import 'package:url_strategy/url_strategy.dart';
import 'package:flutter_web_plugins/flutter_web_plugins.dart';

import './UI/mainPage.dart';
import './UI/sitePage.dart';
import './UI/searchPage.dart';
import './UI/randomPage.dart';
import './UI/bookPage.dart';
import './UI/stagePage.dart';
import './UI/errorPage.dart';

void main() {
  // setPathUrlStrategy();
  setUrlStrategy(PathUrlStrategy());
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  // This widget is the root of your application.
  String url = 'http://192.168.128.146/api/novel';
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
      initialRoute: '/',
      onGenerateRoute: (settings) {
        var uri = Uri.parse(settings.name);
        print(uri.pathSegments);
        if (uri.pathSegments.indexOf('stage') == 0) {
          return MaterialPageRoute(
            builder: (context) => StagePage(
              url: url,
            ), 
            settings: settings);
        } else if (uri.pathSegments.length >= 1 && uri.pathSegments.indexOf('sites') == 0) {
          return MaterialPageRoute(
            builder: (context) => SitePage(
              url: url,
              siteName: uri.pathSegments[1]
            ),
            settings: settings);
        } else if (uri.pathSegments.length >= 1 && uri.pathSegments.indexOf('search') == 0) {
          return MaterialPageRoute(
            builder: (context) => SearchPage(
              url: url, 
              siteName: uri.pathSegments[1],
              title: uri.queryParameters['title'],
              writer: uri.queryParameters['writer']
            ),
            settings: settings);
        } else if (uri.pathSegments.length >= 1 && uri.pathSegments.indexOf('random') == 0) {
          return MaterialPageRoute(
            builder: (context) => RandomPage(
              url: url, 
              siteName: uri.pathSegments[1]
            ),
            settings: settings);
        } else if (uri.pathSegments.length >= 2 && uri.pathSegments.indexOf('books') == 0) {
          return MaterialPageRoute(
            builder: (context) => BookPage(
              url: url,
              siteName: uri.pathSegments[1],
              bookId: uri.pathSegments[2]
            ),
            settings: settings);
        } else {
          return MaterialPageRoute(builder: (context) => MainPage(url: url,),
            settings: settings);
        }
      }
    );
  }
}

/*
http://host/                                            => main page
http://host/sites/<site>                                      => site page
http://host/search/<site>?title=<title>,writer=<writer> => search page
http://host/random/<site>                               => random page
http://host/books/<site>/<num>                                => book page
http://host/books/<site>/<num>/<version>                      => book page
*/