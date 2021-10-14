you can use `go run main/main.py` to run the comparision in `src/POC/http-request` folder

according to the result of send 999 `GET` request to an sample api online `https://reqres.in/api/users?page=2` : 
```
net - Linear request
45.809264063s
Alloc =  316695
total error :  0 / 999 

goreq - Linear request
45.994735004s
Alloc =  99159
total error :  0 / 999 

net - go routine request
45.86120809s
Alloc =  99209
total error :  0 / 999 

goreq - go routine request
46.527979428s
Alloc =  99111
total error :  0 / 999 
```

there is no big difference between `net/http` and `gorequest`. The conclusion will be keep using `net/http` as current package response for sending request