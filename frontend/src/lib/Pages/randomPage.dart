import 'package:bookspider/models/all_model.dart';
import 'package:bookspider/repostory/bookSpiderRepostory.dart';
import 'package:flutter/material.dart';
import '../Components/bookList.dart';

class RandomPage extends StatefulWidget {
  final BookSpiderRepostory client;
  final String siteName;

  RandomPage({Key? key, required this.client, required this.siteName})
      : super(key: key);

  @override
  _RandomPageState createState() =>
      _RandomPageState(this.client, this.siteName);
}

class _RandomPageState extends State<RandomPage> {
  final BookSpiderRepostory client;
  final String siteName;
  int perPage = 20;
  List<Book> books = [];
  final GlobalKey scaffoldKey = GlobalKey();
  final ScrollController scrollController;

  _RandomPageState(this.client, this.siteName)
      : this.scrollController = ScrollController() {
    this._loadPage();
  }

  void _loadPage() {
    this.client.randomBook(site: siteName, perPage: perPage).then((books) {
      setState(() {
        this.books = books;
      });
    });
  }

  Widget randomButton(ScrollController controller) {
    return ListTile(
      title:
          Center(child: Text('Reload', style: TextStyle(color: Colors.blue))),
      onTap: () {
        setState(() {
          this._loadPage();
        });
        controller.animateTo(0,
            duration: Duration(milliseconds: 500), curve: Curves.fastOutSlowIn);
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
        child: BookList(
          this.scaffoldKey,
          this.siteName,
          this.books,
          (_) => SizedBox.shrink(),
          randomButton,
        ),
        margin: EdgeInsets.symmetric(horizontal: 5.0),
      ),
    );
  }
}
