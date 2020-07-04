import os, sqlite3, urllib.request, collections, datetime, json
import threading, concurrent.futures
import gc

MAX_THREAD = 20

class Database:
    def __init__(self, db_path):
        self.db_path = db_path
        if (not os.path.exists(db_path)):
            path = os.getcwd()
            path = path[:path.find('BookSpider')] + 'BookSpider/Books2/database/template.db'
            src = open(path, 'rb')
            dest = open(db_path, 'wb')
            dest.write(src.read())
            dest.close()
            src.close()
        self.db_conn = sqlite3.connect(db_path, check_same_thread=False)
    def exist(self, table, record):
        result = self.db_conn.execute('select * from ' + table +' where site=? and num=?', [record.site, record.num]).fetchone()
        return (result != None)
    def get_record(self, table, condition):
        #db_conn = sqlite3.connect(self.db_path)
        result = self.db_conn.cursor().execute('select * from ' + table + ' ' + condition)
        return result
    def add_record(self, table, record):
        if (self.exist(table, record)):
            self.update_record(table, record)
        elif (self.exist('error', record) and (table == 'books')):
            self.move_record('error', 'books', record)
            self.update_record(table, record)
        elif (self.exist('books', record) and (table == 'error')):
            self.move_record('books','error', record)
            self.update_record(table, record)
        else:
            self.db_conn.execute('insert into ' + table + ' values ' + str(record.to_record(table)))
        self.db_conn.commit()
    def update_record(self, table, record):
        self.db_conn.execute('update ' + table + ' set ' + record.to_value(table) + ' where site=? and num=?', [record.site, record.num])
        self.db_conn.commit()
    def delete_record(self, table, record):
        self.db_conn.execute('delete from ' + table + ' where site=? and num=?', [record.site, record.num])
        self.db_conn.commit()
    def move_record(self, source_table, target_table, record):
        if (not self.exist(source_table, record)):
            return
        self.delete_record(source_table, record)
        self.add_record(target_table, record)
        self.db_conn.commit()
    def check_end(self):
        criteria = ["后记", "後記", "新书", "新書", "结局", "結局", "感言", 
                "尾声", "尾聲", "终章", "終章", "外传", "外傳", "完本",
                "结束", "結束", "完結", "完结", "终结", "終結", "番外",
                "结尾", "結尾", "全书完", "全書完", "全本完"]
        sql = "update books set end='true', download='false' where ("
        for c in criteria:
            sql += "chapter like '%"+c+"%' or "
        sql += "date < '"+str(datetime.datetime.now().year-2)+"') and (end <> 'true' or end is null)"
        c = self.db_conn.execute(sql)
        self.db_conn.commit()
        print(str(c.rowcount)+" row affected")
    def backup(self):
        #TODO: backup the database from database folder to backup folder
        # create yyyymmdd folder if not exist
        path = __file__.replace('models.py', 'backup/') + str(datetime.date.today()) + '/'
        if (not os.path.exists(path)):
            os.mkdir(path)
        path += self.db_path[self.db_path.rfind('/')+1:].replace('.db', '-'+datetime.datetime.now().strftime('%H-%M-%S')+'.db')
        # copy database to folder and name as site-hhmmss.db
        src = open(self.db_path, 'rb')
        dest = open(path, 'wb')
        dest.write(src.read())
        dest.close()
        src.close()
        pass

class Logger:
    def __init__(self, db_path):
        self.db_path = db_path
        self.db_conn = sqlite3.connect(db_path, check_same_thread=False)
    def _create_table(self):
        self.db_conn.execute("create table if not exists log (time real, raw text)")
    def log(self, **raw):
        raw = json.dumps(raw)
        time = datetime.datetime.now().timestamp()
        self.db_conn.execute("insert into log (time, raw) values (?, ?)", (time, raw))
        self.db_conn.commit()
