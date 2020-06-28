from .. import models
import re

class Book(models.Book):
    def _get_title(self, html):
        html = html[re.search("<div class=\"tit1\">", html).end():]
        html = re.search("<h1>(.*?)</h1>", html).group(1)
        return html
    def _get_writer(self, html):
        html = re.search("<span>小说作者：(.*?)</span>", html).group(1)
        return html
    def _get_type(self, html):
        html = html[re.search("<div class=\"crumbs\">", html).end():]
        html = re.search("<a href=\"/ls/.*?html\">(.*?)</a>", html).group(1)
        return html
    def _get_last_update(self, html):
        html = re.search("<span>更新日期：(.*?)</span>", html).group(1)
        return html
    def _get_last_chapter(self, html):
        html = html[re.search("<div class=\"new_con\">", html).end():]
        html = html[re.search("最新章节：", html).end():]
        html = re.search("<a.*?>(.*?)</a>", html).group(1)
        return html
    def _get_chapters_url(self, html):
        try:
            html = html[re.search("<div class=\"book_con_list\">", html).end():]
            html = html[re.search("<div class=\"book_con_list\">", html).end():]
            html = re.sub("(&nbsp;|<ul>|</ul>|<li>|</li>)", "", html[:re.search("</div>", html).start()])
            html = re.sub("</(.*?)>", "</\\1>\n", html).split("\n")
            for i in range(len(c)):
                if (c[i].find("<a") == 0):
                    if (re.search("<a href=\".*?\">.*?</a>", html[i]) == None):
                        continue
                    html[i] = html[i].replace('”', '"')
                    html[i] = self.download_url + re.search("<a href=\"(.*?)\"", html[i]).group(1)
                else:
                    html[i] = ""
            return html
        except:
            return []
    def _get_chapters_title(self, html):
        try:
            html = html[re.search("<div class=\"book_con_list\">", html).end():]
            html = html[re.search("<div class=\"book_con_list\">", html).end():]
            html = re.sub("(&nbsp;|<ul>|</ul>|<li>|</li>)", "", html[:re.search("</div>", html).start()])
            html = re.sub("(</?h\d>|</?a.*?>)", "", re.sub("</(.*?)>", "</\\1>\n", html)).split("\n")
            return html
        except:
            return []
    def _get_content(self, html):
        html = html[re.search("id=\"content\">", html).end():]
        html = html[:re.search("<div", html).start()]
        html = re.sub("&nbsp;"," ", html)
        html = re.sub("<br />    ","\n", html)
        html = html.replace("\r\n\r\n", "\r\n")
        return html

class Site(models.Site):
    def get_book(self, num):
        base_url = self.meta_base_url.format(num)
        download_url = self.meta_download_url.format(num)
        return Book(self.db, base_url, download_url, self.chapter_url, self.site, num, self.decode, self.max_thread, self.timeout)

desktop_setting = models.Setting(
    meta_base_url="http://www.xqiushu.com/txt{}/",
    meta_download_url="http://www.xqiushu.com/t/{}/",
    chapter_url="http://www.xqiushu.com/t/\d"
)

mobile_setting = models.Setting(
    meta_base_url="http://www.xqiushu.com/txt{}/",
    meta_download_url="http://www.xqiushu.com/t/{}/",
    chapter_url="http://www.xqiushu.com/t/\d"
)
