package baling

const (
	Host = "80txt"
	// url template
	vendorProtocol         = "https"
	vendorHost             = "www.80zw.info/"
	bookURLTemplate        = vendorProtocol + "://" + vendorHost + "/txtml_%v.html"
	chapterListURLTemplate = vendorProtocol + "://" + vendorHost + "/txtml_%v.html"
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
