package ck101

const (
	Host = "ck101"
	// url template
	vendorProtocol         = "https"
	vendorHost             = "www.ck101.org"
	bookURLTemplate        = vendorProtocol + "://" + vendorHost + "/%v.html"
	chapterListURLTemplate = vendorProtocol + "://" + vendorHost + "/0/%v/"
	// go query selectors
	bookTitleGoquerySelector       = `placeholder`
	bookWriterGoquerySelector      = `placeholder`
	bookTypeGoquerySelector        = `placeholder`
	bookDateGoquerySelector        = `placeholder`
	bookChapterGoquerySelector     = `placeholder`
	chapterListItemGoquerySelector = `placeholder`
	chapterTitleGoquerySelector    = `placeholder`
	chapterContentGoquerySelector  = `placeholder`
)
