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
  bool load = false;
  Map<String, dynamic> info, process;
  final GlobalKey scaffoldKey = GlobalKey();

  _StagePageState(this.url) {
    // call backend api
    String infoUrl = '$url/info';
    http.get(infoUrl)
    .then( (response) {
      if (response.statusCode != 404) {
        this.info = Map<String, dynamic>.from(jsonDecode(response.body));
        this.load = true;
        setState((){});
      } else {
        load = false;
      }
    });
    String dataUrl = '$url/process';
    http.get(dataUrl)
    .then( (response) {
      if (response.statusCode != 404) {
        this.process = Map<String, dynamic>.from(jsonDecode(response.body.replaceAll(String.fromCharCode(9), " ")));
        this.process['logs'] = List<String>.from(this.process['logs']).where( (s) => s.startsWith('{'));
        this.load = true;
        setState((){});
      } else {
        load = false;
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
    result.add(Text(this.process['time']));
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
  Widget _renderProcess() {
    if (this.process == null) { return Center(child: Text('Loading logs')); }
    List<String> logs = List<String>.from(this.process['logs']);
    return ListView.builder(
      padding: const EdgeInsets.all(1),
      itemCount: this.process['logs'].length,
      itemBuilder: (BuildContext context, int index) {
        Map<String, dynamic> content = Map<String, dynamic>.from(jsonDecode(logs[index]));
        String subTitle;
        if (content['book'] != null) {
          subTitle = 'title: ' + content['book']['title'] + "\nchapter: " + content['book']['chapter'];
        } else if (content['new'] != null) {
          subTitle = content['old']['title'] + ' -> ' + content['new']['title'];
        } else {
          subTitle = 'id: ' + content['id'].toString();
        }
        return ListTile(
          title: Text(content['site'] + " - " + content['message']),
          subtitle: Text(subTitle),
        );
      }
    );
  }
  
  @override
  Widget build(BuildContext context) {
    List<Widget> stage = this._renderStage();
    List<Widget> subStage = this._renderSubStage();
    // show the content
    return Scaffold(
      appBar: AppBar(title: Text('Stage')),
      key: this.scaffoldKey,
      body: Container(
        child: Column(
          children: [
            Container(
              height: 20.0,
              child: ListView(
                children: stage,
                scrollDirection: Axis.horizontal,
              )
            ),
            Divider(),
            Container(
              height: 20.0,
              child: ListView(
                children: subStage,
                scrollDirection: Axis.horizontal,
              )
            ),
            Divider(),
            Expanded(
              child: this._renderProcess(),
            )
          ],
        ),
          
        margin: EdgeInsets.symmetric(horizontal: 5.0),
      ),
    );
  }
}