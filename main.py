import Books.downloadBooks
import os
import io

bookDownload = None
musicDownload = None
storyDownload = None

def book():
    global bookDownload
    bookDownload = Books.downloadBooks.BookCollection(os.getcwd())
    looping = True
    while(looping):
        os.system("cls")
        print("Book Download"+'='*20)
        print("D:\tDownload books")
        print("C:\tCheck new books")
        print("U:\tUpdate books")
        print("A:\tAll books update and check")
        print("E:\tExit")
        res = input(">>>").upper().strip()[0]
        if(res == "D"): bookDownload.Download()
        elif(res == "C"): bookDownload.Explore(100)
        elif(res == "U"): bookDownload.Update()
        elif(res == "A"): bookDownload.ErrorUpdate()
        elif(res == "E"): looping = False
        else:
            os.system("cls")
            print("wrong input")
    bookDownload.close()

def music():
    pass

def baidu():
    pass

### GUI ###
os.popen("chcp 936")
looping = True
while looping:
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
    if(key == 0): looping = False
    elif(key == 1): book()
    elif(key == 2): music()
    elif(key == 3): baidu()
    else:
        os.system("cls")
        print("wrong input")