import urllib.request
import os
import zipfile
import sqlite3
import http


def download(conn, book):
    print(book.name)
    # make a url for download books
    cut = book.website.rfind("/") + 1
    urlA = book.website[:cut]
    urlB = urllib.parse.quote(book.name)
    url  = urlA + urlB + ".zip"
    # get the zip file
    res = urllib.request.urlopen(url)
    try:
        content = res.read()
        path = os.getcwd() + "\\download\\" + book.name + "-" + book.writer + ".zip"
        print("get zip file", end="\t")

        # save the zip file
        f = open(path, 'wb')
        f.write(content)
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
        return

def bookEnd(book):
    # get book id
    
    # go to the website
    
    # extract the data from the website:
    # last update date
    
    # chapter

    # state
    print("developing")

def anyNew(conn):
    errorPage = 0
    # get largest 80txt book id from database

    # loop to check any new books suit the type after the id
    '''
        # if have, update the book into database

        # if it is decode error, record the book page to error table

        # if the book page is not exist, add error by 1

        except:
            errorPage += 1
        '''
    print("developing")
