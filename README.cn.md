

<div align=center></br></br></br>

<center> <img src="https://thirdqq.qlogo.cn/g?b=sdk&k=iaNcdgTAPWOS0JJseiafW1Dw&kti=ZIsqGgAAAAI&s=40&t=1638804590" style="zoom:300%;" /></center>

#  <center>  CiProxy </center>

###### <center>Introduce</center>

<center>use tcp implement, not http server, base on golang</center>
</div>

###### 安装

```makefile
#服务器

apt install go

git clone https://github.com/OpencvLZG/CiProxy

#服务端
cd CiProxy/server 

go build main.go

#客户端
cd CiProxy/client

go build main.go
```

###### running

```makefile
#服务器
cd CiProxy/server 
# -ip 服务器ip默认全部接口 -port 服务器监听端口 -method 直连模式与加密模式
./main


#for client(need enable proxy server maually)
cd CiProxy/client
# -ip 服务器ip -port 服务器端口
./main

```











###### contact 

- [bilibili](https://space.bilibili.com/433915419)
- [twitter]()
- [slack]()















