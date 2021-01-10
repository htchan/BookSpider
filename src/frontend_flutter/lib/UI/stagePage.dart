import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

class StagePage extends StatefulWidget{
  final String url;

  StagePage({Key key, this.url}) : super(key: key);

  @override
  _StagePageState createState() => _StagePageState(this.url);
}

class _StagePageState extends State<StagePage> {
  final String url;
  bool error = true;
  Map<String, dynamic> info, process;
  final GlobalKey scaffoldKey = GlobalKey();

  _StagePageState(this.url) {
    // call backend api
    String infoUrl = '$url/info';
    http.get(infoUrl)
    .then( (response) {
      if (response.statusCode != 404) {
        this.info = Map<String, dynamic>.from(jsonDecode(response.body));
        this.error = false;
        setState((){});
      }
    });
    String dataUrl = '$url/process';
    http.get(dataUrl)
    .then( (response) {
      if (response.statusCode != 404) {
        this.process = Map<String, dynamic>.from(jsonDecode(response.body.replaceAll(String.fromCharCode(9), " ")));
        this.error = false;
        setState((){});
      }
    });
  }
  void _addStage(List list, String line) {
    if (line.contains('start')) {
      list.add(Text(
        line.replaceAll('sub_', '').replaceAll('stage: ', '').replaceAll(' stage', ''),
        style: TextStyle(backgroundColor: Colors.blue, color: Colors.white)
      ));
      list.add(VerticalDivider());
    } else {
      list[list.length - 2] = Text(
        line.replaceAll('sub_', '').replaceAll('stage: ', '').replaceAll(' stage', ''),
        style: TextStyle(backgroundColor: Colors.green, color: Colors.white)
      );
    }
  }
  List<Widget> _renderStage() {
    if (this.process == null) { return []; }
    List<Widget> result = [];
    for (String line in this.info['stage'].split('\n')) {
      if (line.indexOf('stage') == 0) {
        this._addStage(result, line);
      }
    }
    return result;
  }
  List<Widget> _renderSubStage() {
    if (this.process == null) { return []; }
    List<Widget> result = [];
    for (String line in this.info['stage'].split('\n')) {
      if (line.indexOf('stage') == 0) {
        if (line.contains('start')) { result.clear(); }
      }
      if (line.indexOf('sub_stage') == 0) {
        this._addStage(result, line);
      }
    }
    return result;
  }
  List<Widget> _renderProcess() {
    if (this.process == null) { return []; }
    List<Widget> result = [Text(this.process['time'])];
    for (String line in this.process['logs']) {
      if (line.length > 0) { result.add(Text(line)); }
    }
    return result;
  }
  
  @override
  Widget build(BuildContext context) {
    List<Widget> stage = this._renderStage();
    List<Widget> subStage = this._renderSubStage();
    List<Widget> process = this._renderProcess();
    // show the content
    return Scaffold(
      appBar: AppBar(title: Text('Stage')),
      key: this.scaffoldKey,
      body: Container(
        child: Column(
          children: [
            Row(children: stage,),
            Divider(),
            Row(children: subStage,),
            Divider(),
            Expanded(
              child: ListView(
                children: process
              )
            )
          ],
        ),
          
        margin: EdgeInsets.symmetric(horizontal: 5.0),
      ),
    );
  }
}