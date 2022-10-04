import 'package:flutter/material.dart';
import 'package:fluttericon/rpg_awesome_icons.dart' as Icons2;

class BookSearchBar extends StatefulWidget {
  final String siteName;
  final GlobalKey scaffoldKey;

  BookSearchBar({Key? key, required this.scaffoldKey, required this.siteName})
      : super(key: key);

  @override
  _BookSearchBarState createState() => _BookSearchBarState(this.scaffoldKey);
}

class _BookSearchBarState extends State<BookSearchBar> {
  final GlobalKey scaffoldKey;

  _BookSearchBarState(this.scaffoldKey);

  Widget searchField(String labelText, TextEditingController textController) {
    return TextField(
        decoration: InputDecoration(labelText: labelText),
        controller: textController);
  }

  Widget searchButton(TextEditingController titleController,
      TextEditingController writerController) {
    return TextButton.icon(
      icon: Icon(Icons.search),
      label: Text('Search'),
      onPressed: () {
        String title = titleController.text;
        String writer = writerController.text;
        Navigator.pushNamed(this.scaffoldKey.currentContext!,
            '/sites/${this.widget.siteName}/search?title=$title&writer=$writer');
      },
    );
  }

  Widget randomButton() {
    return TextButton.icon(
      icon: Icon(Icons2.RpgAwesome.perspective_dice_random),
      label: Text('Random'),
      onPressed: () {
        Navigator.pushNamed(this.scaffoldKey.currentContext!,
            '/sites/${this.widget.siteName}/random');
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    TextEditingController titleController, writerController;
    titleController = TextEditingController();
    writerController = TextEditingController();
    return Center(
      child: Card(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: <Widget>[
            searchField('Book Title', titleController),
            searchField('Book Writer', writerController),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Expanded(child: Container()),
                searchButton(titleController, writerController),
                Expanded(child: Container()),
                randomButton(),
                Expanded(child: Container()),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
