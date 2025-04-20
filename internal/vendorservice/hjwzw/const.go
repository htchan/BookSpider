package hjwzw

const (
	Host = "hjwzw"
	// url template
	vendorProtocol         = "https"
	vendorHost             = "tw.hjwzw.com"
	bookURLTemplate        = vendorProtocol + "://" + vendorHost + "/Book/%v/"
	chapterListURLTemplate = vendorProtocol + "://" + vendorHost + "/Book/Chapter/%v/"
	// go query selectors
	bookTitleGoquerySelector       = `meta[property="og:novel:book_name"]`
	bookWriterGoquerySelector      = `meta[property="og:novel:author"]`
	bookTypeGoquerySelector        = `meta[property="og:novel:category"]`
	bookDateGoquerySelector        = `meta[property="og:novel:update_time"]`
	bookChapterGoquerySelector     = `meta[property="og:novel:latest_chapter_name"]`
	chapterListItemGoquerySelector = `div#tbchapterlist>table>tbody>tr>td>a`
	chapterTitleGoquerySelector    = `td>h1`
	chapterContentGoquerySelector  = `table>tbody>tr>td>div:nth-child(6)`
)
