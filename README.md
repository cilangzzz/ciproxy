

<div align=center></br></br></br>

<center> <img src="https://thirdqq.qlogo.cn/g?b=sdk&k=iaNcdgTAPWOS0JJseiafW1Dw&kti=ZIsqGgAAAAI&s=40&t=1638804590" style="zoom:300%;" /></center>

#  <center>  CiProxy </center>

###### <center>Introduce</center>

###### 							中文文档请看README.cn.md

<center>use tcp implement, not http server, base on golang</center>
</div>

###### install

```makefile
#makefile for server

apt install go

git clone https://github.com/OpencvLZG/CiProxy

#for server
cd CiProxy/server 

go build main.go

#for client
cd CiProxy/client

go build main.go
```

###### running

```makefile
#for server
cd CiProxy/server 
# -ip (server ip default all interface) -port (server port) -method (direct and tunnel)
./main

#for client(need enable proxy server maually)
cd CiProxy/client
# -ip (server ip) -port (server port)
./main

```











###### contact 

- [bilibili](https://space.bilibili.com/433915419)
- [twitter]()
- [slack]()















