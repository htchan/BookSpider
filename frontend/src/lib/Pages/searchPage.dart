import 'package:bookspider/models/all_model.dart';
import 'package:bookspider/repostory/bookSpiderRepostory.dart';
import 'package:flutter/material.dart';
import '../Components/bookList.dart';

class SearchPage extends StatefulWidget {
  final BookSpiderRepostory client;
  final String siteName, title, writer;
  final int page, perPage;

  SearchPage(
      {Key? key,
      required this.client,
      required this.siteName,
      required this.title,
      required this.writer,
      required this.page,
      required this.perPage})
      : super(key: key);

  @override
  _SearchPageState createState() => _SearchPageState(this.client, this.siteName,
      this.title, this.writer, this.page, this.perPage);
}

class _SearchPageState extends State<SearchPage> {
  final BookSpiderRepostory client;
  final String siteName, title, writer;
  int page, perPage;
  List<Book> books = [];
  final GlobalKey scaffoldKey = GlobalKey();
  final ScrollController scrollController;

  _SearchPageState(this.client, this.siteName, this.title, this.writer,
      this.page, this.perPage)
      : this.scrollController = ScrollController() {
    loadPage(this.page);
  }

  void loadPage(int page) {
    this
        .client
        .searchBooks(
            site: this.siteName,
            title: this.title,
            writer: this.writer,
            page: page,
            perPage: this.perPage)
        .then((books) {
      setState(() {
        this.books = books;
      });
    });
  }

  Widget loadAboveButton(ScrollController controller) {
    if (this.page == 1) return SizedBox.shrink();
    return TextButton(
      child: Text('Load Above'),
      onPressed: () {
        setState(() {
          this.loadPage(--this.page);
        });
      },
    );
  }

  Widget loadMoreButton(ScrollController controller) {
    if (this.books.length < 20) return SizedBox.shrink();
    return TextButton(
      child: Text('Load More'),
      onPressed: () {
        setState(() {
          this.loadPage(++this.page);
        });
        controller.animateTo(0,
            duration: Duration(milliseconds: 500), curve: Curves.fastOutSlowIn);
      },
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
          loadAboveButton,
          loadMoreButton,
        ),
        margin: EdgeInsets.symmetric(horizontal: 5.0),
      ),
    );
  }
}
