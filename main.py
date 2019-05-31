import Books.downloadBooks
import os
import io

bookDownload = None
musicDownload = None
storyDownload = None

def book():
    global bookDownload
    bookDownload = Books.downloadBooks.BookCollection(os.getcwd)
    pass

def music():
    pass

def baidu():
    pass

### GUI ###
while True:
    print("download"+"-"*20)
    print("1: books")
    print("2: music")
    print("3: baidu")
    key = input("--> ")
    try:
        key = int(key)
    except:
        os.system("cls")
        print("wrong input")
    if(key==1):
        book()
    if(key==2):
        music()
    if(key==3):
        baidu()
    else
        os.system("cls")
        print("wrong input")