import urllib.request, gzip
import re
import io, os, time, copy
import threading

MAX_THREAD = 300

class Database:
    def __init__(self, connection):
        self.connection = connection
    def execute(self, sql, values=(), lock=None):
        if (lock):
            lock.acquire()
        result = self.connection.cursor().execute(sql, values)
        if (lock):
            lock.release()
        return result
    def load(self, book_factory, site, book_num):
        row = self.connection.execute("select name, writer, date, chapter, type, version, end, download, read from books where site=? and num=? order by version desc", (site, book_num)).fetchone()
        info = {
            'site' : site,
            'book_num' : book_num,
            'name' : row[0],
            'writer' : row[1],
            'date' : row[2],
            'last_chapter' : row[3],
            'book_type' : row[4],
            'version' : row[5],
            'end_flag' : row[6],
            'download_flag' : row[7],
            'read_flag' : row[8]
        }
        return book_factory.new(**info)
    def save_book(self, book, lock=None):
        sql = (
            "insert into books"
            "(site, num, name, writer, date, chapter, type, version, end, download, read)"
            "values"
            "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
        )
        if (lock):
            lock.acquire()
        self.connection.execute(sql, (book.site, book.book_num, book.name, book.writer, book.date, book.last_chapter, book.book_type, book.version, book.end_flag, book.download_flag, book.read_flag))
        self.connection.commit()
        if (lock):
            lock.release()
    def update_book(self, book, lock=None):
        version = self.connection.execute("select version from books where books.site=? and books.num=? order by version desc", (book.site, book.book_num)).fetchone()
        if (version == None):
            raise RuntimeError("record not found")
        else:
            version = version[0]
        sql = (
            "update books "
            "set name=?, writer=?, date=?, chapter=?, type=?, version=?, end=?, download=?, read=? "
            "where site=? and num=? and version=?"
        )
        if (lock):
            lock.acquire()
        self.connection.execute(sql, (book.name, book.writer, book.date, book.last_chapter, book.book_type, book.version, book.end_flag, book.download_flag, book.read_flag, book.site, book.book_num, version))
        self.connection.commit()
        if (lock):
            lock.release()
    def delete_book(self, book, lock=None):
        sql = "delete from books where site=? and num=? and version="
        if (lock):
            lock.acquire()
        self.connection.execute(sql, (book.site, book.book_num, book.version))
        self.connection.commit()
        if (lock):
            lock.release()
    def save_error(self, book, error_type="", lock=None):
        sql = "insert into error (site, num, type) values (?, ?, ?)"
        if (lock):
            lock.acquire()
        self.connection.execute(sql, (book.site, book.book_num, error_type))
        self.connection.commit()
        if (lock):
            lock.release()
    def update_error(self, book, error_type="", lock=None):
        sql = "update error set type=? where site=? and num=?"
        if (lock):
            lock.acquire()
        self.connection.execute(sql, (error_type, book.site, book.book_num))
        self.connection.commit()
        if (lock):
            lock.release()
    def delete_error(self, book, lock=None):
        sql = "delete from error where site=? and num=?"
        if (lock):
            lock.acquire()
        self.connection.execute(sql, (book.site, book.book_num))
        self.connection.commit()
        if (lock):
            lock.release()

