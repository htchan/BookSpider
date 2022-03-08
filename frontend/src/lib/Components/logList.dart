import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;


class LogList extends StatelessWidget{
  final List<String> logs;

  LogList({Key key, this.logs}) : super(key: key);

  Widget itemBuilder(BuildContext context, int index) {
    DateTime loggingTime = DateTime.parse(this.logs[index].substring(0, 19).replaceAll('/', '-'));
    String text = this.logs[index].substring(20);
    return ListTile(
      title: Text(text),
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