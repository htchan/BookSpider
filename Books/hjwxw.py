import urllib.request
import os
import zipfile
import sqlite3
import http
import gzip
import io
try: import ClassDefinition
except: import Books.ClassDefinition as ClassDefinition
import time

class HJWXW():
    def __init__(self, dbConn, path):
        self._webpage = ""
        self._bookNum = 0
        self.books = []
        self._conn = dbConn
        self._cursor = self._conn.cursor()
        self._path = path
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

                if("未找到該頁,1秒后為您跳轉" in content):
                    return False

            # return false if the webpage is not exist or not avaliable
            except (urllib.error.HTTPError, urllib.request.socket.timeout):
                return False

            if(not self._name):
                # get name
                start = content.find("<title>")
                self._name = content[start+7:]
                end = self._name.find("/")
                self._name = self._name[:end]
                self._updated = True
            if(not self._writer):
                # get writer (writer)
                start = content.rfind("作者標簽:")
                self._writer = content[start+5:]
                end = self._writer.find('">')
                self._writer = self._writer[:end].strip()
                self._updated = True

            # get date (always get)
            start = content.find('更新時間: ')
            date = content[start+6:]
            end = date.find('">')
            date = date[:end]
            if(self._date != date):
                self._date = date
                self._updated = True

            # get chapter (always get)
            start = content.find("章節名:")
            chapter = content[start+4:]
            end = chapter.find('更新時間')
            chapter = chapter[:end]
            if(self._chapter != chapter):
                self._chapter = chapter
                self._updated = True
            if(not self._bookType):
                # check type (bookType)
                bookType = ""
                c = content
                start = c.find('小說分類標簽:')
                while(start>0):
                    if(bookType != ""):
                        bookType += ","
                    c = c[start+8:]
                    end = c.find('  ')
                    bookType += c[:end].strip()
                    c = c[end:]
                    start = c.find('小說分類標簽:')                    
                self._updated = True
                self._bookType = bookType
            return self._updated
        def DownloadBook(self,path,out=print):
            # fill back the info by the website
            res = urllib.request.urlopen(self._website.replace("Book","Book/Chapter"),timeout=60)
            content = res.read()
            # decode the content
            if (res.info().get('Content-Encoding') == 'gzip'):
                gzipFile = gzip.GzipFile('','rb',9,io.BytesIO(content))
                content = gzipFile.read()
            content = content.decode("utf-8")
            res.close()

            # get chapter set
            start = content.find("tbchapterlist")
            chapters = content[start+6:]
            end = chapters.find("</table>")
            chapters = chapters[:end]
            self._chapterSet = chapters.split("</a>")
            out(self._name)
            # download chapters one by one
            for chapter in self._chapterSet:
                if("href" in chapter):
                    # go to the website and download
                    chapter = chapter[chapter.find('href="')+6:]
                    chapter = chapter[:chapter.find('"')]
                    self._DownloadChapter("https://tw.hjwzw.com"+chapter)
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
                time.sleep(5)
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
            title = title[:end].strip()
            self._chapter = title
            self._text += title + "\n\n"
            # get the content
            start = content.find('<p/>')
            con = content[start+4:]
            end = con.find("<p />")
            con = con[:end]
            con = con.split("<p/>")
            for c in con:
                if(c!=""): self._text += c.strip() + "\n\n"
    def Download(self,out=print):
        # put [end] book in db to books
        self.books.clear()
        for row in self._cursor.execute("select * from books where end='true' and download='false' and website like '%hjwzw%'"):
            self.books.append(self.Book(row[4],name=row[0],writer=row[1],date=row[2],chapter=row[3],bookType=row[5]))
        out("downloading")
        for book in self.books:
            if(book.DownloadBook(self._path)):
                self._cursor.execute("update books set download='true' where website='"+book._website+"'")
        out("finish download")
    def Update(self,out=print):
        # get all books from db to boo
        self.books.clear()
        for row in self._cursor.execute("select * from books where website like '%hjwzw%'"):
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
        for row in self._cursor.execute("select website from books where website like '%hjwzw%' order by website desc"):
            i = row[0]
            i = int(i[i.rfind("/")+1:])
            if(i > self._bookNum):
                self._bookNum = i
        self._bookNum += 1
        errorPage = 0
        # check any new book by the book num
        while(errorPage<n):
            b = self.Book("https://tw.hjwzw.com/Book/"+str(self._bookNum))
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
                    self._cursor.execute("delete from error where website='"+b._website+"'")
                    self._conn.commit()
                errorPage = 0
            else:
                out("\rerror "+str(self._bookNum),end="")
                time.sleep(5)
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
        condition = ""
        if(not checkAll): condition = " and type is null or type=''"
        rows = self._cursor.execute("select * from error where website like '%hjwzw%'"+condition).fetchall()
        for row in rows:
            web = row[0]
            errType = ''
            try:
                res = urllib.request.urlopen(web, timeout=60)
                content = res.read()

                # decode the content
                if (res.info().get('Content-Encoding') == 'gzip'):
                    gzipFile = gzip.GzipFile('','rb',9,io.BytesIO(content))
                    content = gzipFile.read()
                content = content.decode("utf-8")
                res.close()
                if("未找到該頁,1秒后為您跳轉" in content): errType = '404'
            except (urllib.error.HTTPError): errType = '404'
            except (urllib.request.socket.timeout): errType = 'timeout'
            self._cursor.execute("update error set type='"+errType+"' where website='"+web+"'")
            if(errType == ''):
                b = self.Book(web)
                if(b._name):
                    self._conn.execute("delete from error where website='"+web+"'")
                    out("\rupdate "+b._name,end="")
                    sql = (
                        "insert into books (name,writer,date,chapter,website,type,download) values"
                        "('"+b._name+"','"+b._writer+"','"+b._date+"','"+b._chapter+"','"+b._website+"','"+b._bookType+"','false')"
                    )
                    self._cursor.execute(sql)
            self._conn.commit()
        out("error update finish")