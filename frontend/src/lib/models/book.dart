class Book {
  final String site,
      hash,
      title,
      writer,
      type,
      updateDate,
      updateChapter,
      status,
      error;
  final int id;
  final bool isDownload;

  Book.from(Map<String, dynamic> map)
      : this.site = map['site'] ?? "",
        this.id = map['id'] ?? 0,
        this.hash = map['hash_code'] ?? "",
        this.title = map['title'] ?? "",
        this.writer = map['writer'] ?? "",
        this.type = map['type'] ?? "",
        this.updateDate = map['update_date'] ?? "",
        this.updateChapter = map['update_chapter'] ?? "",
        this.status = map['status'] ?? "",
        this.isDownload = map['is_downloaded'] ?? false,
        this.error = map['error'] ?? "unknown error";
}
