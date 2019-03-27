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

def mainLoop():
    while(true):
        print("Book download"+"-"*20)
        print("D : download books")
        print("C : check book end")
        print("N : check new books")
        print("E : exit")
        ans = input(">>>")
        if(ans.upper()=="E"):
            break
        if(ans.upper()=="D"):
            i=0
            for row in c.execute("select * from books where end = 'true' and download = 'false' and read = 'Null'"):
                book=Book(row[0],row[1],row[2],row[3],row[4])
                print(str(i)+":", end="\t")
                txt80.download(conn,book)
                i+=1

connect(os.getcwd()+"\\bookDownload.db")
i=0
for row in c.execute("select * from books where end = 'true' and download = 'false' and read = 'Null'"):
    book=Book(row[0],row[1],row[2],row[3],row[4])
    print(time.ctime()[11:-8],end="\t")
    print(str(i)+":", end="\t")
    txt80.download(conn,book)
    i+=1