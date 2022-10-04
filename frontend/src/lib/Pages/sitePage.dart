import 'package:bookspider/models/all_model.dart';
import 'package:bookspider/repostory/bookSpiderRepostory.dart';
import 'package:flutter/material.dart';

import '../Components/bookSearchBar.dart';
import '../Components/siteChartPanel.dart';
import '../Components/siteInfoPanel.dart';

class SitePage extends StatefulWidget {
  final BookSpiderRepostory client;
  final String siteName;
  final Site? site;

  SitePage({Key? key, required this.client, required this.siteName, this.site})
      : super(key: key);

  @override
  _SitePageState createState() =>
      _SitePageState(this.client, this.siteName, this.site);
}

class _SitePageState extends State<SitePage>
    with SingleTickerProviderStateMixin {
  final String siteName;
  final BookSpiderRepostory client;
  final GlobalKey scaffoldKey = GlobalKey();
  Site? site;
  bool isError = false;

  _SitePageState(this.client, this.siteName, this.site) {
    if (this.site == null) {
      this.client.getSite(this.siteName).then((site) {
        setState(() {
          this.site = site;
        });
      }).catchError((e) {
        setState(() {
          this.isError = true;
        });
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    Widget _chartPanel = Center(child: Text("Loading Chart"));
    Widget _dataPanel = Center(child: Text("Loading Data"));
    if (this.site != null) {
      _chartPanel = SiteChartPanel(scaffoldKey, this.site!);
      _dataPanel = SiteInfoPanel(scaffoldKey, this.site!);
    } else if (isError) {
      _chartPanel = Center(child: Text("Fail to load Chart"));
      _dataPanel = Center(child: Text("Fail to load Data"));
    }
    // show the content
    final PageController pageController = PageController(initialPage: 0);
    return Scaffold(
        appBar: AppBar(title: Text(siteName)),
        key: scaffoldKey,
        body: Container(
          child: Column(
            children: [
              Container(
                  height: MediaQuery.of(context).size.height * 0.5,
                  child: PageView(
                    children: [_chartPanel, _dataPanel],
                    controller: pageController,
                  )),
              BookSearchBar(scaffoldKey: scaffoldKey, siteName: siteName),
              // _renderRandomButton(),
            ],
          ),
          margin: EdgeInsets.symmetric(horizontal: 5.0),
        ));
  }
}
