# Go lang version
----

i found that this program in python is slow and take so many memory when using sqlite3
then i found another language to make the program again
here comes Go lang version with faster speed and less memory consumption

to compile the controller, install go first, then run 
```ternimal
go get -d ./...
go build ./controller.go
```
in terminal under the `go-lang` folder.

for the usage details of controller, use `./controller` to check.

# structure
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
- helper
	- helper.go
		* provide helper function (eg. regex, get url response...)
- frontend
	* it is a react folder containing all fontend needed
- frontend_flutter
	* it is a flutter folder containing all frontend needed

# POC
- compare `http` and `grequest
- compare  json  and  yaml  as config
- compare  regex  and  html parser library
`
