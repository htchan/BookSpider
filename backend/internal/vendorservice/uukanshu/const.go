package uukanshu

const (
	Host = "uukanshu"
	// url template
	vendorProtocol         = "https"
	vendorHost             = "www.uukanshu.com"
	bookURLTemplate        = vendorProtocol + "://" + vendorHost + "/b/%v/"
	chapterListURLTemplate = vendorProtocol + "://" + vendorHost + "/b/%v/"
	// go query selectors
	bookTitleGoquerySelector       = `div.xiaoshuo_content>dl.jieshao>dd.jieshao_content>h1>a`
	bookWriterGoquerySelector      = `div.xiaoshuo_content>dl.jieshao>dd.jieshao_content>h2>a`
	bookTypeGoquerySelector        = `div.weizhi>div.path>a:nth-child(2)`
	bookDateGoquerySelector        = `div.xiaoshuo_content>dl.jieshao>dd.jieshao_content>div.shijian`
	bookChapterGoquerySelector     = `div.zhangjie>ul#chapterList>li:first-child>a`
	chapterListItemGoquerySelector = `div.zhangjie>ul#chapterList>li>a`
	chapterTitleGoquerySelector    = `div.zhengwen_box>div.box_left>div.w_main>div.h1title>h1#timu`
	chapterContentGoquerySelector  = `div.zhengwen_box>div.box_left>div.w_main>div.contentbox>div#contentbox`
)
