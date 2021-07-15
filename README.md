# Introduction
----
This is a project that update and download book from internet. It crawl data from online site by regex pre-defined in yaml file. It can make a copy of the book content. 

# Go lang version
----

i found that implementing this program in python is slow in multi-threading and taking a lot of memory in querying sqlite3
then i found another language to implement the program again
here comes Go lang version with faster speed and less memory consumption

to compile and run, install docker and run 
```ternimal
make build
# run backend server
make backend

# run controller to operate on database
make controller command=help
```

to enable the frontend features, you have to install flutter and run it.

# backend structure
- main
	- controller.go
		* provide command line control
	- backend.go
		* provide http api control
- model
	- Site.go
		* manage site behavior (eg. get book from database, update all books...)
		* manage book and database communication (eg. save / update the book in database)
	- Book.go
		* manage book behavior (eg. update specific book, generate book object...)
	- Config.go
		* read the config from yaml file and turn it to `model.Site` and `model.Book` class
- helper
	- helper.go
		* provide helper function (eg. regex, get url response...)
- frontend
	* it is a flutter folder containing all frontend needed

# POC
- compare `http` and `grequest
- compare  json  and  yaml  as config
- compare  regex  and  html parser library
- try the flag to read parameter from command line