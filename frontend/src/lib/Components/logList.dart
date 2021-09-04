import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;


class LogList extends StatelessWidget{
  final List<String> logs;

  LogList({Key key, this.logs}) : super(key: key);

  Widget itemBuilder(BuildContext context, int index) {
    DateTime loggingTime = DateTime.parse(this.logs[index].substring(0, 19).replaceAll('/', '-'));
    Map<String, dynamic> content = Map<String, dynamic>.from(jsonDecode(logs[index].substring(20)));
    String subTitle;
    if (content['book'] != null) {
      subTitle = 'title: ${content['book']['title']}\nchapter: ${content['book']['chapter']}';
    } else if (content['new'] != null) {
      subTitle = '${content['old']['title']} -> ${content['new']['title']}';
    } else {
      subTitle = 'id: ${content['id'].toString()}';
    }
    return ListTile(
      title: Text('${content['site']}-${content['id']} : ${content['message']}'),
      subtitle: Text(loggingTime.toString()),
    );
  }

  @override
  Widget build(BuildContext context) {
    return ListView.builder(
      padding: const EdgeInsets.all(1),
      itemCount: this.logs.length,
      itemBuilder: this.itemBuilder
    );
  }
}