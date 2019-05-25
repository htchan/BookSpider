class Book():
    def __init__(self, name='', writer='', date='', chapter='', website='', bookType=''):
        self.name = name
        self.writer = writer
        self.date = date
        self.chapter = chapter
        self.website = website
        self.bookType = bookType


class BaseBook():
    def __init__(self, web):
        self._website = web
        self._name = ""
        self._writer = ""
        self._date = ""     # last update date
        self._chapter = ""  # last chapter
        self._bookType = ""
        self._chapterSet = []
        self._downloadAddr = ""
        self._text = ""
        self._getBasicInfo()
    def _getBasicInfo():
        pass
    def DownloadBook():
        pass
    def Save():
        pass