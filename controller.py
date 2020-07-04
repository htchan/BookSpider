#!/usr/bin/env python3
import sys, gc, threading

from Books2.subBook import ck101, bestory, hjwzw, txt80, xqishu

download_path = '/mnt/addition/download/Books'
sites = {
    'BESTORY': bestory.Site('BESTORY', __file__.replace('controller.py', 'Books2/database/bestory.db'),
                            download_path+'/bestory/', bestory.desktop_setting, 'big5'),
    'CK101': ck101.Site('CK101', __file__.replace('controller.py', 'Books2/database/ck101.db'), 
                        download_path+'/ck101/', ck101.desktop_setting, 'big5'),
    'HJWZW': hjwzw.Site('HJWZW', __file__.replace('controller.py', 'Books2/database/hjwzw.db'),
                        download_path+'/hjwzw/', hjwzw.desktop_setting, 'utf8'),
    'TXT80': txt80.Site('TXT80', __file__.replace('controller.py', 'Books2/database/txt80.db'),
                        download_path+'/txt80/', txt80.desktop_setting, 'utf8'),
    'XQISHU': xqishu.Site('XQISHU', __file__.replace('controller.py', 'Books2/database/xqishu.db'),
                            download_path+'/xqishu/', xqishu.desktop_setting, 'utf8'),
    
}

MAX_ERROR_COUNT = 500

def print_help():
    print("Command: ")
    print("help" + " "*14 + "show the functin list avaliable")
    print("download" + " "*10 + "download books")
    print("update" + " "*12 + "update books information")
    print("explore" + " "*11 + "explore new books in internet")
    print("check" + " "*13 + "check recorded books finished")
    print("error" + " "*13 + "update all website may have error")
    print("backup" + " "*12 + "backup the current database by the current date and time")
    print("regular" + " "*11 + "do the default operation (explore->update->download->check)")
    print("moveLog" + " "*11 + "move the logging.log to another suitable filename")
    print("\n")
    print("Flags: ")
    print("--site=site" + " "*7 + "set specific site for download")

def explore():
    site_threads = []
    for (key, site) in sites.items():
        print(key, 'explore')
        thread = threading.Thread(target=site.explore, args=(MAX_ERROR_COUNT))
        site_threads.append(thread)
        thread.daemon = True
        thread.start()
    for thread in site_threads:
        thread.join()
def update():
    site_threads = []
    for (key, site) in sites.items():
        print(key, 'update')
        thread = threading.Thread(target=site.update, args=())
        site_threads.append(thread)
        thread.daemon = True
        thread.start()
    for thread in site_threads:
        thread.join()
def download():
    for (key, site) in sites.items():
        print(key, 'download')
        site.download()
def update_error():
    for (key, site) in sites.items():
        print(key, 'update error')
        site.update_error()
def check_end():
    for (key, site) in sites.items():
        print(key, 'check end')
        site.db.check_end()
def backup():
    for (key, site) in sites.items():
        print(key, 'database backup')
        site.db.backup()
def regular():
    print("backup badtabase")
    backup()
    print("explore" + "*"*30)
    explore()
    gc.collect()
    print("update" + "*"*30)
    update()
    gc.collect()
    print("error update" + "*"*30)
    update_error()
    gc.collect()
    print("download" + "*"*30)
    download()
    gc.collect()
    print("check end" + "*"*30)
    check_end()
def info():
    for (key, site) in sites.items():
        site.info()
def fix_error():
    for (key, site) in sites.items():
        if (key == 'BESTORY'):
            continue
        print(key, 'fix error')
        site.fix_database_error()

if (__name__ == "__main__"):
    args = sys.argv[1:]
    funct = {
        "help":        print_help,
        "download":    download,
        "update":      update,
        "explore":     explore,
        "check":       check_end,
        "error":       update_error,
        "regular":     regular,
        #"find":        find,
        "backup":      backup,
        "info":        info,
        "fix":         fix_error,
    }
    try:
        funct = funct.get(args[0])
        if(funct):
            funct()
    except IndexError:
        print("No arguement")
        print_help()
        exit()
    except KeyboardInterrupt:
        exit("Sudden Exit")
