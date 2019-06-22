# -*- coding: utf-8 -*-
import sqlite3
import os
try:
    import txt80, hjwxw, ck101
except:
    import Books.txt80 as txt80
    import Books.hjwxw as hjwxw
    import Books.ck101 as ck101

### new (with oop)
class BookCollection:
    def __init__(self,path):
        # connect to database
        self.directory = path
        self.conn = sqlite3.connect(self.directory+"\\spider.db")
        self.c = self.conn.cursor()
        self.websites = []
        # init all book website
        self.websites.append(txt80.TXT80(self.conn,self.directory))
        self.websites.append(hjwxw.HJWXW(self.conn,self.directory))
        self.websites.append(ck101.CK101(self.conn,self.directory))

    def close(self):
        self.conn.commit()
        self.conn.close()

    def Download(self):
        # for all website, download books
        for website in self.websites:
            website.Download()

    def Update(self):
        for website in self.websites:
            website.Update()

    def Explore(self,n):
        for website in self.websites:
            print(type(website))
            website.Explore(n)

    def ErrorUpdate(self):
        for website in self.websites:
            website.ErrorUpdate()

'''
    # update books end by their last chapter content
    sql = (
        "update books set end='true' where "
        "chapter like '%后记%' or chapter like '%新书%' or "
        "chapter like '%结局%' or chapter like '%感言%' or "
        "chapter like '%完本%' or chapter like '%尾声%' or "
        "chapter like '%终章%' or chapter like '%结束%' or "
        "chapter like '%外传%'"
    )
'''