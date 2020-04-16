import urllib.request, gzip
import re
import io, os, time, copy
import threading

MAX_THREAD = 300

### new functional class
class Book():
    def __init__(self, **kwargs):
        self.base_web = kwargs["base_web"]
        self.__download_web = kwargs["download_web"]
        self.__chapter_web = kwargs["chapter_web"]
        self.__decode = kwargs["decode"] if ("decode" in kwargs) else "utf8"
        self.__timeout = kwargs["timeout"] if ("timeout" in kwargs) else 30
    def new(self,**kwargs):
        book = copy.copy(self)
        book.book_num = str(kwargs["book_num"])
        book.base_web = self.base_web.format(book.book_num)
        book.__download_web = self.__download_web.format(book.book_num)
        book.name = kwargs["name"] if ("name" in kwargs) else ""
        book.writer = kwargs["writer"] if ("writer" in kwargs) else ""
        book.date = kwargs["date"] if ("date" in kwargs) else ""
        book.last_chapter = kwargs["last_chapter"] if ("last_chapter" in kwargs) else ""
        book.book_type = kwargs["book_type"] if ("book_type" in kwargs) else ""
        book.updated = True
        try:
            book.__get_basic_info()
        except:
            pass
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
        name = self._cut_name(content)
        if ((not self.name) or (self.name != name)):
            # get name
            self.name = name
            self.updated = False
        if (not self.writer):
            # get writer (writer)
            self.writer = self._cut_writer(content)
            self.updated = False
        # get date (always get)
        date = self._cut_date(content)
        if (self.date != date):
            self.date = date
            self.updated = False
        # get chapter (always get)
        last_chapter = self._cut_last_chapter(content)
        if (self.last_chapter != last_chapter):
            self.last_chapter = last_chapter
            self.updated = False
        if (not self.book_type):
            # check type (bookType)
            self.book_type = self._cut_book_type(content)
            self.updated = False
        return self.updated
    def download(self,path,out,i=0):
        out("download")
        # open chapters page
        try:
            content = self.open_website(self.__download_web)
        except (urllib.error.HTTPError, urllib.error.URLError, urllib.request.socket.timeout):
            if (i < 10):
                time.sleep(self.__timeout//10)
                out("reload chapter list",i)
                return self.download(path, out, i+1)
            elif (i == 10):
                out("unable to open the webite for chapter list")
                return False
                #raise RuntimeError("Unable to open the website for chapter lists")
        # read all chapters url
        chapters = self._cut_chapter(content)
        titles = self._cut_title(content)
        text = ""
        # single thread
        # read actual content

        for i in range(min(len(titles),len(chapters))):
            chapter_text = self.__download_chapter(titles[i],chapters[i],out)
            if (chapter_text == 'error'):
                return False
            else:
                text += chapter_text
            out(titles[i])
        '''
        # multi thread 
        arr = []
        lock = threading.Lock()
        threads = []
        self.threads_control = threading.Semaphore(MAX_THREAD)
        # read all chapter
        for i in range(min(len(titles),len(chapters))):
            self.threads_control.acquire()
            t = threading.Thread(target=self.__download_chapter_thread,args=(titles[i],chapters[i],out,arr,lock))
            t.daemon = True
            t.start()
            threads.append(t)
        for t in threads:
            t.join()
        # order the chapter by its url
        for i in range(len(arr)):
            if (arr[i][1] == ''):
                out("download error")
                return False
        text = ""
        for i in range(min(len(titles),len(chapters))):
            for j in range(len(arr)):
                if (chapters[i] == arr[j][0]):
                    text += arr[j][1]
                    break
        '''
        # save actual content
        try: os.mkdir(path)
        except: pass
        f = open(path+'/'+self.book_num+".txt","w",encoding="utf8")
        f.write(text)
        f.close()
        out("download success")
        return True
    def __download_chapter(self,chapter_title,chapter_url,out):
        # check valid url pattern or not
        m = re.match(self.__chapter_web,chapter_url)
        if (not m): return '\n'+'-'*20+'\n'+chapter_title+'\n'+'-'*20+'\n'
        # open chapter url
        while(True):
            try:
                content = self.open_website(chapter_url)
                break
            except :
                # wait for a while and try again
                time.sleep(self.__timeout//10)
                out("Reload chapter content")
        if (not content):
            raise RuntimeError("cannot download chapter")
        # read title
        t = self._cut_chapter_title(content)
        # read content
        c = self._cut_chapter_content(content)
        if (c == 'error'):
            return c
        return  '\n'+'-'*20+'\n'+chapter_title+'\n'+'-'*20+'\n'+t+c+'\n'
    def __download_chapter_thread(self,chapter_title,chapter_url,out,arr,lock):
        # check url pattern
        m = re.match(self.__chapter_web,chapter_url)
        if (not m):
            lock.acquire()
            arr.append((chapter_url,'\n'+'-'*20+'\n'+chapter_title+'\n'+'-'*20+'\n'))
            lock.release()
            self.threads_control.release()
            return
        # open chapter url
        for i in range(10):
            try:
                content = self.open_website(chapter_url)
                break
            except:
                # wait for a while and try again
                content = ""
                time.sleep(self.__timeout//10)
                out("Reload chapter_url "+chapter_url)
        if (not content):
            lock.acquire()
            arr.append((chapter_url, ""))
            lock.release()
            self.threads_control.release()
            return
        # read title
        t = self._cut_chapter_title(content)
        # read content
        c = self._cut_chapter_content(content)
        lock.acquire()
        if (c != 'error'):
            # put the result into common array
            arr.append((chapter_url,'\n'+'-'*20+'\n'+chapter_title+'\n'+'-'*20+'\n'+t+c+'\n'))
        else:
            arr.append((chapter_url, ''))
        lock.release()
        self.threads_control.release()
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
        for result in self.conn.execute("select name,writer,date,chapter,num,type,end, download from books where site = '"+self.identify+"' and end='true' and download='false' order by date").fetchall():
            info = {
                "name":result[0],
                "writer":result[1],
                "date":result[2],
                "last_chapter":result[3],
                "book_num":result[4],
                "book_type":result[5]
            }
            b = self.__Book.new(**info)
            out(self.identify+"\t"+b.book_num+"\t"+b.name+"\t"+"-"*15)
            if (b.updated):
                if (b.download(self.path+self.identify,out)):
                    out("Download Successfully")
                    self.conn.execute("update books set download='true' where site='"+self.identify+"' and num="+b.book_num)
                    self.conn.commit()
                else:
                    out("Download Error")
                    self.conn.execute("update books set download='error' where site='"+self.identify+"' and num="+b.book_num)
                    self.conn.commit()
            else:
                out("Not Updated")
                self.conn.execute("update books set end='false' where site='"+self.identify+"' and num="+b.book_num)
                self.conn.commit()
    def download_thread(self,info,lock,out):
        b = self.__Book.new(**info)
        if (b.updated):
            b.download(self.path)
            if (lock): lock.acquire()
            try:
                self.conn.execute("update books set download='true' where site='"+self.identify+"' and num="+b.book_num)
                self.conn.commit()
            finally:
                if (lock): lock.release()
    def book(self,bookId):
        return self.__Book.new(book_num=str(bookId))
    def query(self,**kwargs):
        query = "(site='"+self.identify+"')"
        if ("name" in kwargs):
            query += " and (name='"+kwargs["name"]+"')"
        if ("writer" in kwargs):
            query += " and (writer='"+kwargs["writer"]+"')"
        if ("book_type" in kwargs):
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
    def update(self,update_all,out):
        cursor = self.conn.cursor()
        self.threads_controller = threading.Semaphore(MAX_THREAD)
        lock = threading.Lock()
        threads = []
        sql = "select name,writer,date,chapter,num,type,download,version from books where site='"+self.identify+"'"
        if (not update_all):
            sql += " and (read='false' or read is null)"
        for result in cursor.execute(sql + " order by date desc"):
            info = {
                "name":result[0],
                "writer":result[1],
                "date":result[2],
                "last_chapter":result[3],
                "book_num":result[4],
                "book_type":result[5],
                "download":result[6],
                "version":result[7]
            }
            self.threads_controller.acquire()
            th = threading.Thread(target=self.update_thread,args=(info,lock,out))
            th.daemon = True
            th.start()
            threads.append(th)
        for thread in threads:
            thread.join()
    def update_thread(self,info,lock,out):
        b = self.__Book.new(**info)
        lock.acquire()
        try:
            out("Update", "-"*14)
            out(self.identify, "\t", b.book_num, "\t", b.name)
            if (not b.updated):
                cursor = self.conn.cursor()
                sql = "update books set date='"+b.date+"', chapter='"+b.last_chapter+"', end='false'"
                if (info.download):
                    sql += ", version="+str(info.version + 1)
                if (b.name != info.name):
                    cursor.execute("insert into error (type, site, num) values ('name changed <" + info.name + ">-<" + b.name + ">', '" + self.identify + "', " + str(b.book_num) + ")")
                condition = " where site='"+self.identify+"' and num="+b.book_num
                cursor.execute(sql + condition)
                self.conn.commit()
                out("updated")
            else:
                out("skip")
        finally:
            lock.release()
            self.threads_controller.release()
    def explore(self,n,out):
        book_num = 0
        # get largest book num
        tran = self.conn.cursor().execute("select num from books where site ='"+self.identify+"' order by num desc").fetchone()
        if (tran): book_num = tran[0]
        book_num += 1
        # init data for the explore thread
        self.error_page = -1 if ((self.identify == 'bestory') and (book_num < 99999)) else 0
        print(self.error_page)
        self.threads_controller = threading.Semaphore(MAX_THREAD)
        lock = threading.Lock()
        threads = []
        # explore in threads
        while(self.error_page<n):
            if (self.threads_controller.acquire()):
                th = threading.Thread(target=self.explore_thread,args=(book_num,lock,out))
                th.daemon = True
                book_num += 1
                th.start()
                threads.append(th)
        # confirm all thread are finish
        for thread in threads:
            thread.join()
    def explore_thread(self,num,lock,out):
        # get book info
        b = self.__Book.new(book_num=num)
        # try to update it
        lock.acquire()
        try:
            cursor = self.conn.cursor()
            out("explore", "-"*13)
            out(self.identify, "\t", b.book_num, "\t", b.name)
            if (not b.updated):
                flag = bool(cursor.execute("select * from books where name='"+b.name+"' and site='"+self.identify+"' and num="+b.book_num).fetchone())
                if (not(flag)):
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
                if ((self.identify != 'bestory') or (num > 99999)):
                    self.error_page = 0
                out("find")
            else:
                # update if the record does not exist in database
                out("error")
                f_books = bool(cursor.execute("select * from books where site='"+self.identify+"' and num="+b.book_num).fetchone())
                f_error = bool(cursor.execute("select * from error where site='"+self.identify+"' and num="+b.book_num).fetchone())
                if ((not f_books)and(not f_error)):
                    cursor.execute("insert into error (site,num) values ('"+self.identify+"','"+b.book_num+"')")
                    self.conn.commit()
                if ((self.error_page >= 0) or not((num < 99999) and (self.identify == 'bestory'))):
                    self.error_page += 1
        finally:
            lock.release()
            self.threads_controller.release()       
    def error_update(self,out):
        out(self.identify+"="*15)
        cursor = self.conn.cursor()
        self.threads_controller = threading.Semaphore(MAX_THREAD)
        lock = threading.Lock()
        threads = []
        for result in cursor.execute("select site, num from error where site='"+self.identify+"'"):
            self.threads_controller.acquire()
            num = result[1]
            th = threading.Thread(target=self.explore_thread,args=(num,lock,out))
            th.daemon = True
            th.start()
            threads.append(th)
        for thread in threads:
            thread.join()
    def __del__(self):
        self.conn.close()
