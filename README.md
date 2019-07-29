# python spider

#=====class=====#
  #===Book===#
    - init     : with most basic info (to declear which website it belongs to)(similar to a book factory)
    - new      : get a specific book of that website
    - download : download the online book to local storage
  #===Book Site===#
    - init         : with local database information and specify the type of book website it is
    - download     : download valid books
    - update       : update all books in the data base (except the books had been read)
    - explore      : check any new books posted on the website
    - error update : check any update for the book website which has error before
#=====controller=====#
  * sites          : the collection of book sites
  - download     : download valid books
  - update       : update all books in the data base (except the books had been read)
  - explore      : check any new books posted on the website
  - check_end    : check books' last chapter and last update time to define it is end or not
  - error update : check any update for the book website which has error before
  * if the controller is called in cmd, it work as a program (details check "py controller.py --help")
#==========main-cmd==========#
  commend line interface of the program.