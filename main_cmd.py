import Books.controller
import os

if(__name__=="__main__"):
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
		if(res == "D"): Books.controller.download(print)
		elif(res == "C"): Books.controller.explore(print)
		elif(res == "U"): Books.controller.update(print)
		elif(res == "A"): Books.controller.error_update(print)
		elif(res == "E"): looping = False
		else:
			print("wrong input")
		input()
		os.system("cls")