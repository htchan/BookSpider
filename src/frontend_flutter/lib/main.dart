import 'package:flutter/material.dart';
import './UI/mainPage.dart';
import './UI/sitePage.dart';
import './UI/searchPage.dart';
import './UI/randomPage.dart';
import './UI/bookPage.dart';

void main() {
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  // This widget is the root of your application.
  String url = 'http://192.168.128.146:9427';
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
        if (settings.name == '/') {
          return MaterialPageRoute(builder: (context) => MainPage(url: url,));
        }
        var uri = Uri.parse(settings.name);
        print(uri.pathSegments);
        if (uri.pathSegments.length == 2) {
          return MaterialPageRoute(builder: (context) => SitePage(
            url: url,
            siteName: uri.pathSegments[0]
          ));
        } else if (uri.pathSegments.indexOf('search') == 1) {
          return MaterialPageRoute(builder: (context) => SearchPage(
            url: url, 
            siteName: uri.pathSegments[0],
            title: uri.queryParameters['title'],
            writer: uri.queryParameters['writer']
          ));
        } else if (uri.pathSegments.indexOf('random') == 1) {
          return MaterialPageRoute(builder: (context) => RandomPage(
            url: url, 
            siteName: uri.pathSegments[0]
          ));
        }else if (uri.pathSegments.length == 3) {
          return MaterialPageRoute(builder: (context) => BookPage(
            url: url,
            siteName: uri.pathSegments[0],
            bookId: uri.pathSegments[1]
          ));
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