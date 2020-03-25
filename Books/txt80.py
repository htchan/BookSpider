import re
try:import ClassDefinition
except:import Books.ClassDefinition as ClassDefinition

class Txt80Book(ClassDefinition.Book):
    def _cut_name(self,c):
        c = c[re.search("titlename",c).end():]
        c = c[re.search("<h1>",c).end():re.search("</h1>",c).start()]
        c = re.sub("全文阅读","",c)
        return c
    def _cut_writer(self,c):
        c = c[re.search("作者：",c).end():]
        c = c[re.search(">",c).end():re.search("</a>",c).start()]
        c = re.sub("'","",c)
        return c
    def _cut_date(self,c):
        m = re.search("\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}",c)
        c = c[m.start():m.end()]
        return c
    def _cut_last_chapter(self,c):
        while(re.search("<li>",c)):
            m = re.search("<li>",c)
            c = c[m.end():]
        c = c[:re.search("</li>",c).start()]
        c = re.sub("(<a.*?>|</a>)","",c)
        return c
    def _cut_book_type(self,c):
        c = c[re.search("分类：",c).end():]
        c = c[re.search(">",c).end():re.search("</a>",c).start()]
        return c
    def _cut_chapter(self,c):
        c = c[re.search("yulan..",c).end():]
        c = c[:re.search("</div>",c).start()]
        c = re.findall("(http.*?html)",c)
        return c
    def _cut_title(self,c):
        c = c[re.search("yulan..",c).end():]
        c = c[:re.search("</div>",c).start()]
        c = re.sub("<strong>.*?</strong>","",c)
        c = re.sub(".*>(.+?)<(a|\/a).*","\\1", c).split("\n")
        c = c[1:len(c)-1]
        return c
    def _cut_chapter_title(self,c):
        return ""
    def _cut_chapter_content(self,c):
        c = c[re.search("id=\"content\">",c).end():]
        c = c[:re.search("<div",c).start()]
        c = re.sub("&nbsp;"," ",c)
        c = re.sub("<br />    ","\n",c)
        c = c.replace("\r\n\r\n", "\r\n")
        return c
txt80 = {
    "book":Txt80Book,
    "identify":"80txt",
    "web":{
        "base_web":"https://www.balingtxt.com/txtml_{}.html",
        "download_web":"http://www.balingtxt.com/txtml_{}.html",
        "chapter_web":"http://www.xqiushu.com"
    },
    "setting":{
        "decode":"utf8",
        "timeout":30
    }
}

def site():
    return ClassDefinition.BookSite(**txt80)

