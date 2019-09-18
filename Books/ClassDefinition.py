import urllib.request, gzip
import re
import io, os, time, copy
import threading

MAX_THREAD = 10

### new functional class
class Book():
    def __init__(self,**kwargs):
        self.base_web = kwargs["base_web"]
        self.__download_web = kwargs["download_web"]
        self.__chapter_web = kwargs["chapter_web"]
        self.__decode = kwargs["decode"] if("decode" in kwargs) else "utf8"
        self.__timeout = kwargs["timeout"] if("timeout" in kwargs) else 30
    def new(self,**kwargs):
        book = copy.copy(self)
        book.book_num = str(kwargs["book_num"])
        book.base_web = self.base_web.format(book.book_num)
        book.__download_web = self.__download_web.format(book.book_num)
        book.name = kwargs["name"] if("name" in kwargs) else ""
        book.writer = kwargs["writer"] if("writer" in kwargs) else ""
        book.date = kwargs["date"] if("date" in kwargs) else ""
        book.last_chapter = kwargs["last_chapter"] if("last_chapter" in kwargs) else ""
        book.book_type = kwargs["book_type"] if("book_type" in kwargs) else ""
        book.updated = True
        try:book.__get_basic_info()
        except:pass
        return book
    def __str__(self):
        return "Name:\t"+self.name+"\nWriter:\t"+self.writer+"\nType:\t"+self.book_type
    ### custom define function
    def _cut_name(self,c):
        raise NotImplementedError()
    def _cut_writer(self,c):
        raise NotImplementedError()
    def _cut_date(self,c):
        raise NotImplementedError()
    def _cut_last_chapter(self,c):
        raise NotImplementedError()
    def _cut_book_type(self,c):
        raise NotImplementedError()
    def _cut_chapter(self,c):
        raise NotImplementedError()
    def _cut_title(self,c):
        raise NotImplementedError()
    def _cut_chapter_title(self,c):
        raise NotImplementedError()
    def _cut_chapter_content(self,c):
        raise NotImplementedError()
    def open_website(self,url):
        res = urllib.request.urlopen(url,timeout=self.__timeout)
        content = res.read()
        if (res.info().get('Content-Encoding') == 'gzip'):
            content = gzip.GzipFile('','rb',9,io.BytesIO(content)).read()
        content = content.decode(self.__decode,"ignore")
        res.close()
        return content
    ### basic function
    def __get_basic_info(self):
        # read website
        try:
            content = self.open_website(self.base_web)
        except (urllib.error.HTTPError, urllib.error.URLError, urllib.request.socket.timeout):
            raise RuntimeError("Unable to open the website for basic information")
        # check the website book information
        if(not self.name):
            # get name
            self.name = self._cut_name(content)
            self.updated = False
        if(not self.writer):
            # get writer (writer)
            self.writer = self._cut_writer(content)
            self.updated = False
        # get date (always get)
        date = self._cut_date(content)
        if(self.date != date):
            self.date = date
            self.updated = False
        # get chapter (always get)
        last_chapter = self._cut_last_chapter(content)
        if(self.last_chapter != last_chapter):
            self.last_chapter = last_chapter
            self.updated = False
        if(not self.book_type):
            # check type (bookType)
            self.book_type = self._cut_book_type(content)
            self.updated = False
        return self.updated
    def download(self,path,out):
        # open chapters page
        try:
            content = self.open_website(self.__download_web)
        except (urllib.error.HTTPError, urllib.error.URLError, urllib.request.socket.timeout):
            raise RuntimeError("Unable to open the website for chapter lists")
        # read all chapters url
        chapters = self._cut_chapter(content)
        titles = self._cut_title(content)
        text = ""
        '''
        # read actual content
        for i in range(min(len(titles),len(chapters))):
            text += self.__download_chapter(titles[i],chapters[i],out)
            out(titles[i])
        '''
        # for thread 
        arr = []
        lock = threading.Lock()
        threads = []
        # read all chapter
        for i in range(min(len(titles),len(chapters))):
            t = threading.Thread(target=self.__download_chapter_thread,args=(titles[i],chapters[i],out,arr,lock))
            t.daemon = True
            t.start()
            threads.append(t)
        for t in threads:
            t.join()
        # order the chapter by its url
        text = ""
        for i in range(min(len(titles),len(chapters))):
            for j in range(len(arr)):
                if(chapters[i] == arr[j][0]):
                    print(arr[j][1][20:40])
                    text += arr[j][1]
                    break
        # save actual content
        try: os.mkdir(path)
        except: pass
        f = open(path+'/'+self.book_num+".txt","w",encoding="utf8")
        f.write(text)
        f.close()
        return True
    def __download_chapter(self,chapter_title,chapter_url,out):
        # check valid url pattern or not
        m = re.match(self.__chapter_web,chapter_url)
        if(not m): return '\n'+'-'*20+'\n'+chapter_title+'\n'+'-'*20+'\n'
        # open chapter url
        while(True):
            try:
                content = self.open_website(chapter_url)
                break
            except :
                # wait for a while and try again
                time.sleep(self.__timeout//10)
                out("Reload")
        # read title
        t =  self._cut_chapter_title(content)
        # read content
        c = self._cut_chapter_content(content)
        return  '\n'+'-'*20+'\n'+chapter_title+'\n'+'-'*20+'\n'+t+c+'\n'
    def __download_chapter_thread(self,chapter_title,chapter_url,out,arr,lock):
        # check url pattern
        m = re.match(self.__chapter_web,chapter_url)
        if(not m):
            lock.acquire()
            arr.append((chapter_url,'\n'+'-'*20+'\n'+chapter_title+'\n'+'-'*20+'\n'))
            lock.release()
            return
        # open chapter url
        while(True):
            try:
                content = self.open_website(chapter_url)
                break
            except :
                # wait for a while and try again
                time.sleep(self.__timeout//10)
                out("Reload")
        # read title
        t = self._cut_chapter_title(content)
        # read content
        c = self._cut_chapter_content(content)
        lock.acquire()
        # put the result into common array
        arr.append((chapter_url,'\n'+'-'*20+'\n'+chapter_title+'\n'+'-'*20+'\n'+t+c+'\n'))
        lock.release()
        return

class BookSite():
    def __init__(self,**kwargs):
        self.__Book = kwargs["book"](**kwargs["web"],**kwargs["setting"])
        self.conn = kwargs["conn"]
        self.path = kwargs["path"]
        self.identify = kwargs["identify"]
        self.running_thread = 0
        self.error_page = 0
    def __str__(self):
        return type(self.__Book).__name__+' Site'
    def download(self,out):
        out(self.identify+"="*15)
        for result in self.conn.execute("select name,writer,date,chapter,num,type from books where site = '"+self.identify+"' and end='true' and download='false'").fetchall():
            info = {
                "name":result[0],
                "writer":result[1],
                "date":result[2],
                "last_chapter":result[3],
                "book_num":result[4],
                "book_type":result[5]
            }
            b = self.__Book.new(**info)
            if(b.updated):
                b.download(self.path,out)
                self.conn.execute("update books set download='true' where site='"+self.identify+"' and num="+b.book_num)
                self.conn.commit()
    def download_thread(self,info,lock,out):
        b = self.__Book.new(**info)
        if(b.updated):
            b.download(self.path)
            if(lock): lock.acquire()
            try:
                self.conn.execute("update books set download='true' where site='"+self.identify+"' and num="+b.book_num)
                self.conn.commit()
            finally:
                if(lock): lock.release()
    def book(self,bookId):
        return self.__Book.new(book_num=str(bookId))
    def query(self,**kwargs):
        query = "(site='"+self.identify+"')"
        if("name" in kwargs):
            query += " and (name='"+kwargs["name"]+"')"
        if("writer" in kwargs):
            query += " and (writer='"+kwargs["writer"]+"')"
        if("book_type" in kwargs):
            query += " and (type='"+kwargs["type"]+"')"
        sql = "select num,name,writer,type from books where "+query
        result = self.conn.cursor().execute(sql).fetchall()
        for i in range(len(result)):
            result[i] = {
                "num":result[i][0],
                "name":result[i][1],
                "writer":result[i][2],
                "type":result[i][3]
            }
        return result
    def update(self,out):
        out(self.identify+"="*15)
        cursor = self.conn.cursor()
        self.running_thread = 0
        lock = threading.Lock()
        threads = []
        for result in cursor.execute("select name,writer,date,chapter,num,type from books where site='"+self.identify+"' and (read='false' or read is null) order by date desc"):
            while(self.running_thread >= MAX_THREAD):pass
            info = {
                "name":result[0],
                "writer":result[1],
                "date":result[2],
                "last_chapter":result[3],
                "book_num":result[4],
                "book_type":result[5]
            }
            th = threading.Thread(target=self.update_thread,args=(info,lock,out))
            th.daemon = True
            self.running_thread += 1
            th.start()
            threads.append(th)
        for thread in threads:thread.join()
    def update_thread(self,info,lock,out):
        b = self.__Book.new(**info)
        lock.acquire()
        try:
            if(not b.updated):
                out(b.name)
                cursor = self.conn.cursor()
                sql = "update books set date='"+b.date+"', chapter='"+b.last_chapter+"' where site='"+self.identify+"' and num="+b.book_num
                cursor.execute(sql)
                self.conn.commit()
            else:out("skip\t"+b.name)
        finally:
            lock.release()
            self.running_thread -= 1
    def explore(self,n,out):
        out(self.identify+"="*15)
        book_num = 0
        # get largest book num
        tran = self.conn.cursor().execute("select num from books where site ='"+self.identify+"' order by num desc").fetchone()
        if(tran): book_num = tran[0]
        book_num += 1
        # init data for the explore thread
        self.error_page = 0
        self.running_thread = 0
        lock = threading.Lock()
        threads = []
        # explore in threads
        while(self.error_page<n):
            if(self.running_thread < MAX_THREAD):
                th = threading.Thread(target=self.explore_thread,args=(book_num,lock,out))
                th.daemon = True
                book_num += 1
                self.running_thread += 1
                th.start()
                threads.append(th)
        # confirm all thread are finish
        for thread in threads:thread.join()
    def explore_thread(self,num,lock,out):
        # get book info
        b = self.__Book.new(book_num=num)
        # try to update it
        lock.acquire()
        try:
            cursor = self.conn.cursor()
            if(not b.updated):
                out(b.name)
                flag = bool(cursor.execute("select * from books where name='"+b.name+"' and site='"+self.identify+"' and num="+b.book_num).fetchone())
                if(not(flag)):
                    # add record to books table
                    sql = (
                        "insert into books (site,num,name,writer,date,chapter,type,download) values"
                        "('"+self.identify+"',"+b.book_num+",'"+b.name+"','"+b.writer+"','"+b.date+"','"+b.last_chapter+"','"+b.book_type+"','false')"
                    )
                    # delete record from error
                    cursor.execute(sql)
                    sql = (
                        "delete from error where site='"+self.identify+"' and num="+b.book_num
                    )
                    cursor.execute(sql)
                    self.conn.commit()
                self.error_page = 0
            else:
                # update if the record does not exist in database
                out("error "+str(num)+"\t"+b.name)
                f_books = bool(cursor.execute("select * from books where site='"+self.identify+"' and num="+b.book_num).fetchone())
                f_error = bool(cursor.execute("select * from error where site='"+self.identify+"' and num="+b.book_num).fetchone())
                if((not f_books)and(not f_error)):
                    cursor.execute("insert into error (site,num) values ('"+self.identify+"','"+b.book_num+"')")
                    self.conn.commit()
                self.error_page += 1
        finally:
            lock.release()
            self.running_thread -= 1        
    def error_update(self,out):
        out(self.identify+"="*15)
        cursor = self.conn.cursor()
        self.running_thread = 0
        lock = threading.Lock()
        threads = []
        for result in cursor.execute("select site,num from error where site='"+self.identify+"'"):
            while(self.running_thread >= MAX_THREAD):pass
            num = result[1]
            th = threading.Thread(target=self.explore_thread,args=(num,lock,out))
            th.daemon = True
            self.running_thread += 1
            th.start()
            threads.append(th)
        for thread in threads:thread.join()
        pass
    def __del__(self):
        self.conn.close()
