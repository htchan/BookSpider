import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import '../Components/logList.dart';

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
  Widget _stage, _subStage, _logs;
  final GlobalKey scaffoldKey = GlobalKey();

  _StagePageState(this.url) {
    // call backend api
    String dataUrl = '$url/process';
    _stage = Center(child: Text('Loading Stage...'));
    _subStage = Center(child: Text('Loading Sub-Stage...'));
    _logs = Center(child: Text('Loading Logs...'));
    http.get(Uri.parse(dataUrl))
    .then( (response) {
      if (response.statusCode >= 200 && response.statusCode < 300) {
        Map<String, dynamic> info = Map<String, dynamic>.from(
          jsonDecode(response.body.replaceAll(String.fromCharCode(9), " "))
        );
        info['logs'] = List<String>.from(info['logs'])
          .where( (s) => s.length > 20 );
        setState((){
          _stage = _renderStage(info);
          _subStage = _renderSubStage(info);
          _logs = LogList(logs: List<String>.from(info['logs']));
        });
      } else {
        setState((){
          _stage = Center(child: Text('Error${response.statusCode}'));
          _subStage = Center(child: Text('Error${response.statusCode}'));
          _logs = Center(
            child: Column(
              children: [
                Text(response.statusCode.toString()),
                Text(response.body)
              ],
            )
          );
        });
      }
    });
  }
  void _addStage(List<Widget> list, String line) {
    if (line.contains('start')) {
      list.add(Text(
        line.replaceAll('sub_', '').replaceAll('stage: ', '').replaceAll(' stage', ''),
        style: TextStyle(backgroundColor: Colors.blue, color: Colors.white)
      ));
      list.add(VerticalDivider());
    } else {
      String target = line.replaceAll('sub_', '').replaceAll('stage: ', '').replaceAll(' stage', '');
      list.asMap().keys.toList().map( (i) {
        if (list[i] is Text && (list[i] as Text).data.split(' ')[0] == target.split(' ')[0]) {
          list[i] = Text(
            target,
            style: TextStyle(backgroundColor: Colors.green, color: Colors.white)
          );
        }
      }).toList();
      // list[list.length - 2] = Text(
      //   line.replaceAll('sub_', '').replaceAll('stage: ', '').replaceAll(' stage', ''),
      //   style: TextStyle(backgroundColor: Colors.green, color: Colors.white)
      // );
    }
  }
  Widget _renderStage(Map<String, dynamic> info) {
    List<Widget> stages = [];
    for (String line in info['stage']) {
      if (line.indexOf('stage') == 0) {
        _addStage(stages, line);
      }
    }
    stages.add(Text(DateTime.fromMillisecondsSinceEpoch(info['time'] * 1000).toString()));
    return ListView(
      children: stages,
      scrollDirection: Axis.horizontal,
    );
  }
  
  Widget _renderSubStage(Map<String, dynamic> info) {
    List<Widget> subStages = [];
    for (String line in info['stage']) {
      if (line.indexOf('stage') == 0) {
        if (line.contains('start')) { subStages.clear(); }
      }
      if (line.indexOf('sub_stage') == 0) {
        _addStage(subStages, line);
      }
    }
    return ListView(
      children: subStages,
      scrollDirection: Axis.horizontal,
    );
  }
  
  @override
  Widget build(BuildContext context) {
    // show the content
    return Scaffold(
      appBar: AppBar(title: Text('Stage')),
      key: scaffoldKey,
      body: Container(
        child: Column(
          children: [
            Container(
              height: 20.0,
              child: _stage,
            ),
            Divider(),
            Container(
              height: 20.0,
              child: _subStage,
            ),
            Divider(),
            Expanded(
              child: _logs,
            )
          ],
        ),
          
        margin: EdgeInsets.symmetric(horizontal: 5.0),
      ),
    );
  }
}