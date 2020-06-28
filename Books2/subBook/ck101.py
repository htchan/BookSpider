from .. import models
import re

class Book(models.Book):
    def _get_title(self, html):
        html = html[re.search("<h1>", html).end():]
        html = html[re.search("\">", html).end():]
        html = html[:re.search("</a>", html).start()]
        html = re.sub("('|\x00)", "", html)
        return html
    def _get_writer(self, html):
        html = html[re.search("作者︰", html).end():]
        html = html[re.search("\">", html).end():re.search("</a>", html).start()]
        html = re.sub("('|\x00)", "", html)
        return html
    def _get_type(self, html):
        html = re.sub("\x00", "", html)
        html = re.findall(" &gt; (\w{4}) &gt; ", html)[0]
        return html
    def _get_last_update(self, html):
        m = re.search("最新章節\\((\\d{4}-\\d{2}-\\d{2})\\)", html)
        html = m.group(1) if(m) else ""
        return html
    def _get_last_chapter(self, html):
        html = html[re.search("<strong>", html).end():]
        html = html[re.search(">", html).end():re.search("</a>", html).start()]
        html = re.sub("('|\x00)", "", html)
        return html
    def _get_chapters_url(self, html):
        try:
            html = html[re.search("<dl", html).end():]
            html = html[:re.search("</div>", html).start()]
            html = re.sub("\x00", "", html)
            html = re.sub("\"(/.*?html)", "\"" + self.chapter_url.replace("/\\d", "") + "\\1", html)
            html = re.findall("(https://w+.ck101.org/.*?html)", html)
            return html
        except:
            return []
    def _get_chapters_title(self, html):
        try:
            html = html[re.search("<dl", html).end():]
            html = html[:re.search("</div>", html).start()]
            html = re.sub("\x00", "", html)
            html = re.findall("html\">(.*?)<", html)
            return html
        except:
            return []
    def _get_content(self, html):
        html = html[re.search("yuedu_zhengwen", html).end():]
        html = html[re.search(">", html).end():]
        html = html[:re.search("</div", html).start()]
        html = re.sub("&nbsp;", " ", html)
        html = re.sub("(<br />|<.+?</.+?>|\x00)", "", html)
        if(len(html)<10):print("this chapter is too short!!!")
        return html

class Site(models.Site):
    def get_book(self, num):
        base_url = self.meta_base_url.format(num)
        download_url = self.meta_download_url.format(num)
        return Book(self.db, base_url, download_url, self.chapter_url, self.site, num, self.decode, self.max_thread, self.timeout)

desktop_setting = models.Setting(
    meta_base_url="https://www.ck101.org/book/{}.html",
    meta_download_url="https://www.ck101.org/0/{}/",
    chapter_url="https://www.ck101.org/\\d"
)

mobile_setting = models.Setting(
    meta_base_url="https://w.ck101.org/book/{}.html",
    meta_download_url="https://w.ck101.org/0/{}/",
    chapter_url="https://w.ck101.org/\\d"
)
