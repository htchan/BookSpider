from .. import models
import re

class Book(models.Book):
    def _get_title(self, html):
        html = html[re.search("<title>",html).end():]
        html = html[:re.search("</title>",html).start()]
        html = re.sub("('|\x00)","",html).replace(" 在線觀看", "")
        return html
    def _get_writer(self, html):
        html = html[re.search("作者:",html).end():]
        html = html[re.search("\">",html).end():re.search("</font>",html).start()]
        return html
    def _get_type(self, html):
        html = re.search("<a href='/category/\\d+.html' id='cat\\d+' class=\"nav\" >(.*?)</a>", html).group(1)
        return html
    def _get_last_update(self, html):
        m = re.search("更新: <b><font.*?>(.*?)</font>",html)
        html = m.group(1) if (m) else ""
        return html
    def _get_last_chapter(self, html):
        html = re.findall("<a href='/novel/.*?' >(.*?)</a>", html)
        html = html[-1]
        return html
    def _get_chapters_url(self, html):
        try:
            html = re.sub("href='(/novel/.*?)' ", "href=\"" + self.chapter_url.replace("/novel/\\d", "") + "\\1\"", html)
            html = re.findall("<a href=\"(.*?/novel/.*?)\">.*?</a>", html)
            return html
        except:
            return []
    def _get_chapters_title(self, html):
        try:
            html = re.findall("<a href='/novel/.*?' >(.*?)</a>", html)
            return html
        except:
            return []
    def _get_content(self, html):
        try:
            html = html[re.search("<p class=content>", html).end():]
            html = html[:re.search("</p>", html).start()]
            if ("<br>\r\n" in html[:-10]):
                html = re.sub("<br>\r\n", "\n", html)
            else:
                html = re.sub(" ", "\n", html)
            if(len(html)<10):print("this chapter is too short!!!")
            return html
        except:
            return None

class Site(models.Site):
    def get_book(self, num):
        base_url = self.meta_base_url.format(num)
        download_url = self.meta_download_url.format(num)
        return Book(self.db, base_url, download_url, self.chapter_url, self.site, num, self.decode, self.max_thread, self.timeout)

desktop_setting = models.Setting(
    meta_base_url="https://www.book100.com/book/book{}.html",
    meta_download_url="https://www.book100.com/book/book{}.html",
    chapter_url="https://www.book100.com/novel/\\d"
)

mobile_setting = models.Setting(
    meta_base_url="https://www.book100.com/book/book{}.html",
    meta_download_url="https://www.book100.com/book/book{}.html",
    chapter_url="https://www.book100.com/novel/\\d"
)