class Book():
    def __init__(self, name='', writer='', date='', chapter='', website='', bookType=''):
        self.name = name
        self.writer = writer
        self.date = date
        self.chapter = chapter
        self.website = website
        self.bookType = bookType


class BaseBook():
    def __init__(self, web, name="", writer="", date="", chapter="", bookType=""):
        self._website = web
        self._name = name
        self._writer = writer
        self._date = date           # last update date
        self._chapter = chapter     # last chapter
        self._bookType = bookType
        self._chapterSet = []
        self._downloadAddr = ""
        self._text = ""
        self._getBasicInfo()
    def _getBasicInfo():
        pass
    def DownloadBook():
        # download and save the book
        pass
    def Update():
        # check any info can be update (date, chapter)
        pass