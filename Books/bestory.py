import re
try:import ClassDefinition
except:import Books.ClassDefinition as ClassDefinition


class BestoryBook(ClassDefinition.Book):
    def _cut_name(self,c):
        c = c[re.search("<title>",c).end():]
        c = c[:re.search("</title>",c).start()]
        c = re.sub("('|\x00)","",c).replace(" 在線觀看", "")
        return c
    def _cut_writer(self,c):
        c = c[re.search("作者:",c).end():]
        c = c[re.search("\">",c).end():re.search("</font>",c).start()]
        return c
    def _cut_date(self,c):
        m = re.search("更新: <b><font.*?>(.*?)</font>",c)
        c = m.group(1) if (m) else ""
        return c
    def _cut_last_chapter(self,c):
        c = re.findall("<a href='/novel/.*?' >(.*?)</a>", c)
        c = c[-1]
        return c
    def _cut_book_type(self,c):
        c = re.search("<a href='/category/\\d+.html' id='cat\\d+' class=\"nav\" >(.*?)</a>", c).group(1)
        return c
    def _cut_chapter(self,c):
        c = re.sub("href='(/novel/.*?)' ", "href=\"" + self._Book__chapter_web.replace("/novel/\\d", "") + "\\1\"", c)
        c = re.findall("<a href=\"(.*?/novel/.*?)\">.*?</a>", c)
        return c
    def _cut_title(self,c):
        try:
            c = re.findall("<a href='/novel/.*?' >(.*?)</a>", c)
            return c
        except:
            return []
    def _cut_chapter_title(self,c):
        return ""
    def _cut_chapter_content(self,c):
        c = c[re.search("<p class=content>", c).end():]
        c = c[:re.search("</p>", c).start()]
        if ("<br>\r\n" in c[:-10]):
            c = re.sub("<br>\r\n", "\n", c)
        else:
            c = re.sub(" ", "\n", c)
        if(len(c)<10):print("this chapter is too short!!!")
        return c

desktop_web = {
    "base_web":"https://www.book100.com/book/book{}.html",
    "download_web":"https://www.book100.com/book/book{}.html",
    "chapter_web":"https://www.book100.com/novel/\d"
}
tablet_web = {
    "base_web":"https://www.book100.com/book/book98999.html",
    "download_web":"https://www.book100.com/book/book98999.html",
    "chapter_web":"https://www.book100.com/novel/\d"
}

bestory = {
    "book":BestoryBook,
    "identify":"bestory",
    'web':desktop_web,
    "setting":{
        "decode":"big5",
        "timeout":30
    }
}

def site():
    return ClassDefinition.BookSite(**bestory)
