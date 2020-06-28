from .. import models
import re

class Book(models.Book):
    def _get_title(self, html):
        html = re.search("<meta property=\"og:novel:book_name\" content=\"(.*?)\" />", html).group(1)
        return html
    def _get_writer(self, html):
        html = re.search("<meta property=\"og:novel:author\" content=\"(.*?)\" />", html).group(1)
        return html
    def _get_type(self, html):
        html = re.search("<meta property=\"og:novel:category\" content=\"(.*?)\" />", html).group(1)
        return html
    def _get_last_update(self, html):
        html = re.search("<meta property=\"og:novel:update_time\" content=\"(.*?)\" />", html).group(1)
        return html
    def _get_last_chapter(self, html):
        html = re.search("<meta property=\"og:novel:latest_chapter_name\" content=\"(.*?)\" />", html).group(1)
        return html
    def _get_chapters_url(self, html):
        try:
            html = re.sub("/Book/Read/(.*?)", self.chapter_url.replace("\\d", "") + "\\1", html)
            html = re.findall("<a href=\"(" + self.chapter_url.replace("\\d", "") + ".*?)\".*?>", html)
            return html
        except:
            return []
    def _get_chapters_title(self, html):
        try:
            html = re.findall("<a href=\"/Book/Read.*?>(.*?)</a>", html)
            return html
        except:
            return []
    def _get_content(self, html):
        html = html[re.search("<table align=\"center\" width=\"1000px\">", html).end():]
        html = html[re.search("<div id=\"AllySite\"", html).end():]
        html = html[re.search("<div (.*?)>", html).end():]
        html = re.sub("(&nbsp;|<p ?/?>|</?a(.*?)>|<br/>|                    )", "", html[:re.search("</div>", html).start()])
        html = html.replace("\r\n", "\r\n\r\n").strip()
        return html

class Site(models.Site):
    def get_book(self, num):
        base_url = self.meta_base_url.format(num)
        download_url = self.meta_download_url.format(num)
        return Book(self.db, base_url, download_url, self.chapter_url, self.site, num, self.decode, self.max_thread, self.timeout)

desktop_setting = models.Setting(
    meta_base_url="https://tw.hjwzw.com/Book/{}/",
    meta_download_url="https://tw.hjwzw.com/Book/Chapter/{}/",
    chapter_url="https://tw.hjwzw.com/Book/Read/\\d"
)

mobile_setting = models.Setting(
    meta_base_url="https://tw.hjwzw.com/Book/{}/",
    meta_download_url="https://tw.hjwzw.com/Book/Chapter/{}/",
    chapter_url="https://tw.hjwzw.com/Book/Read/\\d"
)
