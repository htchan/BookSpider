import urllib.request
import os
import zipfile
import sqlite3
import http
import gzip
import io
import ClassDefinition
import time

class Novel101():
    class Book(ClassDefinition.BaseBook):
        def _getBasicInfo(self):
            # fill back the info by the website
            try:
                res = urllib.request.urlopen(self._website,timeout=60)
                content = res.read()

                # decode the content
                if (res.info().get('Content-Encoding') == 'gzip'):
                    gzipFile = gzip.GzipFile('','rb',9,io.BytesIO(content))
                    content = gzipFile.read()
                content = content.decode("big5","ignore")
                res.close()

            # return false if the webpage is not exist or not avaliable
            except (urllib.error.HTTPError, urllib.error.URLError, urllib.request.socket.timeout):
                return False

            if(not self._name):
                # get name
                start = content.find("f20h")
                self._name = content[start+6:]
                end = self._name.find("<em>")
                self._name = self._name[:end].strip()
                self._updated = True
            if(not self._writer):
                # get writer (writer)
                start = content.find("作者")
                self._writer = content[start+3:]
                end = self._writer.find('/em>')
                self._writer = self._writer[:end-1].strip()
                if(self._writer.replace("\x00","") == ""):
                    self._writer = "EMPTY"
                self._updated = True

            # get date (always get)
            start = content.find('更新時間')
            date = content[start+5:]
            start = date.find('</b>')
            date = date[start+4:]
            end = date.find('</')
            date = date[:end]
            if(self._date != date):
                self._date = date.strip()
                self._updated = True

            # get chapter (always get)
            if(self._chapter == ""):
                self._chapter = "EMPTY"
                self._updated = True
            if(not self._bookType):
                # check type (bookType)
                start = content.find('小說分類')
                self._bookType = content[start+3:]
                start = self._bookType.find('</b>')
                self._bookType = self._bookType[start+4:]
                end = self._bookType.find('</td>')
                self._bookType = self._bookType[:end].strip()
                if(self._bookType == ''):
                    self._bookType = "EMPTY"
                self._updated = True
            return self._updated
        def DownloadBook(self,path,out=print):
            # fill back the info by the website
            res = urllib.request.urlopen(self._website.replace("txt","ck101")[:-5]+"/",timeout=60)
            content = res.read()
            # decode the content
            if (res.info().get('Content-Encoding') == 'gzip'):
                gzipFile = gzip.GzipFile('','rb',9,io.BytesIO(content))
                content = gzipFile.read()
            content = content.decode("big5","ignore")
            res.close()

            # get chapter set
            start = content.find("defaulthtml4")
            chapters = content[start:]
            start = chapters.find("<TBODY>")
            chapters = chapters[start:]
            end = chapters.find("</TBODY>")
            chapters = chapters[:end]
            self._chapterSet = chapters.split("href=")
            out(self._name)
            # download chapters one by one
            for chapter in self._chapterSet:
                if("html" in chapter):
                    # go to the website and download
                    chapter = chapter[chapter.find('"')+1:]
                    chapter = chapter[:chapter.find('"')]
                    self._DownloadChapter("https://www.101novel.net"+chapter)
                    # log for progress
                    out("\r"+self._chapter,end=" "*20)
                    time.sleep(5)
            # save it into file
            try: os.mkdir(path)
            except: pass
            try: os.mkdir(path+"\\"+self._bookType)
            except: pass
            f = open(path+"\\"+self._bookType+"\\"+self._name+"-"+self._writer+".txt","w",encoding='utf8')
            f.write(self._text)
            f.close()
            return True
        def _DownloadChapter(self,url):
            # open chapter url
            chRes = urllib.request.urlopen(url,timeout=60)
            content = chRes.read()
            # decode the content
            if (chRes.info().get('Content-Encoding') == 'gzip'):
                gzipFile = gzip.GzipFile('','rb',9,io.BytesIO(content))
                content = gzipFile.read()
            content = content.decode("big5","ignore")
            chRes.close()
            # get the title
            self._text += "\n"
            start = content.find("<H1>")
            title = content[start+4:]
            end = title.find("</H1>")
            title = title[:end]
            self._chapter = title
            self._text += title + "\n"
            # get the content
            start = content.find('<P>')
            c = content[start+5:]
            end = c.find("</P>")
            c = c[:end]
            c = c.replace("&nbsp;"," ")
            c = c.replace("<br />","\n").replace("\r\n","\n")
            for i in range(3):
                c = c.replace("\n\n\n","\n\n")
            self._text += c
    def __init__(self, dbConn, path):
        self._webpage = ""
        self._bookNum = 0
        self.books = []
        self._conn = dbConn
        self._cursor = self._conn.cursor()
        self._path = path
    def Download(self,out=print):
        # put [end] book in db to books
        self.books.clear()
        for row in self._cursor.execute("select * from books where end='true' and download='false' and website like '%101novel%'"):
            self.books.append(self.Book(row[4],name=row[0],writer=row[1],date=row[2],chapter=row[3],bookType=row[5]))
            time.sleep(5)
        out("downloading")
        for book in self.books:
            if(book.DownloadBook(self._path)):
                self._cursor.execute("update books set download='true' where website='"+book._website+"'")
        out("finish download")
    def Update(self,out=print):
        # get all books from db to boo
        self.books.clear()
        for row in self._cursor.execute("select * from books where website like '%101novel%'"):
            self.books.append(self.Book(row[4],name=row[0],writer=row[1],date=row[2],chapter=row[3],bookType=row[5]))
            time.sleep(5)
        # check any update
        out("updating")
        for book in self.books:
            # if the book info had been updated
            if(book._updated):
                sql = "update books set date='"+book._date+"', chapter='"+book._chapter+"' where website='"+book._website+"'"
                self._cursor.execute(sql)
                self._conn.commit()
                # if it update, but ended, add a '-' at first of the file name
                flag = self._cursor.execute("select download from books where website='"+book._website+"'").fetchone()
                for f in flag:
                    if(f=="false"):
                        f = False
                        break
                if(f):
                    os.rename(self._path+"\\"+book._bookType+"\\"+book._name+"-"+book._writer+".txt",self._path+"\\"+book._bookType+"\\-"+book._name+"-"+book._writer+".txt")
                    book._name = '-'+book._name
                    self._cursor.execute("update books set name='"+book._name+"' where website='"+book._website+"'")
                    self._conn.commit()
        out("update finish")
    def Explore(self,n,out=print):
        # get the max book num from the db
        self.books.clear()
        self._bookNum = 0
        for row in self._cursor.execute("select website from books where website like '%101novel%' order by website desc"):
            i = row[0]
            i = int(i[i.rfind("/")+1:i.rfind(".")])
            if(i > self._bookNum):
                self._bookNum = i
        self._bookNum += 1
        errorPage = 0
        # check any new book by the book num and try to save it
        while(errorPage<n):
            time.sleep(5)
            b = self.Book("https://www.101novel.net/txt/"+str(self._bookNum)+".html")
            if(b._name):
                self.books.append(b)
                flag = bool(self._cursor.execute("select * from books where name='"+b._name+"' and writer='"+b._writer+"'").fetchone())
                if(not(flag)):
                    sql = (
                        "insert into books (name,writer,date,chapter,website,type,download) values"
                        "('"+b._name+"','"+b._writer+"','"+b._date+"','"+b._chapter+"','"+b._website+"','"+b._bookType+"','false')"
                    )
                    self._cursor.execute(sql)
                    self._conn.commit()
                errorPage = 0
            else: errorPage += 1
            self._bookNum += 1