class Book:
    def __init__(self, db=None, base_url='', download_url='', chapter_url='', site='', num=0, decode='utf8', max_thread = MAX_THREAD, timeout=30):
        self.base_url = base_url
        self.download_url = download_url
        self.chapter_url = chapter_url
        self.site = site.lower()
        self.num = num
        self.decode = decode
        self.max_thread = max_thread
        self.timeout = timeout
        if ((db != None) and (db.exist('books', self))):
            row = db.get_record('books', 'where site="'+self.site+'" and num='+str(self.num)+' order by version desc limit 2')
            record = row.fetchone()
            row.close()
            self.title = record[0]
            self.writer = record[1]
            self.last_update = record[2]
            self.last_chapter = record[3]
            self.book_type = record[4]
            self.end_flag = record[5]
            self.download_flag = record[6]
            self.read_flag = record[7]
            self.version = record[10]
        else:
            self.title = ''
            self.writer = ''
            self.last_update = ''
            self.last_chapter = ''
            self.book_type = ''
            self.end_flag = False
            self.download_flag = False
            self.read_flag = False
            self.version = 0
           
    def _get_title(self, html):
        pass
    def _get_writer(self, html):
        pass
    def _get_type(self, html):
        pass
    def _get_last_update(self, html):
        pass
    def _get_last_chapter(self, html):
        pass
    def update(self):
        try:
            with urllib.request.urlopen(self.base_url, timeout=self.timeout) as res:
                html = res.read().decode(self.decode)
        except:
            return False
        try:
            self.title = self._get_title(html)
            self.writer = self._get_writer(html)
            self.book_type = self._get_type(html)
            last_update = self._get_last_update(html)
            if (self.last_update != last_update):
                self.last_update = last_update
            else:
                '''
                if (self.logger != None):
                    self.logger.log(site=self.site, num=self.num,
                            version=self.version, operation="update", result="fail")
                '''
                return False
            self.last_chapter = self._get_last_chapter(html)
            if (self.end_flag == True):
                self.version += 1
                self.end_flag = False
            '''
            if (self.logger != None):
                self.logger.log(site=self.site, num=self.num,
                        version=self.version, operation="update", result="success")
            '''
            return True
        except (AttributeError, IndexError):
            return False
    def _get_chapters_url(self, html):
        pass
    def _get_chapters_title(self, html):
        pass
    def _get_content(self, html):
        pass
    def _download_chapter(self, url, title, semaphore, lock):
        if (self.chapter_url.replace('\\d', '') not in url):
            #lock.acquire() if (lock != None) else None
            self.content_list.append((self.chapters_url.index(url), title, None))
            #lock.release() if lock != None else None
            semaphore.release() if semaphore != None else None
            return
        for i in range(10):
            try:
                res = urllib.request.urlopen(url, timeout=self.timeout)
                html = res.read().decode(self.decode, 'ignore')
                break
            except:
                if (i == 10):
                    #lock.acquire() if lock != None else None
                    self.content_list.append((self.chapters_url.index(url), title, None))
                    #lock.release() if lock != None else None
                    semaphore.release() if semaphore != None else None
                    return
                print('reload', url)
        content = self._get_content(html)
        #lock.acquire() if lock != None else None
        self.content_list.append((self.chapters_url.index(url), title, content))
        #lock.release() if lock != None else None
        semaphore.release() if semaphore != None else None
        return content
    def download(self, path):
        try:
            res = urllib.request.urlopen(self.download_url, timeout=self.timeout)
            html = res.read().decode(self.decode)
        except:
            '''
            if (self.logger != None):
                self.logger.log(site=self.site, num=self.num,
                            version=self.version, operation="download", result="fail")
            '''
            return False
        self.chapters_url = self._get_chapters_url(html)
        self.chapters_title = self._get_chapters_title(html)
        self.content_list = []
        threads = []
        lock = threading.Lock()
        semaphore = threading.Semaphore(self.max_thread)
        pool = concurrent.futures.ThreadPoolExecutor(MAX_THREAD)
        print('different') if len(self.chapters_title) != len(self.chapters_url) else None
        #'''
        for (url, title) in zip(self.chapters_url, self.chapters_title):
            ### thread version ###
            threads.append(threading.Thread(target=self._download_chapter, args=(url, semaphore, lock)))
            threads[-1].deamon = True
            semaphore.acquire()
            threads[-1].start()

            ### linear version ###
            ''''
            self._download_chapter(url, None, None)
            '''
            print(title)
        for thread in threads:
            thread.join()
        self.content_list.sort(key=lambda item: item[0])
        result = self.title + '\n' + self.writer + '\n' + '-'*20 + '\n'*2
        for (_, chapter_title, chapter_content) in self.content_list:
            ### sort the url list ###
            result += chapter_title + '\n' + '-'*20 + '\n'
            result += chapter_content + '\n'*2

        '''
        for (url, title) in zip(self.chapters_url, self.chapters_title):
            semaphore.acquire()
            self.content_list.append(pool.submit(self._download_chapter, url, semaphore, lock))
        pool.shutdown()
        self.content_list = [ content.result() for content in self.content_list ]
        result = self.title + '\n' + self.writer + '\n' + '-'*20 + '\n'*2
        for (title, content) in zip(self.chapters_title, self.content_list):
            result += title + '\n' + '-'*20 + '\n' + content + '\n'*2
        if ('.txt' not in path):
            path += '/' + str(self.num) + '-v' + str(self.version) + '.txt'
        f = open(path, 'w', encoding='utf8')
        f.write(result)
        f.close()

        if (self.logger != None):
            self.logger.log(site=self.site, num=self.num,
                    version=self.version, operation="download", result="success")
        '''
        return True
    def to_record(self, table):
        if (table == 'books'):
            return (self.title, self.writer, self.last_update, self.last_chapter, self.book_type,
                    str(self.end_flag).lower(), str(self.download_flag).lower(),
                    str(self.read_flag).lower(), self.site, self.num, self.version)
        elif (table == 'error'):
            return ('', self.site, self.num)
    def to_value(self, table):
        if (table == 'books'):
            # TODO set values for books
            return (
                'site="' + self.site + '", '
                'num=' + str(self.num) + ', '
                'name="' + self.title + '", '
                'writer="' + self.writer +'", '
                'date="' + self.last_update + '", '
                'chapter="' + self.last_chapter + '", '
                'type="' + self.book_type + '", '
                'end="' + str(self.end_flag) + '", '
                'download="' + str(self.download_flag) + '", '
                'read="' + str(self.read_flag) + '", '
                'version=' + str(self.version)
            )
            pass
        elif (table == 'error'):
            return 'site="'+self.site+'", num='+str(self.num)

