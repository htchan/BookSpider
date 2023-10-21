package xqishu

const (
	Host = "xqishu"
	// url template
	vendorProtocol         = "https"
	vendorHost             = "www.aidusk.com"
	bookURLTemplate        = vendorProtocol + "://" + vendorHost + "/txt%v/"
	chapterListURLTemplate = vendorProtocol + "://" + vendorHost + "/t/%v/"
	// go query selectors
	bookTitleGoquerySelector       = `div.tit1>h1`
	bookWriterGoquerySelector      = `div.date>span:nth-child(1)`
	bookTypeGoquerySelector        = `div.crumbs>a:nth-child(2)`
	bookDateGoquerySelector        = `div.date>span:nth-child(3)`
	bookChapterGoquerySelector     = `a.zx_zhang`
	chapterListItemGoquerySelector = `div.book_con_list>ul>li>a`
	chapterTitleGoquerySelector    = `div.date>h1`
	chapterContentGoquerySelector  = `div.book_content`
)
