import 'dart:convert';
import 'dart:async';

import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

void main() {
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      home: MyHomePage(),
    );
  }
}

class MyHomePage extends StatefulWidget {
  @override
  _MyHomePageState createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  int numberOfStocks = 5; // Set an initial number of stocks
  TextEditingController _stocksController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Stock Prices'),
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            TextField(
              controller: _stocksController,
              keyboardType: TextInputType.number,
              decoration: InputDecoration(
                labelText: 'Number of Stocks',
              ),
            ),
            ElevatedButton(
              onPressed: () {
                _fetchStockPrices();
              },
              child: Text('Fetch Stock Prices'),
            ),
            FutureBuilder(
              future: _fetchStockPrices(),
              builder: (context,snapshot) {
                if (snapshot.connectionState == ConnectionState.waiting) {
                  return CircularProgressIndicator();
                } else if (snapshot.hasError) {
                  return Text('Error: ${snapshot.error}');
                } else {
                      var prices=[...?snapshot.data];
                  return Column(
                    children: [
                      Text('Stock Prices:'),
                      for (var price in prices) Text('${price['Name']}: \$${price['Price']}'),
                    ],
                  );
                }
              },
            ),
          ],
        ),
      ),
    );
  }

  Future<List<dynamic>> _fetchStockPrices() async {
    int n = int.tryParse(_stocksController.text) ?? 5; // Use the input or default to 5
    final response = await http.get(Uri.parse('http://localhost:3000/fetch-stocks/$n'));

    if (response.statusCode == 200) {
      // If the server returns a 200 OK response, parse the JSON
      return json.decode(response.body);
    } else {
      // If the server did not return a 200 OK response, throw an exception.
      throw Exception('Failed to fetch stock prices');
    }
  }
}
