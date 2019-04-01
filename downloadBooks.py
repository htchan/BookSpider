import sqlite3
import os
import txt80
import time
import ClassDefinition

def connect(path):
    global conn, c
    conn = sqlite3.connect(path)
    c = conn.cursor()

def disconnect():
    c.close()
    conn.commit()
    conn.close()

def downloadAll(sql):
    i = 0
    for row in c.execute(sql):
        book = ClassDefinition.Book(row[0],row[1],row[2],row[3],row[4])
        print(time.ctime()[11:-8],end="\t")
        print(str(i)+":", end="\t")
        if("80txt" in book.website):
            txt80.download(conn,book)
        i += 1
        time.sleep(5)

def checkNew():
    txt80.anyNew(conn)

def updateAll(sql):
    for row in c.execute(sql):
        book = ClassDefinition.Book(row[0],row[1],row[2],row[3],row[4])
        print(book.name, end="\t")
        if("80txt" in book.website):
            txt80.bookUpdate(conn, book)
    # update books end by their last chapter content
    sql = (
        "update books set end='true' where "
        "chapter like '%后记%' or chapter like '%新书%' or "
        "chapter like '%结局%' or chapter like '%感言%' or "
        "chapter like '%完本%' or chapter like '%尾声%' or "
        "chapter like '%终章%' or chapter like '%结束%' or "
        "chapter like '%外传%'"
    )
    c.execute(sql)

def showInfo():
    print("number of books end, but not downloaded :")
    print(str(len(c.execute("select website from books where end='true' and not download='true'").fetchall())))
    print("number of books is not end :")
    print(str(len(c.execute("select website from books where end is null").fetchall())))
    print("number of book not end and last update is at least last year :")
    print(str(len(c.execute("select website from books where date not like '%"+time.ctime()+"%'").fetchall())))


def mainLoop():
    while(True):
        print("Book download"+"-"*20)
        print("D : download books")
        print("U : check book update")
        print("N : check new books")
        print("I : information")
        print("E : exit")
        ans = input(">>> ")
        if(ans.upper()=="E"):
            disconnect()
            break
        elif(ans.upper()=="D"):
            downloadAll("select * from books where end = 'true' and download = 'false' and read = 'false' order by date desc")
        elif(ans.upper()=="N"):
            checkNew()
        elif(ans.upper()=="U"):
            updateAll("select * from books where end='false' order by date desc")
        elif(ans.upper()=="I"):
            showInfo()
            input()


connect(os.getcwd()+"\\bookDownload.db")
mainLoop()