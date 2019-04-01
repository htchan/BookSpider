import urllib.request
import os
import zipfile
import sqlite3
import http
import gzip
import io
import ClassDefinition


def download(conn, book):
    print(book.name)
    # make a url for download books
    cut = book.website.rfind("/") + 1
    urlA = book.website[:cut]
    urlB = urllib.parse.quote(book.name)
    url  = urlA + urlB + ".zip"
    # get the zip file
    try:
        res = urllib.request.urlopen(url)

        contentLen = int(res.getheader("content-length"))
        if(contentLen):
            blockSize = int(contentLen)//100
            blockSize = max(4096, blockSize)
        else:
            blockSize = 4096
        
        buf = io.BytesIO()
        size = 0
        while (True):
            content = res.read(blockSize)
            if(not(content)):
                break
            buf.write(content)
            size+=len(content)
            print("\r"+str(size*100//contentLen)+'%',end="")
        print()

        path = os.getcwd() + "\\download\\" + book.name + "-" + book.writer + ".zip"
        print("get zip file", end="\t")

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
        print(e.reason().strerror+" --- skip")
        return

def bookUpdate(conn, book):
    # get book id
    bookID = book.website[:book.website.rfind("/")]
    bookID = bookID[bookID.rfind("/")+1:]
    url = "https://www.80txt.com/txtxz/"+bookID+".html"
    print(url)
    # go to the website
    res = urllib.request.urlopen(url)
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
    while(errorPage<100):
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
            "('"+name+"', '"+writer+"', '"+link+"', '"+bookType+"', 'false, 'false', 'false')"
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
    print("developing")
    # have to update the date and last chapter also
    book = ClassDefinition.Book(name=name, writer=writer, website=link)
    bookUpdate(conn, book)
