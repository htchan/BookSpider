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
          .where( (s) => s.startsWith('{') )
          .map( (s) => s.replaceAll('\\\"', '\"'));
        setState((){
          _stage = _renderStage(info);
          _subStage = _renderSubStage(info);
          _logs = _renderProcesses(info);
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
    stages.add(Text(info['time']));
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

  Widget _renderProcesses(Map<String, dynamic> info) {
    List<String> logs = List<String>.from(info['logs']);
    return ListView.builder(
      padding: const EdgeInsets.all(1),
      itemCount: info['logs'].length,
      itemBuilder: (BuildContext context, int index) {
        Map<String, dynamic> content = Map<String, dynamic>.from(jsonDecode(logs[index]));
        String subTitle;
        if (content['book'] != null) {
          subTitle = 'title: ${content['book']['title']}\nchapter: ${content['book']['chapter']}';
        } else if (content['new'] != null) {
          subTitle = '${content['old']['title']} -> ${content['new']['title']}';
        } else {
          subTitle = 'id: ${content['id'].toString()}';
        }
        return ListTile(
          title: Text('${content['site']} - ${content['message']}'),
          subtitle: Text(subTitle),
        );
      }
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