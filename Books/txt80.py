import urllib.request
import os
import zipfile
import sqlite3
import http
import gzip
import io
try: import ClassDefinition
except: import Books.ClassDefinition as ClassDefinition

class TXT80():
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
                content = content.decode("utf-8")
                res.close()

            # return false if the webpage is not exist or not avaliable
            except (urllib.error.HTTPError, urllib.request.socket.timeout):
                return False

            if(not self._name):
                # get name
                start = content.find("titlename")
                self._name = content[start+9:]
                start = self._name.find("<h1>")
                self._name = self._name[start+4:]
                end = self._name.find("</h1>")
                self._name = self._name[:end]
                self._name = self._name.replace("全文阅读","")
                self._updated = True
            if(not self._writer):
                # get writer (writer)
                start = content.find("作者：")
                self._writer = content[start+3:]
                start = self._writer.find('>')
                self._writer = self._writer[start+1:]
                end = self._writer.find('</a>')
                self._writer = self._writer[:end]
                self._updated = True

            # get date (always get)
            start = content.find('更新时间：')
            date = content[start+5:]
            end = date.find('</span>')
            date = date[:end]
            if(self._date != date):
                self._date = date
                self._updated = True

            # get chapter (always get)
            start = content.rfind("<li>")
            chapter = content[start+4:]
            start = chapter.find('">')
            chapter = chapter[start+2:]
            end = chapter.find('</a>')
            chapter = chapter[:end]
            if(self._chapter != chapter):
                self._chapter = chapter
                self._updated = True
            if(not self._bookType):
                # check type (bookType)
                start = content.find('分类：')
                self._bookType = content[start+3:]
                start = self._bookType.find('>')
                self._bookType = self._bookType[start+1:]
                end = self._bookType.find('</a>')
                self._bookType = self._bookType[:end]
                self._updated = True

                self._bookType = self._bookType[:end]
            return self._updated
        def DownloadBook(self,path,out=print):
            # fill back the info by the website
            res = urllib.request.urlopen(self._website,timeout=60)
            content = res.read()
            # decode the content
            if (res.info().get('Content-Encoding') == 'gzip'):
                gzipFile = gzip.GzipFile('','rb',9,io.BytesIO(content))
                content = gzipFile.read()
            content = content.decode("utf-8")
            res.close()

            # get chapter set
            start = content.find("yulan")
            chapters = content[start:]
            end = chapters.find("</div>")
            chapters = chapters[:end]
            self._chapterSet = chapters.split("href=")
            out(self._name)
            # download chapters one by one
            for chapter in self._chapterSet:
                if("<B>" in chapter):
                    self._text += chapter[chapter.find("<B>")+3:chapter.find("<a")]+"\n"
                elif("</B>" not in chapter)and("https" in chapter):
                    # go to the website and download
                    chapter = chapter[chapter.find('"')+1:]
                    chapter = chapter[:chapter.find('"')]
                    self._DownloadChapter(chapter)
                    # log for progress
                    out("\r"+self._chapter,end=" "*20)
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
            try:
                # open chapter url
                chRes = urllib.request.urlopen(url,timeout=60)
                content = chRes.read()
            except:
                time.sleep(10)
                self._DownloadChapter(url)
                return
            # decode the content
            if (chRes.info().get('Content-Encoding') == 'gzip'):
                gzipFile = gzip.GzipFile('','rb',9,io.BytesIO(content))
                content = gzipFile.read()
            content = content.decode("utf-8")
            chRes.close()
            # get the title
            self._text += "\n"
            start = content.find("<h1>")
            title = content[start+4:]
            end = title.find("</h1>")
            title = title[:end]
            self._chapter = title
            self._text += title + "\n"
            # get the content
            start = content.find('id="content">')
            c = content[start+13:]
            end = c.find("<div")
            c = c[:end]
            c = c.replace("&nbsp;"," ")
            c = c.replace("<br />    ","\n")
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
        for row in self._cursor.execute("select * from books where end='true' and download='false' and website like '%80txt%'"):
            self.books.append(self.Book(row[4],name=row[0],writer=row[1],date=row[2],chapter=row[3],bookType=row[5]))
        out("downloading")
        for book in self.books:
            if(book.DownloadBook(self._path)):
                self._cursor.execute("update books set download='true' where website='"+book._website+"'")
        out("finish download")
    def Update(self,out=print):
        # get all books from db to boo
        self.books.clear()
        for row in self._cursor.execute("select * from books where website like '%80txt%'"):
            self.books.append(self.Book(row[4],name=row[0],writer=row[1],date=row[2],chapter=row[3],bookType=row[5]))
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
                    out("\rupdate "+book._name,end="")
                    os.rename(self._path+"\\"+book._bookType+"\\"+book._name+"-"+book._writer+".txt",self._path+"\\"+book._bookType+"\\-"+book._name+"-"+book._writer+".txt")
                    book._name = '-'+book._name
                    self._cursor.execute("update books set name='"+book._name+"' where website='"+book._website+"'")
                    self._conn.commit()
        out("update finish")
    def Explore(self,n,out=print):
        # get the max book num from the db
        self.books.clear()
        self._bookNum = 0
        for row in self._cursor.execute("select website from books where website like '%80txt%' order by website desc"):
            i = row[0]
            i = int(i[i.find("_")+1:i.rfind(".")])
            if(i > self._bookNum):
                self._bookNum = i
        self._bookNum += 1
        errorPage = 0
        # check any new book by the book num and try to save it
        while(errorPage<n):
            b = self.Book("https://www.80txt.com/txtml_"+str(self._bookNum)+".html")
            if(b._name):
                out("\r"+b._name,end="")
                self.books.append(b)
                flag = bool(self._cursor.execute("select * from books where name='"+b._name+"' and website='"+b._website+"'").fetchone())
                if(not(flag)):
                    sql = (
                        "insert into books (name,writer,date,chapter,website,type,download) values"
                        "('"+b._name+"','"+b._writer+"','"+b._date+"','"+b._chapter+"','"+b._website+"','"+b._bookType+"','false')"
                    )
                    self._cursor.execute(sql)
                    self._conn.commit()
                errorPage = 0
            else:
                out("\rerror "+str(self._bookNum),end="")
                f = self._cursor.execute("select * from error where website='"+b._website+"'").fetchone()
                if(not f):
                    self._cursor.execute("insert into error (website) values ('"+b._website+"')")
                    self._conn.commit()
                errorPage += 1
            self._bookNum += 1
        self.ErrorUpdate()
    def ErrorUpdate(self,out=print,checkAll=False):
        # get all books from db to boo
        self.books.clear()
        out("error update")
        condition = ""
        if(not checkAll): condition = " and type is null or type=''"
        for row in self._cursor.execute("select * from error where website like '%80txt%'"+condition):
            web = row[0]
            errType = ''
            try:
                res = urllib.request.urlopen(web, timeout=60)
            except (urllib.error.HTTPError): errType = '404'
            except (urllib.request.socket.timeout): errType = 'timeout'
            self._cursor.execute("update error set type='"+errType+"' where website='"+web+"'")
            if(errType == ''):
                self._conn.execute("delete from error where website='"+web+"'")
                b = self.Book(web)
                out("\rerror update "+b._name,end="")
                sql = (
                    "insert into books (name,writer,date,chapter,website,type,download) values"
                    "('"+b._name+"','"+b._writer+"','"+b._date+"','"+b._chapter+"','"+b._website+"','"+b._bookType+"','false')"
                )
                self._cursor.execute(sql)
            self._conn.commit()
        out("\nerror update finish")