Setting = collections.namedtuple('Setting', ['meta_base_url', 'meta_download_url', 'chapter_url'])

class Site:
    def __init__(self, site, db_path, download_path, setting, decode='utf8', max_thread=MAX_THREAD, timeout=30):
        self.site = site
        self.db = Database(db_path)
        self.download_path = download_path
        self.meta_base_url = setting.meta_base_url
        self.meta_download_url = setting.meta_download_url
        self.chapter_url = setting.chapter_url
        self.decode = decode
        self.max_thread = max_thread
        self.timeout = timeout
    def get_book(self, num):
        pass
    def explore(self, max_error_count):
        last_num = self.db.get_record('books', 'order by num desc').fetchone()[9]
        self.error_count = 0
        lock = threading.Lock()
        semaphore = threading.Semaphore(MAX_THREAD)
        pool = concurrent.futures.ThreadPoolExecutor(MAX_THREAD)
        threads = []
        while (self.error_count < max_error_count):
            '''
            semaphore.acquire()
            pool.submit(self._explore_thread, last_num, semaphore, lock)
            last_num += 1
            '''
            semaphore.acquire()
            thread = threading.Thread(target=self._explore_thread, args=(last_num, semaphore, lock))
            thread.daemon = True
            threads.append(thread)
            thread.start()
            if (len(threads) >= 1000):
                for thread in threads:
                    if (not thread.is_alive()):
                        thread.join()
                        threads.remove(thread)
                gc.collect()
        '''
        pool.shutdown()
        '''
        for thread in threads:
            thread.join()
        gc.collect()

    def _explore_thread(self, num, semaphore, lock):
        b = self.get_book(num)
        if (b.update()):
            self.error_count = 0
            lock.acquire() if (lock != None) else None
            self.db.add_record('books', b)
            print(b.site, b.num, 'success')
            lock.release() if (lock != None) else None
        else:
            self.error_count += 1
            lock.acquire() if (lock != None) else None
            self.db.add_record('error', b)
            print(b.site, b.num, 'failed')
            lock.release() if (lock != None) else None
        semaphore.release() if (semaphore != None) else None
    def update(self):
        rows = self.db.get_record('books', 'where site="' + self.site.lower() + '" order by date desc')
        lock = threading.BoundedSemaphore(1)
        semaphore = threading.BoundedSemaphore(MAX_THREAD)
        threads = []
        '''
        pool = concurrent.futures.ThreadPoolExecutor(MAX_THREAD)
        '''
        for i, row in enumerate(rows):
            '''
            semaphore.acquire()
            pool.submit(self._update_thread, row[9], semaphore, lock)
            '''
            semaphore.acquire()
            thread = threading.Thread(target=self._update_thread, args=(row[9], semaphore, lock))
            thread.daemon = True
            threads.append(thread)
            thread.start()
            if (len(threads) >= 1000):
                for thread in threads:
                    if (not thread.is_alive()):
                        thread.join()
                        threads.remove(thread)
                gc.collect()
        for thread in threads:
            thread.join()
        '''
        pool.shutdown()
        '''
        gc.collect()
    def _update_thread(self, num, semaphore, lock):
        b = self.get_book(num)
        v = b.version
        if (b.update()):
            if (v == b.version):
                lock.acquire() if (lock != None) else None
                self.db.update_record('books', b)
                lock.release() if (lock != None) else None
                #log.log(site=self.site, num=num, message='updated')
                print(self.site, num, 'updated')
            else:
                lock.acquire() if (lock != None) else None
                self.db.add_record('books', b)
                lock.release() if (lock != None) else None
                #log.log(site=self.site, num=num, version=b.version, message='added')
                print(self.site, num, b.version, 'added')
        else:
            #lock.acquire() if (lock != None) else None
            #log.log(site=self.site, num=num, message='unchanged')
            print(self.site, num, 'unchanged')
            #lock.release() if (lock != None) else None
        semaphore.release() if (semaphore != None) else None
    def download(self):
        if (not os.path.exists(self.download_path)):
            os.mkdir(self.download_path)
        rows = self.db.get_record('books', 'where end=true')
        for row in rows:
            b = self.get_book(row[9])
            b.download(self.download_path)
            b.download_flag = True
            self.db.update_record('books', b)
    def update_error(self):
        rows = self.db.get_record('error', 'where site="' + self.site + '"').fetchall()
        lock = threading.Lock()
        semaphore = threading.Semaphore(MAX_THREAD)
        threads = []
        '''
        pool = concurrent.futures.ThreadPoolExecutor(MAX_THREAD)
        '''
        for row in rows:
            '''
            semaphore.acquire()
            pool.submit(self._update_thread, row[2], semaphore, lock)
            '''
            semaphore.acquire()
            thread = threading.Thread(target=self._update_thread, args=(row[2], semaphore, lock))
            thread.daemon = True
            threads.append(thread)
            thread.start()
            if (len(threads) >= 1000):
                for thread in threads:
                    if (not thread.is_alive()):
                        thread.join()
                        threads.remove(thread)
                gc.collect()
        for thread in threads:
            thread.join()
        pool.shutdown()
        gc.collect()
    def _update_error_thread(self, num, semaphore, lock):
            b = self.get_book(num)
            if (b.update()):
                lock.acquire() if (lock != None) else None
                self.db.move_record('error', 'books', b)
                self.db.update_record('books', b)
                lock.release() if (lock != None) else None
            semaphore.release() if (semaphore != None) else None
    def info(self):
        ### get info of the site ###
        '''
        eg. site name, num of normal books, num of error books, total num of books, max num of books
        '''
        print('site\t:\t', self.site)
        normal_rows = self.db.db_conn.execute('select num from books group by num').fetchall()
        print('normal books count\t:\t', len(normal_rows))
        error_rows = self.db.db_conn.execute('select num from error group by num').fetchall()
        print('error books count\t:\t', len(error_rows))
        print('total books count\t:\t', len(normal_rows)+len(error_rows))
        try:
            (normal_num,) = self.db.db_conn.execute('select num from books order by num desc').fetchone()
            (error_num,) = self.db.db_conn.execute('select num from error order by num desc').fetchone()
        except:
            normal_num, error_num = 0,0
        print('max book num\t:\t', max(normal_num, error_num))
    def fix_storage_error(self):
        ### fix storage error ###
        storage_books = os.listdir(self.download_path)
        database_download_books = self.db.db_conn.execute("select num, version, download from books where download=?", ('true',)).fetchall()
        fake_record_books = []
        ### check any download books is not record in database
        for book in storage_books:
            book_info = book.replace('.txt','').split('-v')
            book_info[0] = int(book_info[0])
            if (len(book_info) < 2):
                book_info.append(0)
            find = list(filter(lambda item: item[0]==book_info[0] and item[1]==book_info[1], database_download_books))
            if (len(find) == 0):
                fake_record_books.append(book_info)
            else:
                database_download_books.remove(find[0])
        ### check any books mark download is fake
        print('fake download books : ', len(database_download_books))
        for book in database_download_books:
            self.db.db_conn.execute('update books set download="false" where num=? and version=?', (book[0], book[1]))
        print('wrong record books : ', len(fake_record_books))
        for book in fake_record_books:
            self.db.db_conn.execute('update books set download="true" where num=? and version=?', (book[0], book[1]))
        self.db.db_conn.commit()
    def fix_database_error(self):
        ### check database error ###
        database_books = self.db.db_conn.execute("select site, num, version from books order by num").fetchall()
        ### delete duduplicate record (site, num, version)
        if(len(database_books) > len(set(database_books))):
            print('deduplicate books exist !!!')
            for book in set(database_books):
                database_books.remove(book)
            for book in database_books:
                print(book)
                row = self.db.db_conn.execute('select * from books where site=? and num=? and version=?', book)
                self.db.db_conn.execute('delete from books where site=? and num=? and version=?', book)
                self.db.db_conn.execute('insert into books values (?,?,?,?,?,?,?,?,?,?,?)', row)
            self.db.db_conn.commit()
        database_books = self.db.db_conn.execute("select site, num from error order by num").fetchall()
        ### delete duduplicate record (site, num, version)
        if(len(database_books) > len(set(database_books))):
            print('deduplicate error exist !!!')
            for book in set(database_books):
                database_books.remove(book)
            for book in database_books:
                print(book)
                self.db.db_conn.execute('delete from error where site=? and num=?', book)
                self.db.db_conn.execute('insert into error (site, num) values (?,?)', book)
            self.db.db_conn.commit()
        ### add missing num to error
        database_books = self.db.db_conn.execute('select num from books group by num order by num').fetchall()
        database_books.extend(self.db.db_conn.execute('select num from error order by num').fetchall())
        database_books = [i for (i,) in database_books]
        (max_num,) = self.db.db_conn.execute('select num from books order by num desc').fetchone()
        if (database_books[-1] == len(database_books)):
            return
        all_num = sorted([i for i in range(1,max_num+1)])
        for num in database_books:
            all_num.remove(num)
        self.error_count = 0
        threads = []
        lock = threading.Lock()
        semaphore = threading.Semaphore(300)
        for num in all_num:
            print(self.site, num, 'not found')
            threads.append(threading.Thread(target=self._explore_thread, args=(num, semaphore, lock)))
            threads[-1].daemon = False
            semaphore.acquire()
            threads[-1].start()
        for thread in threads:
            thread.join()
