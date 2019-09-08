try:import ClassDefinition
except:import Books.ClassDefinition as ClassDefinition
import os, re, datetime
import sqlite3

try:import ck101, txt80
except:from Books import ck101, txt80

### const init
MAX_EXPLORE_NUM = 100
dbPath = os.getcwd()
conn = sqlite3.connect(dbPath+"\\spider.db",check_same_thread=False)
path = 'download books'

txt80.txt80['conn'] = conn
txt80.txt80['path'] = path

ck101.ck101['conn'] = conn
ck101.ck101['path'] = path

### variable init
sites = {}
sites["ck101"] = ck101.site()
sites["80txt"] = txt80.site()

def __print_help(out):
    out("--help"+" "*14+"show the functin list avaliable")
    out("--download"+" "*10+"download books")
    out("--update"+" "*12+"update books information")
    out("--explore"+" "*11+"explore new books in internet")
    out("--check"+" "*13+"check recorded books finished")
    out("--error"+" "*13+"update all website may have error")
def download(out):
    for x in sites:
        sites[x].download(out)
    out("Download finish")
def update(out):
    for x in sites:
        sites[x].update(out)
    out("Update finish")
def explore(out,n=MAX_EXPLORE_NUM):
    for x in sites:
        sites[x].explore(n,out)
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
    sql += "date < '"+str(datetime.datetime.now().year-2)+"') and (end <> 'true' or end is null)"
    c = conn.cursor()
    c.execute(sql)
    conn.commit()
    out(str(c.rowcount)+" row affected")
def error_update(out):
    for x in sites:
        sites[x].error_update(out)
    print("Error update finished")
def find(out,*args):
    # return basic info of the books
    query = {}
    for element in args:
        element = element[2:].split('=')
        query[element[0]] = element[1]
    if("site" in query):
        result = sites[query["site"]].query(**query)
        print(str(sites[query["site"]])+'-'*30)
        for r in result:
            print(str(r["num"])+'\t'+r["writer"]+'\t'+r["name"])
        return
    for x in sites:
        result = sites[x].query(**query)
        print(str(sites[x])+'-'*30)
        for r in result:
            print(str(r["num"])+'\t'+r["writer"]+'\t'+r["name"])

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
