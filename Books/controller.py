#!/usr/bin/env python3

try:import ClassDefinition
except:import Books.ClassDefinition as ClassDefinition
import os, re, datetime
import sqlite3
import gc

try:import ck101, txt80, bestory, xqishu, hjwzw
except:from Books import ck101, txt80, bestory, xqishu, hjwzw

### load setting
dbPath = os.getcwd()
path = 'download books'
dbName = '/spider.db'
f = open("../.setting", 'r')
setting = f.readlines()
f.close()
for s in setting:
    exec(s)


### const init
MAX_EXPLORE_NUM = 100
conn = sqlite3.connect(dbPath+dbName,check_same_thread=False)

txt80.txt80['conn'] = conn
txt80.txt80['path'] = path

ck101.ck101['conn'] = conn
ck101.ck101['path'] = path

bestory.bestory['conn'] = conn
bestory.bestory['path'] = path

xqishu.xqishu['conn'] = conn
xqishu.xqishu['path'] = path

hjwzw.hjwzw['conn'] = conn
hjwzw.hjwzw['path'] = path

### variable init
sites = {}
sites["CK101"] = ck101.site()
sites["80TXT"] = txt80.site()
sites['BESTORY'] = bestory.site()
sites['XQISHU'] = xqishu.site()
sites['HJWZW'] = hjwzw.site()

def __get_flags(*args):
    output = {}
    for arg in args:
        if ((arg.find("--") == 0) and (arg.find("=") > 2)):
            output[arg[arg.find("--")+2:arg.find("=")].upper()] = arg[arg.find("=")+1:].upper()
    return output
def __print_help(out,*args):
    out("Command: ")
    out("help" + " "*14 + "show the functin list avaliable")
    out("download" + " "*10 + "download books")
    out("update" + " "*12 + "update books information")
    out("explore" + " "*11 + "explore new books in internet")
    out("check" + " "*13 + "check recorded books finished")
    out("error" + " "*13 + "update all website may have error")
    out("backup" + " "*12 + "backup the current database by the current date and time")
    out("regular" + " "*11 + "do the default operation (explore->update->download->check)")
    out("moveLog" + " "*11 + "move the logging.log to another suitable filename")
    out("\n")
    out("Flags: ")
    out("--site=site" + " "*7 + "set specific site for download")
def download(out,*args):
    if (len(args) == 0):
        for x in sites:
            sites[x].download(out)
            out("Download Finish")
            return
    flags = __get_flags(*args)
    if ("SITE" in flags):
        try:
            sites[flags["SITE"]].download(out)
            out("Download Finish")
        except IndexError:
            out("Site " + flags["SITE"] + " Not Found")
def update(out,*args):
    if (len(args) == 0):
        for x in sites:
            sites[x].update(False, out)
        out("Update Finish")
        return
    flags = __get_flags(*args)
    update_all = False
    if ("ALL" in flags):
        update_all = True
    if ("SITE" in flags):
        try:
            sites[flags["SITE"]].update(update_all, out)
            out("Update Finish")
        except IndexError:
            out("Site " + flags["SITE"] + " not found")
def explore(out,*args):
    if (len(args) == 0):
        n = MAX_EXPLORE_NUM
        for x in sites:
            sites[x].explore(n,out)
        out("Explore finish")
        return
    flags = __get_flags(*args)
    if ("SITE" in flags):
        site = flags["SITE"]
        n = int(flags["NUM"]) if (("NUM" in flags) and (flags["NUM"].isdigit())) else MAX_EXPLORE_NUM
        try:
            sites[site].explore(n, out)
        except IndexError:
            out("Site " + flags["SITE"] + " not found")

def check_end(out,*args):
    # update books end by their last chapter content
    criteria = ["后记", "後記", "新书", "新書", "结局", "結局", "感言", 
                "尾声", "尾聲", "终章", "終章", "外传", "外傳", "完本",
                "结束", "結束", "完結", "完结", "终结", "終結", "番外",
                "结尾", "結尾", "全书完", "全書完", "全本完"]
    sql = "update books set end='true', download='false' where ("
    for c in criteria:
        sql += "chapter like '%"+c+"%' or "
    sql += "date < '"+str(datetime.datetime.now().year-2)+"') and (end <> 'true' or end is null)"
    print(sql)
    c = conn.cursor()
    c.execute(sql)
    conn.commit()
    out(str(c.rowcount)+" row affected")
def error_update(out,*args):
    if (len(args) == 0):
        for x in sites:
            sites[x].error_update(out)
        out("Error update finished")
        return
    flags = __get_flags(*args)
    if ("SITE" in flags):
        try:
            sites[flags["SITE"]].error_update(out)
        except IndexError:
            out("Site " + flags["SITE"] + "not found")
'''
def find(out,*args):
    # return basic info of the books
    flags = __get_flags(*args)
    if ("SITE" in flags):
        result = sites[flags["site"]].query(**query)
        print(str(sites[query["site"]])+'-'*30)
        for r in result:
            print(str(r["num"])+'\t'+r["writer"]+'\t'+r["name"])
        return
    for x in sites:
        result = sites[x].query(**query)
        print(str(sites[x])+'-'*30)
        for r in result:
            print(str(r["num"])+'\t'+r["writer"]+'\t'+r["name"])
'''
def backup(out,*args):
    original_database = open(dbPath + dbName, "rb").read()
    flags = __get_flags(*args)
    destination = flags["DEST"] if ("DEST" in flags) else "./backup/"
    open(destination+str(datetime.datetime.now())+"_backup.db", "wb").write(original_database)
def regular(out,*args):
    out("backup badtabase")
    backup(out, *args)
    out("explore" + "*"*30)
    explore(out, *args)
    gc.collect()
    out("update" + "*"*30)
    update(out, *args)
    gc.collect()
    out("error update" + "*"*30)
    error_update(out, *args)
    gc.collect()
    out("download" + "*"*30)
    download(out, *args)
    gc.collect()
    out("check end" + "*"*30)
    check_end(out, *args)
def moveLog(out,*args):
    if (os.path.exists("./log/logging.log")):
        log_files = filter(lambda filename: int('0' + filename.replace("logging", "").replace(".txt", "").replace(".log", "")), os.listdir("./log"))
        log_number = max(log_files) + 1
        os.rename("./log/logging.txt", "./log/logging{}.txt".format(log_number))
        print("log file moved")
    else:
        print("no current log file found")

if(__name__=="__main__"):
    # cmd interface
    import sys
    args = sys.argv[1:]
    funct = {
        "help":        __print_help,
        "download":    download,
        "update":      update,
        "explore":     explore,
        "check":       check_end,
        "error":       error_update,
        "regular":     regular,
        #"find":        find,
        "backup":      backup,
        "moveLog":     moveLog,
    }
    try:
        funct = funct.get(args[0])
        if(funct):
            funct(print, *args[1:])
    except IndexError:
        print("No arguement")
        __print_help(print)
        exit()
    except KeyboardInterrupt:
        exit("Sudden Exit")
