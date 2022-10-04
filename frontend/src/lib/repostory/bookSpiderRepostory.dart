import 'dart:convert';

import 'package:bookspider/models/all_model.dart';
import 'package:http/http.dart' as http;
import 'package:url_launcher/url_launcher.dart';

const String api_prefix = String.fromEnvironment(
    "NOVEL_SPIDER_API_ROUTE_PREFIX",
    defaultValue: "/api/novel");

class BookSpiderRepostory {
  String host;
  final String protocol = "http";

  BookSpiderRepostory(this.host);

  String get url {
    return "${protocol}://${host}${api_prefix}";
  }

  Future<void> downloadBook(
      {required String site, required int id, String? hash}) {
    String url = "${this.url}/sites/${site}/books/${id}";
    if (hash != null) {
      url += "-${hash}";
    }
    url += "/download";
    return launchUrl(Uri.parse(url), mode: LaunchMode.externalApplication);
  }

  Future<Book> getBook({required String site, required int id, String? hash}) {
    String url = "${this.url}/sites/${site}/books/${id}";
    if (hash != null) {
      url += "-${hash}";
    }
    return http.get(Uri.parse(url)).then((response) {
      Map<String, dynamic> responseMap = Map.from(jsonDecode(response.body));
      return Book.from(responseMap);
    });
  }

  Future<List<Book>> searchBooks(
      {required String site,
      required title,
      required writer,
      required int page,
      required perPage}) {
    return http
        .get(Uri.parse(
            "${this.url}/sites/${site}/books/search?title=${title}&writer=${writer}&page=${page}&per_page=${perPage}"))
        .then((response) {
      Map<String, dynamic> responseMap = Map.from(jsonDecode(response.body));
      List booksResponse = responseMap['books'];
      return List<Book>.from(booksResponse
          .map((resp) => Book.from(Map<String, dynamic>.from(resp))));
    });
  }

  Future<List<Book>> randomBook({required String site, required int perPage}) {
    return http
        .get(Uri.parse(
            "${this.url}/sites/${site}/books/random?per_page=${perPage}"))
        .then((response) {
      Map<String, dynamic> responseMap = Map.from(jsonDecode(response.body));
      List booksResponse = responseMap['books'];
      return List<Book>.from(booksResponse
          .map((resp) => Book.from(Map<String, dynamic>.from(resp))));
    });
  }

  Future<Site> getSite(String name) {
    return http.get(Uri.parse("${this.url}/sites/${name}")).then((response) {
      Map<String, dynamic> responseMap = Map.from(jsonDecode(response.body));
      return Site.from(name, responseMap);
    });
  }

  Future<List<Site>> getSites() {
    return http.get(Uri.parse("${this.url}/info")).then((response) {
      Map<String, dynamic> responseMap = Map.from(jsonDecode(response.body));
      return List<Site>.from(
          responseMap.keys.map((key) => Site.from(key, responseMap[key])));
    });
  }
}
