import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:js' as js;

class ErrorPage extends StatelessWidget{

  ErrorPage({Key key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text('Error')),
      body: Container(
        child: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(
                Icons.error_outline,
                size: 250,
              ),
              Text(
                'Page Not Exist',
                style: TextStyle(
                  fontSize: 50,
                ),
              ),
            ],
          ),
        ),
        margin: EdgeInsets.symmetric(horizontal: 5.0),
      ),
    );
  }
}
