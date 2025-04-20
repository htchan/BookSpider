package xbiquge

const (
	Host = "xbiquge"
	// url template
	vendorProtocol         = "https"
	vendorHost             = "www.xbiquge.bz"
	bookURLTemplate        = vendorProtocol + "://" + vendorHost + "/book/%v/"
	chapterListURLTemplate = vendorProtocol + "://" + vendorHost + "/book/%v/"
	chapterURLTemplate     = vendorProtocol + "://" + vendorHost + "/book/%v/%v"
	// go query selectors
	bookTitleGoquerySelector       = `meta[property="og:novel:book_name"]`
	bookWriterGoquerySelector      = `meta[property="og:novel:author"]`
	bookTypeGoquerySelector        = `meta[property="og:novel:category"]`
	bookDateGoquerySelector        = `meta[property="og:novel:update_time"]`
	bookChapterGoquerySelector     = `meta[property="og:novel:latest_chapter_name"]`
	chapterListItemGoquerySelector = `dd>a`
	chapterTitleGoquerySelector    = `div.bookname>h1`
	chapterContentGoquerySelector  = `div#content`
)
