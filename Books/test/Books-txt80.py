import os, sys
sys.path.append(os.path.dirname(os.path.dirname(os.path.realpath(__file__))))
import sqlite3
import txt80
os.system("copy Books\\test\\demo.db Books\\test\\test.db")
test = txt80.TXT80(sqlite3.connect("Books\\test\\test.db"),os.getcwd()+"\\Books\\test\\test book\\txt80")
# test get book
test.Explore(1)
# test updated book
test.Update()
# test download books
test.Download()