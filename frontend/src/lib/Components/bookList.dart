import 'package:bookspider/models/all_model.dart';
import 'package:flutter/material.dart';

class BookList extends StatelessWidget {
  GlobalKey scaffoldKey;
  String siteName;
  List<Book> books;
  final Function firstButton, lastButton;
  final ScrollController scrollController = ScrollController();

  BookList(this.scaffoldKey, this.siteName, this.books, this.firstButton,
      this.lastButton);

  @override
  Widget build(BuildContext context) {
    if (books.length == 0) {
      return Center(child: Text('no books found'));
    }

    List<Widget> list = [];
    list.add(firstButton(scrollController));
    list.addAll(books.map((book) => ListTile(
        title: Text('${book.title} - ${book.writer}'),
        subtitle: Text('${book.updateDate} - ${book.updateChapter}'),
        onTap: () {
          Navigator.pushNamed(
            this.scaffoldKey.currentContext!,
            '/sites/$siteName/books/${book.id}-${book.hash}',
            arguments: book,
          );
        })));

    list.add(lastButton(scrollController));

    return ListView.separated(
      controller: scrollController,
      separatorBuilder: (context, index) => Divider(
        height: 10,
      ),
      itemCount: list.length,
      itemBuilder: (context, index) => list[index],
    );
  }
}
