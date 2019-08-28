import re
import ClassDefinition


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
        c = re.sub("\x00","",c)
        m = re.search("\\d{4}-\\d{2}-\\d{2}",c)
        c = c[m.start():m.end()]
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
        c = re.sub("\"(/.*?html)","\"https://www.ck101.org\\1",c)
        c = re.findall("(https://www.ck101.org/.*?html)",c)
        return c
    def _cut_title(self,c):
        c = c[re.search("<dl",c).end():]
        c = c[:re.search("</div>",c).start()]
        c = re.sub("\x00","",c)
        c = re.findall("html\">(.*?)<",c)
        return c
    def _cut_chapter_title(self,c):
        return ""
    def _cut_chapter_content(self,c):
        c = c[re.search("yuedu_zhengwen",c).end():]
        c = c[re.search(">",c).end():]
        c = c[:re.search("</div",c).start()]
        c = re.sub("&nbsp;"," ",c)
        c = re.sub("(<br />|<.+?</.+?>|\x00)","",c)
        if(len(c)<100):print("this chapter is too short!!!")
        return c
ck101 = {
    "book":Ck101Book,
    "identify":"ck101",
    "web":{
        "base_web":"https://www.ck101.org/book/{}.html",
        "download_web":"https://www.ck101.org/0/{}/",
        "chapter_web":"https://www.ck101.org/"
    },
    "setting":{
        "decode":"big5",
        "timeout":30
    }
}

def site():
    return ClassDefinition.BookSite(**ck101)