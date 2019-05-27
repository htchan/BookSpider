import urllib.request
import os
import zipfile
import sqlite3
import http
import gzip
import io
import ClassDefinition


def download(conn, book):
    print(book.name+"\t"+book.date[:4])
    # make a url for download books
    cut = book.website.rfind("/") + 1
    urlA = book.website[:cut]
    urlB = urllib.parse.quote(book.name)
    url  = urlA + urlB + ".zip"
    # get the zip file
    try:
        res = urllib.request.urlopen(url,timeout=60)

        contentLen = int(res.getheader("content-length"))
        if(contentLen):
            blockSize = int(contentLen)//100
            blockSize = max(4096, blockSize)
        else:
            blockSize = 4096
        
        buf = io.BytesIO()
        size = 0
        print("0%",end="")
        while (True):
            content = res.read(blockSize)
            if(not(content)):
                break
            buf.write(content)
            size+=len(content)
            print("\r"+str(size*100//contentLen)+'%',end="")
        print(end="\t")
        if(size < contentLen):
            print("Incomplete read")
            return
        path = os.getcwd() + "\\download\\" + book.name + "-" + book.writer + ".zip"

        # save the zip file
        f = open(path, 'wb')
        f.write(buf.getvalue())
        f.close()

        # unzip the file
        zFile = zipfile.ZipFile(path,'r')
        zFile.extractall()
        zFile.close()
        print("unzip", end="\t")

        #rename extracted and put it into the download file
        os.rename("all.txt", path[:-4]+".txt")

        # delete the zip file
        os.remove(path)
        print("get txt file", end="\t")

        # update database
        c = conn.cursor()
        c.execute("update books set download = 'true' where website='"+book.website+"'")
        conn.commit()
        print("database update")
    except http.client.IncompleteRead:
        print("Incomplete read")
        return
    except urllib.error.URLError as e:
        print(str(e)+" --- skip")
        return
    except urllib.request.socket.timeout:
        print("\ttime out, skip")
        return

def bookUpdate(conn, book):
    try:
        # get book id
        bookID = book.website[:book.website.rfind("/")]
        bookID = bookID[bookID.rfind("/")+1:]
        url = "https://www.80txt.com/txtxz/"+bookID+".html"
        print(url)
        # go to the website
        res = urllib.request.urlopen(url,timeout=60)
        content = res.read()
        if res.info().get('Content-Encoding') == 'gzip':
            gzipFile = gzip.GzipFile('','rb',9,io.BytesIO(content))
            content=gzipFile.read()
        content = content.decode("utf-8")

        # extract the data from the website:
        # last update date
        update = content[content.find("更新时间：</b>"):]
        update = update[update.find('</b>')+4:]
        update = update[:update.find('</li>')]
        print(update, end="\t")
        
        if(update!=book.date):
            c = conn.cursor()
            book.date = update

            # chapter
            chapter = content[content.find("最新章节："):]
            chapter = chapter[chapter.find('</b>')+4:]
            chapter = chapter[:chapter.find('</li>')]
            book.chapter = chapter
            print(chapter, end="\t")
            
            # state
            state = content[content.find("写作进度："):]
            state = state[:state.find("</li>")]
            if("已完成" in state):
                c.execute("update books set end='true' where website='"+book.website+"'")
                

            #save back to database
            sql = ("update books set date='"+book.date+"', chapter='"+book.chapter+"' where website='"+book.website+"'")
            c.execute(sql)
            conn.commit()
        print("\n")
    except Exception as e:
        print(e)

def anyNew(conn):
    errorPage = 0
    # get largest 80txt book id from database
    bookId=0
    c = conn.cursor()
    for row in c.execute("select website from books where website like '%80txt%'"):
        tranId = row[0]
        tranId = tranId[:tranId.rfind("/")]
        tranId = tranId[tranId.rfind("/")+1:]
        tranId = int(tranId)
        if(tranId>bookId):
            bookId = tranId
    # loop to check any new books suit the type after the id
    bookId += 1
    while(errorPage<10):
        try:
            # try to get the page
            url = 'https://www.80txt.com/txtxz/'+str(bookId)+'/down.html'
            print(url, end="\t")
            res = urllib.request.urlopen(url)
            content = res.read()
            errorPage = 0

            # decode the content
            if (res.info().get('Content-Encoding') == 'gzip'):
                gzipFile = gzip.GzipFile('','rb',9,io.BytesIO(content))
                content = gzipFile.read()
            content = content.decode("utf-8")

            # check type (bookType)
            start = content.find('分类：')
            bookType = content[start+3:]
            end = bookType.find('</span>')
            bookType = bookType[:end]
            allowTypes = {"奇幻修真","奇幻魔法","异术超能","东方传奇","江湖武侠","未来幻想"}
            if(not(bookType in allowTypes)):
                print("wrong type:"+str(bookType))
                print()
                bookId += 1
                continue
            print(bookType)
            
            # get writer (writer)
            start = content.find("作者：")
            writer = content[start+3:]
            end = writer.find("</a>")
            writer = writer[:end]
            print(writer, end="\t")
            
            # get download link (link)
            start = content.find('https://dz.80txt.com')
            link = content[start:]
            end = link.find('.zip')
            link = link[:end+4]
            print(link)
            print()
            missing = 0

            # get name
            name = link[link.rfind("/")+1:]
            name = name[:name.rfind(".")]

            #save to database
            sql = ('insert into books (name, writer, website, type, end, download, read)'
            ' values '
            "('"+name+"', '"+writer+"', '"+link+"', '"+bookType+"', 'false', 'false', 'false')"
            )
            c.execute(sql)
            conn.commit()

        # if it is decode error, record the book page to error table
        except UnicodeDecodeError:
            c.execute("insert into error (website) values ('"+link+"')")
            conn.commit()

        # if the book page is not exist, add error by 1
        except urllib.error.HTTPError:
            print("error")
            errorPage += 1
            print(str(errorPage)+"/100")
        bookId += 1
    # have to update the date and last chapter also
    for row in c.execute("select name,writer,website,date from books where date is null"):
        book = ClassDefinition.Book(name=row[0], writer=row[1], website=row[2])
        bookUpdate(conn, book)

class TXT80():
    def __init__(self, dbConn, path):
        self._webpage = ""
        self._bookNum = 0
        self.books = []
        self._conn = dbConn
        self._cursor = self._conn.cursor()
        self._path = path
    class Book(ClassDefinition.BaseBook):
        def _getBasicInfo(self):
            if(not self._name):
                # fill back the info by the website
                res = urllib.request.urlopen(self._website)
                content = res.read()

                # decode the content
                if (res.info().get('Content-Encoding') == 'gzip'):
                    gzipFile = gzip.GzipFile('','rb',9,io.BytesIO(content))
                    content = gzipFile.read()
                content = content.decode("utf-8")

                # get name
                start = content.find("titlename")
                self._name = content[start+9]
                start = self._name.find("<h1>")
                self._name = self._name[start+4]
                end = self._name.find("</h1>")
                self._name = self._name[:end]

                # get writer (writer)
                start = content.find("作者：")
                self._writer = content[start+3:]
                start = self._writer.find('>')
                self._writer = self._writer[start+1:]
                end = self._writer.find('</a>')
                self._writer = self._writer[:end]

                # get date
                start = content.find('更新时间：')
                self._date = content[start+5:]
                end = self._date.find('</span>')
                self._date = self._date[:end]

                # get chapter
                start = content.rfind("<li>")
                self._chapter = content[start+4:]
                start = self._chapter.find('">')
                self.chapter = self._chapter[start+2:]
                end = self._chapter.find('</a>')
                self._chapter = self._chapter[:end]

                # check type (bookType)
                start = content.find('分类：')
                self._bookType = content[start+3:]
                start = self.bookType.find('>')
                self._bookType = self._bookType[start+1:]
                end = self._bookType.find('</a>')

                self._bookType = self._bookType[:end]
        def DownloadBook(self,path):
            # fill back the info by the website
            res = urllib.request.urlopen(self._website)
            content = res.read()
            # decode the content
            if (res.info().get('Content-Encoding') == 'gzip'):
                gzipFile = gzip.GzipFile('','rb',9,io.BytesIO(content))
                content = gzipFile.read()
            content = content.decode("utf-8")

            # get chapter set
            start = content.find("yulan")
            chapters = content[start:]
            end = chapters.find("</div>")
            chpapters = chapters[:end]
            self._chapterSet = chapters.split("href=")

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
                    print("\r"+self._chapter,end=" "*20)
            # save it into file
            try:
                os.mkdir(path)
            except: pass
            f = open(path+"\\"+self._name+"-"+self._writer+".txt","w",encoding='utf8')
            f.write(self._text)
            # TODO: update the sql record
        def _DownloadChapter(self,url):
            # open chapter url
            chRes = urllib.request.urlopen(url)
            content = chRes.read()
            # decode the content
            if (chRes.info().get('Content-Encoding') == 'gzip'):
                gzipFile = gzip.GzipFile('','rb',9,io.BytesIO(content))
                content = gzipFile.read()
            content = content.decode("utf-8")
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
    def Download(self):
        # put [end] book in db to books
        for row in self._cursor.execute("select * from books where end='true' and download='false'"):
            self.books.append(self.Book(row[4],name=row[0],writer=row[1],date=row[2],chapter=row[3],bookType=row[5]))
        for book in self.books:
            book.DownloadBook(self._path)
    def Update():#TODO
        # get all books from db to books
        # check any update
        # if it update, but ended, add a '-' at first of the file name
        pass
    def Explore():#TODO
        # get the max book num from the db
        # check any new book by the book num
        pass