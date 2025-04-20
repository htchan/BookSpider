package bestory

const (
	Host = "bestory"
	// url template
	vendorProtocol         = "https"
	vendorHost             = "www.8book.com"
	bookURLTemplate        = vendorProtocol + "://" + vendorHost + "/novelbooks/%v/"
	chapterListURLTemplate = vendorProtocol + "://" + vendorHost + "/novelbooks/%v/"
	// go query selectors
	bookTitleGoquerySelector       = `meta[property="og:novel:book_name"]`
	bookWriterGoquerySelector      = `meta[property="og:novel:author"]`
	bookTypeGoquerySelector        = `meta[property="og:novel:category"]`
	bookDateGoquerySelector        = `div>div.txt_info:nth-child(4)`
	bookChapterGoquerySelector     = `div.yulan:last-child>a`
	chapterListItemGoquerySelector = `div.yulan>li>a`
	chapterTitleGoquerySelector    = `div.date>h1`
	chapterContentGoquerySelector  = `div.book_content`
)
