import sqlite3
import os
import txt80
import time

class Book():
    def __init__(self, name='', writer='', date='', chapter='', website='', bookType=''):
        self.name = name
        self.writer = writer
        self.date = date
        self.chapter = chapter
        self.website = website
        self.bookType = bookType


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
        book = Book(row[0],row[1],row[2],row[3],row[4])
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
        book = Book(row[0],row[1],row[2],row[3],row[4])
        print(book.name)
        if("80txt" in book.website):
            txt80.bookUpdate(conn, book)


def mainLoop():
    while(True):
        print("Book download"+"-"*20)
        print("D : download books")
        print("U : check book update")
        print("N : check new books")
        print("E : exit")
        ans = input(">>> ")
        if(ans.upper()=="E"):
            disconnect()
            break
        elif(ans.upper()=="D"):
            downloadAll("select * from books where end = 'true' and download = 'false' and read = 'Null' order by date")
        elif(ans.upper()=="N"):
            checkNew()
        elif(ans.upper()=="U"):
            updateAll("select * from books where end is null")


connect(os.getcwd()+"\\bookDownload.db")
mainLoop()
'''
i=0
for row in c.execute("select * from books where end = 'true' and download = 'false' and read = 'Null'"):
    book = Book(row[0],row[1],row[2],row[3],row[4])
    print(time.ctime()[11:-8],end="\t")
    print(str(i)+":", end="\t")
    txt80.download(conn,book)
    i += 1
'''