try:import ClassDefinition
except:import Books.ClassDefinition as ClassDefinition
import os, re, datetime
import sqlite3

### const init
MAX_EXPLORE_NUM = 100
conn = sqlite3.connect("spider.db",check_same_thread=False)
path = os.getcwd()

### class init
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
        c = re.findall("(https.*?html)",c)
        return c
    def _cut_title(self,c):
        c = c[re.search("yulan..",c).end():]
        c = c[:re.search("</div>",c).start()]
        c = re.sub("<strong>.*?</strong>","",c)
        c = re.findall(".*>(\w.+?)<.*",c)
        return c
    def _cut_chapter_title(self,c):
        return ""
    def _cut_chapter_content(self,c):
        c = c[re.search("id=\"content\">",c).end():]
        c = c[:re.search("<div",c).start()]
        c = re.sub("&nbsp;"," ",c)
        c = re.sub("<br />    ","\n",c)
        return c
txt80 = {
    "book":Txt80Book,
    "conn":conn,
    "path":path,
    "identify":"80txt",
    "web":{
        "base_web":"https://www.80txt.la/txtml_{}.html",
        "download_web":"https://www.80txt.la/txtml_{}.html",
        "chapter_web":"https://www.qiushuzw.com"
    },
    "setting":{
        "decode":"utf8",
        "timeout":30
    }
}

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
    "conn":conn,
    "path":path,
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

Txt80Site = ClassDefinition.BookSite(**txt80)
Ck101Site = ClassDefinition.BookSite(**ck101)

### variable init
sites = []
sites.append(Ck101Site)
sites.append(Txt80Site)

def __print_help(out):
    out("--help"+" "*14+"show the functin list avaliable")
    out("--download"+" "*10+"download books")
    out("--update"+" "*12+"update books information")
    out("--explore"+" "*11+"explore new books in internet")
    out("--check"+" "*13+"check recorded books finished")
    out("--error"+" "*13+"update all website may have error")
def download(out):
    for site in sites:
        site.download(out)
    out("Download finish")
def update(out):
    for site in sites:
        site.update(out)
    out("Update finish")
def explore(out,n=MAX_EXPLORE_NUM):
    for site in sites:
        site.explore(n,out)
    out("Explore finish")
def check_end(out):
    # update books end by their last chapter content
    criteria = ["后记", "後記", "新书", "新書", "结局", "結局", "感言", 
                "尾声", "尾聲", "终章", "終章", "外传", "外傳", "完本",
                "结束", "結束", "完結", "完结", "终结", "終結", "番外"
                "结尾", "結尾", "全书完", "全書完"]
    sql = "update books set end='true' where ("
    for c in criteria:
        sql += "chapter like '%"+c+"%' or "
    sql += "date < '"+str(datetime.datetime.now().year-2)+"') and end <> 'true'"
    c = conn.cursor()
    c.execute(sql)
    conn.commit()
    out(str(c.rowcount)+" row affected")
def error_update(out):
    for site in sites:
        site.error_update(out)
    print("Error update finished")
def find(out,*args):
    # return basic info of the books
    query = {}
    for element in args:
        element = element[2:].split('=')
        query[element[0]] = element[1]
    for site in sites:
        result = site.query(**query)
        print(str(site)+'-'*30)
        for r in result:
            print(re.search("(\\d*?)\\.html",r[4])[1]+'\t'+r[1]+'\t'+r[0])

if(__name__=="__main__"):
    # cmd interface
    import sys
    args = sys.argv[1:]
    funct = {
        "--help":__print_help,
        "--download":download,
        "--update":update,
        "--explore":explore,
        "--check":check_end,
        "--error":error_update,
        "--find":find
    }
    try:
        funct = funct.get(args[0])
        if(funct):
            if(funct==explore):
                if(len(args)==1):funct(print)
                elif((len(args)==2)and(args[1].isdigit())):funct(print,int(args[1]))
                else:exit("Invalid arguement")
            elif(funct==find):
                funct(print,*args[1:])
            elif(len(args)==1):
                funct(print)
            else:exit("Invalid arguement")
    except IndexError:exit("No arguement")
    except KeyboardInterrupt:exit("Sudden Exit")
