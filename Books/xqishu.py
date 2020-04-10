import re
try:import ClassDefinition
except:import Books.ClassDefinition as ClassDefinition


class XqishuBook(ClassDefinition.Book):
    def _cut_name(self,c):
        c = c[re.search("<div class=\"tit1\">", c).end():]
        c = re.search("<h1>(.*?)</h1>", c).group(1)
        return c
    def _cut_writer(self,c):
        c = re.search("<span>小说作者：(.*?)</span>", c).group(1)
        return c
    def _cut_date(self,c):
        c = re.search("<span>更新日期：(.*?)</span>", c).group(1)
        return c
    def _cut_last_chapter(self,c):
        c = c[re.search("<div class=\"new_con\">", c).end():]
        c = c[re.search("最新章节：", c).end():]
        c = re.search("<a.*?>(.*?)</a>", c).group(1)
        return c
    def _cut_book_type(self,c):
        c = c[re.search("<div class=\"crumbs\">", c).end():]
        c = re.search("<a href=\"/ls/.*?html\">(.*?)</a>", c).group(1)
        return c
    def _cut_chapter(self,c):
        c = c[re.search("<div class=\"book_con_list\">", c).end():]
        c = c[re.search("<div class=\"book_con_list\">", c).end():]
        c = re.sub("(&nbsp;|<ul>|</ul>|<li>|</li>)", "", c[:re.search("</div>", c).start()])
        c = re.sub("</(.*?)>", "</\\1>\n", c).split("\n")
        for i in range(len(c)):
            if (c[i].find("<a") == 0):
                if (c[i].find("<a href") < 0):
                    continue
                c[i] = self._Book__download_web + re.search("<a href=\"(.*?)\"", c[i]).group(1)
            else:
                c[i] = ""
        return c
    def _cut_title(self,c):
        c = c[re.search("<div class=\"book_con_list\">", c).end():]
        c = c[re.search("<div class=\"book_con_list\">", c).end():]
        c = re.sub("(&nbsp;|<ul>|</ul>|<li>|</li>)", "", c[:re.search("</div>", c).start()])
        c = re.sub("(</?h\d>|</?a.*?>)", "", re.sub("</(.*?)>", "</\\1>\n", c)).split("\n")
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

xqishu = {
    "book":XqishuBook,
    "identify":"xqishu",
    'web':{
        "base_web":"http://www.xqiushu.com/txt{}/",
        "download_web":"http://www.xqiushu.com/t/{}/",
        "chapter_web":"http://www.xqiushu.com/t/\d"
    },
    "setting":{
        "decode":"utf8",
        "timeout":30
    }
}

def site():
    return ClassDefinition.BookSite(**xqishu)
