import re
try:import ClassDefinition
except:import Books.ClassDefinition as ClassDefinition


class HjwzwBook(ClassDefinition.Book):
    def _cut_name(self,c):
        c = re.search("<meta property=\"og:novel:book_name\" content=\"(.*?)\" />", c).group(1)
        return c
    def _cut_writer(self,c):
        c = re.search("<meta property=\"og:novel:author\" content=\"(.*?)\" />", c).group(1)
        return c
    def _cut_date(self,c):
        c = re.search("<meta property=\"og:novel:update_time\" content=\"(.*?)\" />", c).group(1)
        return c
    def _cut_last_chapter(self,c):
        c = re.search("<meta property=\"og:novel:latest_chapter_name\" content=\"(.*?)\" />", c).group(1)
        return c
    def _cut_book_type(self,c):
        c = re.search("<meta property=\"og:novel:category\" content=\"(.*?)\" />", c).group(1)
        return c
    def _cut_chapter(self,c):
        c = re.sub("/Book/Read/(.*?)", self._Book__chapter_web.replace("\\d", "") + "\\1", c)
        c = re.findall("<a href=\"(" + self._Book__chapter_web.replace("\\d", "") + ".*?)\".*?>", c)
        return c
    def _cut_title(self,c):
        c = re.findall("<a href=\"/Book/Read.*?>(.*?)</a>", c)
        return c
    def _cut_chapter_title(self,c):
        return ""
    def _cut_chapter_content(self,c):
        c = c[re.search("<table align=\"center\" width=\"1000px\">", c).end():]
        c = c[re.search("<div id=\"AllySite\"", c).end():]
        c = c[re.search("<div (.*?)>", c).end():]
        c = re.sub("(&nbsp;|<p ?/?>|</?a(.*?)>|<br/>|                    )", "", c[:re.search("</div>", c).start()])
        c = c.replace("\r\n", "\r\n\r\n").strip()
        return c

book_factory_info = {
    "base_web":"https://tw.hjwzw.com/Book/{}/",
    "download_web":"https://tw.hjwzw.com/Book/Chapter/{}/",
    "chapter_web":"https://tw.hjwzw.com/Book/Read/\d",
    "book_product":HjwzwBook,
    "decode":"utf8",
    "timeout":30
}

hjwzw = {
    "book_factory":ClassDefinition.BookFactory(**book_factory_info),
    "identify":"hjwzw",
}
'''
hjwzw = {
    "book":HjwzwBook,
    "identify":"hjwzw",
    'web':{
        "base_web":"https://tw.hjwzw.com/Book/{}/",
        "download_web":"https://tw.hjwzw.com/Book/Chapter/{}/",
        "chapter_web":"https://tw.hjwzw.com/Book/Read/\d"
    },
    "setting":{
        "decode":"utf8",
        "timeout":30
    }
}
'''
def site():
    return ClassDefinition.BookSite(**hjwzw)
