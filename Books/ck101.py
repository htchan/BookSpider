import re
try:import ClassDefinition
except:import Books.ClassDefinition as ClassDefinition


class Ck101Book(ClassDefinition.Book):
    def _cut_name(self,c):
        c = c[re.search("<h1>",c).end():]
        c = c[re.search("\">",c).end():]
        c = c[:re.search("</a>",c).start()]
        c = re.sub("('|\x00)","",c)
        return c
    def _cut_writer(self,c):
        c = c[re.search("作者︰",c).end():]
        c = c[re.search("\">",c).end():re.search("</a>",c).start()]
        c = re.sub("('|\x00)","",c)
        return c
    def _cut_date(self,c):
        m = re.search("最新章節\\((\\d{4}-\\d{2}-\\d{2})\\)",c)
        c = m.group(1) if(m) else ""
        return c
    def _cut_last_chapter(self,c):
        c = c[re.search("<strong>",c).end():]
        c = c[re.search(">",c).end():re.search("</a>",c).start()]
        c = re.sub("('|\x00)","",c)
        return c
    def _cut_book_type(self,c):
        c = re.sub("\x00","",c)
        c = re.findall(" &gt; (\w{4}) &gt; ",c)[0]
        return c
    def _cut_chapter(self,c):
            c = c[re.search("<dl",c).end():]
            c = c[:re.search("</div>",c).start()]
            c = re.sub("\x00","",c)
            c = re.sub("\"(/.*?html)", "\"" + self._Book__chapter_web.replace("/\\d", "") + "\\1",c)
            c = re.findall("(https://w+.ck101.org/.*?html)",c)
            return c
    def _cut_title(self,c):
        try:
            c = c[re.search("<dl",c).end():]
            c = c[:re.search("</div>",c).start()]
            c = re.sub("\x00","",c)
            c = re.findall("html\">(.*?)<",c)
            return c
        except:
            return []
    def _cut_chapter_title(self,c):
        return ""
    def _cut_chapter_content(self,c):
        c = c[re.search("yuedu_zhengwen",c).end():]
        c = c[re.search(">",c).end():]
        c = c[:re.search("</div",c).start()]
        c = re.sub("&nbsp;"," ",c)
        c = re.sub("(<br />|<.+?</.+?>|\x00)","",c)
        if(len(c)<10):print("this chapter is too short!!!")
        return c

desktop_web = {
    "base_web":"https://www.ck101.org/book/{}.html",
    "download_web":"https://www.ck101.org/0/{}/",
    "chapter_web":"https://www.ck101.org/\\d"
}
tablet_web = {
    "base_web":"https://w.ck101.org/book/{}.html",
    "download_web":"https://w.ck101.org/0/{}/",
    "chapter_web":"https://w.ck101.org/\\d"
}

ck101 = {
    "book":Ck101Book,
    "identify":"ck101",
    'web':desktop_web,
    "setting":{
        "decode":"big5",
        "timeout":30
    }
}

def site():
    return ClassDefinition.BookSite(**ck101)
