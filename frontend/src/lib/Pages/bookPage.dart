import 'package:bookspider/models/all_model.dart';
import 'package:bookspider/repostory/bookSpiderRepostory.dart';
import 'package:flutter/material.dart';
import 'package:url_launcher/url_launcher.dart';

class BookPage extends StatefulWidget {
  final BookSpiderRepostory client;
  final String siteName, id, hash;
  final Book? book;

  BookPage(
      {Key? key,
      required this.client,
      required this.siteName,
      required this.id,
      required this.hash,
      this.book})
      : super(key: key);

  @override
  _BookPageState createState() =>
      _BookPageState(this.client, this.siteName, this.id, this.hash, this.book);
}

class _BookPageState extends State<BookPage> {
  final BookSpiderRepostory client;
  final String siteName, id, hash;
  Book? book;
  final GlobalKey scaffoldKey = GlobalKey();

  _BookPageState(this.client, this.siteName, this.id, this.hash, this.book) {
    if (this.book == null) {
      this
          .client
          .getBook(site: siteName, id: int.parse(id), hash: hash)
          .then((book) {
        setState(() {
          this.book = book;
        });
      });
    }
  }

  Widget _renderBookContent() {
    if (book == null) {
      return Center(
        child: Text("loading"),
      );
    }
    List<Widget> rows = [
      SelectableText('ID: ${book!.id} - ${book!.hash}'),
      Row(
        // TODO: extract this to an external widget
        children: [
          SelectableText('Title: ${book!.title}'),
          ElevatedButton(
            child: Text("Search Google"),
            onPressed: () {
              // TODO: open new page to search google
              Uri searchUrl =
                  Uri.parse("https://www.google.com/search?q=${book!.title}");
              canLaunchUrl(searchUrl).then((result) {
                if (result) launchUrl(searchUrl);
              });
            },
          )
        ],
      ),
      Row(
        // TODO: extract this to external widget
        children: [
          SelectableText('Writer: ${book!.writer}'),
          ElevatedButton(
            child: Text("Search"),
            onPressed: () {
              // TODO: search the writer name internally
              String writer = book!.writer;
              Navigator.pushNamed(this.scaffoldKey.currentContext!,
                  '/sites/$siteName/search/?writer=$writer');
            },
          )
        ],
      ),
      SelectableText('Type: ${book!.type}'),
      SelectableText('Last Update: ${book!.updateDate}'),
      SelectableText('Last Chapter: ${book!.updateChapter}'),
      SelectableText('Status: ${book!.status}')
    ];
    if (book!.isDownload) {
      rows.add(TextButton(
        child: Text('Download'),
        onPressed: () => this
            .client
            .downloadBook(site: siteName, id: int.parse(id), hash: hash),
      ));
    }
    return ListView.separated(
      separatorBuilder: (context, index) => Divider(
        height: 10,
      ),
      itemCount: rows.length,
      itemBuilder: (context, index) => rows[index],
    );
  }

  @override
  Widget build(BuildContext context) {
    // show the content
    return Scaffold(
        appBar: AppBar(title: Text(this.siteName)),
        key: this.scaffoldKey,
        body: Container(
          child: _renderBookContent(),
          margin: EdgeInsets.symmetric(horizontal: 5.0),
        ));
  }
}
