from .. import models
import re

class Book(models.Book):
    def _get_title(self, html):
        html = html[re.search("titlename", html).end():]
        html = html[re.search("<h1>", html).end():re.search("</h1>", html).start()]
        html = re.sub("全文阅读", "", html)
        return html
    def _get_writer(self, html):
        html = html[re.search("作者：", html).end():]
        html = html[re.search(">", html).end():re.search("</a>", html).start()]
        html = re.sub("'", "", html)
        return html
    def _get_type(self, html):
        html = html[re.search("分类：", html).end():]
        html = html[re.search(">", html).end():re.search("</a>", html).start()]
        return html
    def _get_last_update(self, html):
        m = re.search("\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}", html)
        html = html[m.start():m.end()]
        return html
    def _get_last_chapter(self, html):
        while(re.search("<li>", html)):
            m = re.search("<li>", html)
            html = html[m.end():]
        html = html[:re.search("</li>", html).start()]
        html = re.sub("(<a.*?>|</a>)", "", html)
        return html
    def _get_chapters_url(self, html):
        try:
            html = html[re.search("yulan..", html).end():]
            html = html[:re.search("</div>", html).start()]
            html = re.findall("(http.*?html)", html)
            return html
        except:
            return []
    def _get_chapters_title(self, html):
        try:
            html = html[re.search("yulan..", html).end():]
            html = html[:re.search("</div>", html).start()]
            html = re.sub("<strong>.*?</strong>", "", html)
            html = re.sub(".*>(.+?)<(a|\/a).*", "\\1", html).split("\n")
            html = html[1:len(html)-1]
            return html
        except:
            return []
    def _get_content(self, html):
        html = html[re.search("id=\"content\">", html).end():]
        html = html[:re.search("<div", html).start()]
        html = re.sub("&nbsp;", " ", html)
        html = re.sub("<br />    ", "\n", html)
        html = html.replace("\r\n\r\n", "\r\n")
        return html

class Site(models.Site):
    def get_book(self, num):
        base_url = self.meta_base_url.format(num)
        download_url = self.meta_download_url.format(num)
        return Book(self.db, base_url, download_url, self.chapter_url, self.site, num, self.decode, self.max_thread, self.timeout)

desktop_setting = models.Setting(
    meta_base_url="https://www.balingtxt.com/txtml_{}.html",
    meta_download_url="http://www.balingtxt.com/txtml_{}.html",
    chapter_url="http://www.xqiushu.com"
)

mobile_setting = models.Setting(
    meta_base_url="https://www.balingtxt.com/txtml_{}.html",
    meta_download_url="http://www.balingtxt.com/txtml_{}.html",
    chapter_url="http://www.xqiushu.com"
)
