import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

class BookList extends StatelessWidget {
  GlobalKey scaffoldKey;
  String siteName;
  List<Map<String, dynamic>> books;
  final Function firstButton, lastButton;
  final ScrollController scrollController = ScrollController();

  BookList(this.scaffoldKey, this.siteName, this.books, this.firstButton, this.lastButton);

  @override
  Widget build(BuildContext context) {
    if (books.length == 0) { return Center(child: Text('no books found')); }

    List<Widget> list = [];
    if (firstButton != null) { list.add(firstButton(scrollController)); }
    list.addAll(books.map( (book) => ListTile(
      title: Text('${book['title']} - ${book['writer']}'),
      subtitle: Text('${book['updateDate']} - ${book['updateChapter']}'),
      onTap: () {
        Navigator.pushNamed(
          this.scaffoldKey.currentContext,
          '/books/$siteName/${book['id']}'
        );
    })));

    if (lastButton != null && books.length == 20) { list.add(lastButton(scrollController)); }
    
    return ListView.separated(
      controller: scrollController,
      separatorBuilder: (context, index) => Divider(height: 10,),
      itemCount: list.length,
      itemBuilder: (context, index) => list[index],
    );
  }
}