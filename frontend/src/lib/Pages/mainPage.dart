import 'package:bookspider/models/all_model.dart';
import 'package:bookspider/repostory/bookSpiderRepostory.dart';
import 'package:flutter/material.dart';

class MainPage extends StatefulWidget {
  BookSpiderRepostory client;

  MainPage({Key? key, required this.client}) : super(key: key);

  @override
  _MainPageState createState() => _MainPageState(client);
}

class _MainPageState extends State<MainPage> {
  bool isError = false;
  List<Site> sites = [];
  final GlobalKey scaffoldKey = GlobalKey();
  final BookSpiderRepostory client;

  _MainPageState(this.client) {
    this.client.getSites().then((sites) {
      setState(() {
        this.sites = sites;
      });
    }).catchError((e) {
      setState(() {
        this.isError = true;
      });
    });
  }

  List<String> get siteNames {
    return List.from(this.sites.map((site) => site.name));
  }

  Iterable<Widget> _renderSiteButtons() {
    if (this.isError) {
      return [Text("Fail to query site")];
    }
    if (this.siteNames.isEmpty) {
      return [Center(child: Text("Loading"))];
    }
    return sites.map((site) => TextButton(
          child: Text(site.name),
          onPressed: () {
            Navigator.pushNamed(
                scaffoldKey.currentContext!, '/sites/${site.name}',
                arguments: site);
          },
        ));
  }

  @override
  Widget build(BuildContext context) {
    // show the content
    var content = _renderSiteButtons().toList();
    return Scaffold(
      appBar: AppBar(title: Text('Book')),
      key: this.scaffoldKey,
      body: Container(
        child: ListView.separated(
          separatorBuilder: (context, index) => Divider(
            height: 10,
          ),
          itemCount: content.length,
          itemBuilder: (context, index) => content[index],
        ),
        margin: EdgeInsets.symmetric(horizontal: 5.0),
      ),
    );
  }
}
