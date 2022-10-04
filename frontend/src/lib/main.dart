import 'package:bookspider/models/all_model.dart';
import 'package:bookspider/repostory/bookSpiderRepostory.dart';
import 'package:flutter/material.dart';
import 'package:flutter_web_plugins/flutter_web_plugins.dart';

import './Pages/allPage.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';

void main() async {
  await dotenv.load(fileName: ".env");
  // setPathUrlStrategy();
  setUrlStrategy(PathUrlStrategy());
  runApp(MyApp());
}

const String host = String.fromEnvironment("NOVEL_SPIDER_API_HOST");
const String FE_ROUTE_PREFIX = String.fromEnvironment(
    "NOVEL_SPIDER_FE_ROUTE_PREFIX",
    defaultValue: "/novel");

class MyApp extends StatelessWidget {
  // This widget is the root of your application.
  final BookSpiderRepostory client = BookSpiderRepostory(host);
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
        initialRoute: "/",
        onGenerateRoute: (settings) {
          var uri = Uri.parse(settings.name ?? "");
          print("path ${uri.path}");
          if (RegExp("^/sites/([^/]*)/?\$").hasMatch(uri.path)) {
            return MaterialPageRoute(
                builder: (context) => SitePage(
                    client: client,
                    siteName: uri.pathSegments[1],
                    site: settings.arguments as Site?),
                settings: settings);
          } else if (RegExp("^/sites/([^/]*)/random/?\$").hasMatch(uri.path)) {
            var query = uri.queryParameters;
            return MaterialPageRoute(
                builder: (context) => RandomPage(
                      client: client,
                      siteName: uri.pathSegments[1],
                    ),
                settings: settings);
          } else if (RegExp("^/sites/([^/]*)/search/?\$").hasMatch(uri.path)) {
            var query = uri.queryParameters;
            return MaterialPageRoute(
                builder: (context) => SearchPage(
                      client: client,
                      siteName: uri.pathSegments[1],
                      title: query['title'] ?? "",
                      writer: query['writer'] ?? "",
                      page: int.parse(query['page'] ?? "0"),
                      perPage: int.parse(query['per_page'] ?? "20"),
                    ),
                settings: settings);
          } else if (RegExp("^/sites/([^/]*)/books/([^/]*)/?\$")
              .hasMatch(uri.path)) {
            var idHash = uri.pathSegments[3].split("-");
            idHash.add("0");
            return MaterialPageRoute(
                builder: (context) => BookPage(
                    client: client,
                    siteName: uri.pathSegments[1],
                    id: idHash[0],
                    hash: idHash[1],
                    book: settings.arguments as Book?),
                settings: settings);
          }
          return MaterialPageRoute(
              builder: (context) => MainPage(client: client),
              settings: settings);
        });
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