class Book():
    def __init__(self, **kwargs):
        self.site = kwargs["site"]
        self.book_num = kwargs["book_num"]
        self.base_web = kwargs["base_web"]
        self.__download_web = kwargs["download_web"]
        self.__chapter_web = kwargs["chapter_web"]
        self.name = kwargs["name"]
        self.writer = kwargs["writer"]
        self.date = kwargs["date"]
        self.last_chapter = kwargs["last_chapter"]
        self.book_type = kwargs["book_type"]
        self.version = kwargs["version"]
        self.end_flag = kwargs["end_flag"]
        self.download_flag = kwargs["download_flag"]
        self.read_flag = kwargs["read_flag"]
        self.__decode = kwargs["decode"]
        self.__timeout = kwargs["timeout"]
    def __str__(self):
        return "Site:\t" + self.site + "\nNum:\t" + str(self.book_num) + "\nVersion:\t" + str(self.version) + "\nName:\t" + self.name + "\nWriter:\t" + self.writer + "\nType:\t" + self.book_type
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
    def update(self):
        # read website
        changed = False
        try:
            content = self.open_website(self.base_web)
        except (urllib.error.HTTPError, urllib.error.URLError, urllib.request.socket.timeout):
            raise RuntimeError("Unable to open the website for basic information")
        # check the website book information
        name = self._cut_name(content)
        if ((not self.name) or (self.name != name)):
            # get name
            self.name = name
            changed = True
        if (not self.writer):
            # get writer (writer)
            self.writer = self._cut_writer(content)
            changed = True
        # get date (always get)
        date = self._cut_date(content)
        if (self.date != date):
            self.date = date
            changed = True
        # get chapter (always get)
        last_chapter = self._cut_last_chapter(content)
        if(self.last_chapter != last_chapter):
            self.last_chapter = last_chapter
            changed = True
        if (not self.book_type):
            # check type (bookType)
            self.book_type = self._cut_book_type(content)
            changed = True
        return changed
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
        '''
        for i in range(min(len(titles), len(chapters))):
            chapter_text = self.__download_chapter(titles[i], chapters[i], out)
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
        for i in range(min(len(titles), len(chapters))):
            self.threads_control.acquire()
            t = threading.Thread(target=self.__download_chapter_thread, 
                args=(titles[i], chapters[i], out, arr, lock))
            t.daemon = True
            t.start()
            threads.append(t)
        for t in threads:
            t.join()
        # order the chapter by its url
        for i in range(len(arr)):
            if(arr[i][1] == ''):
                out("download error")
                return False
        text = ""
        for i in range(min(len(titles), len(chapters))):
            for j in range(len(arr)):
                if(chapters[i] == arr[j][0]):
                    text += arr[j][1]
                    break
        
        # save actual content
        try: os.mkdir(path)
        except: pass
        file_name = path + '/' + self.book_num + '.txt' if (self.version == 0) else path + '/' + self.book_num + '-' + self.version + '.txt'
        f = open(file_name, "w", encoding="utf8")
        f.write(text)
        f.close()
        out("download success")
        return True
    def __download_chapter(self,chapter_title,chapter_url,out):
        # check valid url pattern or not
        m = re.match(self.__chapter_web,chapter_url)
        if(not m): return '\n' + '-'*20 + '\n' + chapter_title + '\n' + '-'*20 + '\n'
        # open chapter url
        while(True):
            try:
                content = self.open_website(chapter_url)
                break
            except :
                # wait for a while and try again
                time.sleep(self.__timeout//10)
                out("Reload chapter content")
        if(not content):
            raise RuntimeError("cannot download chapter")
        # read title
        title = self._cut_chapter_title(content)
        # read content
        content = self._cut_chapter_content(content)
        if (content == 'error'):
            return content
        return  '\n' + '-'*20 + '\n' + chapter_title + '\n' + '-'*20 + '\n' + title + content + '\n'
    def __download_chapter_thread(self,chapter_title,chapter_url,out,arr,lock):
        # check url pattern
        m = re.match(self.__chapter_web, chapter_url)
        if(not m):
            lock.acquire()
            arr.append((chapter_url, '\n' + '-'*20 + '\n' + chapter_title + '\n' + '-'*20 + '\n'))
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
                out("Reload chapter_url " + chapter_url)
        if(not content):
            lock.acquire()
            arr.append((chapter_url, ""))
            lock.release()
            self.threads_control.release()
            return
        # read title
        title = self._cut_chapter_title(content)
        # read content
        content = self._cut_chapter_content(content)
        lock.acquire()
        if (content != 'error'):
            # put the result into common array
            arr.append((chapter_url, '\n' + '-'*20 + '\n' + chapter_title + '\n' + '-'*20 + '\n' + title + content + '\n'))
        else:
            arr.append((chapter_url, ''))
        lock.release()
        self.threads_control.release()
        return

class BookFactory():
    def __init__(self, **kwargs):
        self.base_web = kwargs["base_web"]
        self.__download_web = kwargs["download_web"]
        self.__chapter_web = kwargs["chapter_web"]
        self.__book_product = kwargs["book_product"]
        self.__decode = kwargs["decode"] if ("decode" in kwargs) else "utf8"
        self.__timeout = kwargs["timeout"] if ("timeout" in kwargs) else 30
    def new(self,**kwargs):
        info = {
            'site' : kwargs["site"],
            'book_num' : str(kwargs["book_num"]),
            'base_web' : self.base_web.format(kwargs["book_num"]),
            'download_web' : self.__download_web.format(kwargs["book_num"]),
            'chapter_web' : self.__chapter_web,
            'name' : kwargs["name"] if ("name" in kwargs) else "",
            'writer' : kwargs["writer"] if ("writer" in kwargs) else "",
            'date' : kwargs["date"] if ("date" in kwargs) else "",
            'last_chapter' : kwargs["last_chapter"] if ("last_chapter" in kwargs) else "",
            'book_type' : kwargs["book_type"] if ("book_type" in kwargs) else "",
            'version' : kwargs["version"] if ("version" in kwargs) else 0,
            'end_flag' : kwargs["end_flag"] if ("end_flag" in kwargs) else 'false',
            'download_flag' : kwargs["download_flag"] if ("download_flag" in kwargs) else 'false',
            'read_flag' : kwargs["read_flag"] if ("read_flag" in kwargs) else 'false',
            'decode': self.__decode,
            'timeout' : self.__timeout
        }
        return self.__book_product(**info)

class BookSite():
    def __init__(self, **kwargs):
        self.__book_factory = kwargs["book_factory"]
        self.db = kwargs["db"]
        self.path = kwargs["path"]
        self.identify = kwargs["identify"]
        self.running_thread = 0
        self.error_page = 0
    def __str__(self):
        return type(self.__book_factory.__book_product).__name__+' Site'
    def download(self, out):
        out(self.identify+"="*15)
        for result in self.db.execute("select site, num from books where site=? and end='true' and download='false' group by site, num order by date", (self.identify,)).fetchall():
            try:
                b = self.db.load(self.__book_factory, result[0], result[1])
            except RuntimeError as err:
                out(b.site + "\t" + b.num + "\t" + b.version + "\t" + b.name + "\t" + str(err))
                continue
            out(self.identify+"\t"+b.book_num+"\t"+b.name+"\t"+"-"*15)
            if(not b.update()):
                if(b.download(self.path + self.identify, out)):
                    out("Download Successfully")
                    b.download = True
                    self.db.update_book(b)
                else:
                    out("Download Error")
                    b.download_flag = 'error'
                    self.db.update(b)
            else:
                out("Not Updated")
                b.end_flag = 'false'
                self.db.update_book(b)
    def download_thread(self, site, num, lock, out):
        b = self.db.load(self.__book_factory, site, num)
        if(not b.update()):
            if(b.download(self.path + self.identify, out)):
                out("Download Successfully")
                b.download = True
                self.db.update_book(b, lock=lock)
            else:
                out("Download Error")
                b.download_flag = 'error'
                self.db.update_book(b, lock=lock)
        else:
            out("Not Updated")
            b.end_flag = 'false'
            self.db.update_book(b, lock=lock)
    def book(self, book_num):
        return self.db.load(self.__book_factory, self.identify, book_num)
    #TODO make the query function usable
    def query(self, **kwargs):
        query = "(site='"+self.identify+"')"
        if("name" in kwargs):
            query += " and (name='"+kwargs["name"]+"')"
        if("writer" in kwargs):
            query += " and (writer='"+kwargs["writer"]+"')"
        if("book_type" in kwargs):
            query += " and (type='"+kwargs["type"]+"')"
        sql = "select num,name,writer,type from books where "+query
        result = self.db.execute(sql).fetchall()
        for i in range(len(result)):
            result[i] = {
                "num":result[i][0],
                "name":result[i][1],
                "writer":result[i][2],
                "type":result[i][3]
            }
        return result
    def update(self, update_all, out):
        self.threads_controller = threading.Semaphore(MAX_THREAD)
        lock = threading.Lock()
        threads = []
        sql = "select site, num from books where site='"+self.identify+"'"
        if (not update_all):
            sql += " and (read='false' or read is null)"
        for result in self.db.execute(sql + " order by date desc"):
            self.threads_controller.acquire()
            th = threading.Thread(target=self.update_thread, args=(result[0], result[1], lock, out))
            th.daemon = True
            th.start()
            threads.append(th)
        for thread in threads:
            thread.join()
    def update_thread(self, site, book_num, lock, out):
        b = self.db.load(self.__book_factory, site, book_num)
        try:
            out("Update", "-"*14)
            out(b.site, "\t", b.book_num, "\t", b.name)
            original_name = b.name
            original_writer = b.writer
            if(b.update()):
                if ((b.name != original_name) or (b.writer != original_writer)):
                    print("name changed", b.name, original_name, b.writer, original_writer)
                    b.version += 1
                    b.end_flag = 'false'
                    b.download_flag = 'false'
                    self.db.save_book(b, lock=lock)
                    out(b.site + " " + str(b.book_num) + " content changed")
                    return
                if (b.download_flag.upper() == "TRUE"):
                    b.version += 1
                    b.end_flag = 'false'
                    b.download_flag = 'new update'
                self.db.update_book(b, lock=lock)
                out(b.site + " " + str(b.book_num) + " updated")
            else:
                out(b.site + " " + str(b.book_num) + " skip")
        finally:
            self.threads_controller.release()
    def explore(self, n, out):
        book_num = 0
        # get largest book num
        tran = self.db.execute("select num from books where site=? order by num desc", (self.identify,)).fetchone()
        if(tran): book_num = tran[0]
        book_num += 1
        # init data for the explore thread
        self.error_page = -1 if ((self.identify == 'bestory') and (book_num < 99999)) else 0
        self.threads_controller = threading.Semaphore(MAX_THREAD)
        lock = threading.Lock()
        threads = []
        # explore in threads
        while(self.error_page<n):
            self.threads_controller.acquire()
            th = threading.Thread(target=self.explore_thread, args=(self.identify, book_num, lock, out))
            th.daemon = True
            book_num += 1
            th.start()
            threads.append(th)
            # try to reduce he memory that code takes
            if (book_num % 1500 == 0):
                for thread in threads:
                    if (not thread.is_alive()):
                        thread.join()
                        threads.remove(thread)
                print("threads list clean : ", len(threads))
        # confirm all thread are finish
        for thread in threads:
            thread.join()
    def explore_thread(self, site, book_num, lock, out):
        # get book info
        b = self.__book_factory.new(site=site, book_num=book_num)
        # try to update it
        try:
            out("explore", "-"*13)
            out(self.identify, "\t", b.book_num, "\t", b.name)
            if(b.update()):
                flag = bool(self.db.execute("select * from books where name=? and site=? and num=?", (b.name, b.site, b.book_num)).fetchone())
                if(not(flag)):
                    # add record to books table
                    self.db.save_book(b, lock=lock)
                    # delete record from error
                    self.db.delete_error(b, lock=lock)
                if ((self.identify != 'bestory') or (book_num > 99999)):
                    self.error_page = 0
                out("find")
            else:
                # update if the record does not exist in database
                out("error" + str(self.error_page))
                books_found = bool(self.db.execute("select * from books where site=? and num=?", (b.site, b.book_num)).fetchone())
                error_found = bool(self.db.execute("select * from error where site=? and num=?", (b.site, b.book_num)).fetchone())
                if((not books_found) and (error_found)):
                    self.db.save_error(b, lock=lock)
                if ((self.error_page >= 0) or not ((book_num < 99999) and (self.identify == 'bestory'))):
                    self.error_page += 1
        except:
                # update if the record does not exist in database
                out("error" + str(self.error_page))
                books_found = bool(self.db.execute("select * from books where site=? and num=?", (b.site, b.book_num)).fetchone())
                error_found = bool(self.db.execute("select * from error where site=? and num=?", (b.site, b.book_num)).fetchone())
                if((not books_found) and (error_found)):
                    self.db.save_error(b, lock=lock)
                if ((self.error_page >= 0) or not ((book_num < 99999) and (self.identify == 'bestory'))):
                    self.error_page += 1
        finally:
            self.threads_controller.release()       
    def error_update(self,out):
        out(self.identify + "="*15)
        self.threads_controller = threading.Semaphore(MAX_THREAD)
        lock = threading.Lock()
        threads = []
        for result in self.db.execute("select site, num from error where site=?", (self.identify,)):
            self.threads_controller.acquire()
            site = result[0]
            num = result[1]
            th = threading.Thread(target=self.explore_thread, args=(site, num, lock, out))
            th.daemon = True
            th.start()
            threads.append(th)
        for thread in threads:
            thread.join()
