

<div align=center></br></br></br>
<center> <img src="http://www.cilang.buzz/static/file/logo.jpg" zoom="100%"  width="300" height="300" /></center>

#  <center>  CiProxy </center>

###### <center>Introduce</center>

<center>ciproxy 是一个基于tcp实现的go语言代理框架，暂时没有使用其他包</center>

<center>use tcp implement, not http server, base on golang</center>

</div>


![golang](https://img.shields.io/badge/golang-blue?logo=go) ![request](https://img.shields.io/badge/proxy-blue?logo=proxy)  ![open sourse (shields.io)](https://img.shields.io/badge/open%20sourse-darkgreen?logo=opensourceinitiative)   ![github (shields.io)](https://img.shields.io/badge/github-grey?logo=github) ![gitee (shields.io)](https://img.shields.io/badge/gitee-orange?logo=gitee) ![git (shields.io)](https://img.shields.io/badge/git-lightblue?logo=git) ![Mit: license (shields.io)](https://img.shields.io/badge/Mit-license-blue?logo=bookstack) ![img](https://komarev.com/ghpvc/?username=cilproxy&&style=flat-square)

## Features

1. 纯原生




## Todo(需要实现)

> 个人没什么空余时间，希望可以有更多的人一起开发


+ [ ] restruct code style(重构代码风格规范,部分代码注释目前没有好的管控) 
+ [ ] impelment proxy server controller(实现一个代理切换控制台)
+ [ ] impelment usage example
+ [ ] optimize performance, add sync
+ [ ] impelment test example


## Install(安装)

```makefile
go get github.com/opencvlzg/ciproxy
```
## Example(示例)

> (internal implement) ciproxy-master/cmd 
```makefile
go build ./cmd/proxyType/
./main -ip -port -method 
```

> (external implement) ciproxy-example/cmd

## External-Usage(使用)

```go
// main implement a simple httpsproxy(only https)

package main

import (
	"flag"
	"github.com/opencvlzg/ciproxy/constants/proxyMethod"
	"github.com/opencvlzg/ciproxy/proxyServer/middleHandle"
	"github.com/opencvlzg/ciproxy/proxyServer/serve"
	"net"
)

func main() {
	ip := flag.String("ip", "127.0.0.1", "Server Ip Address")
	port := flag.String("port", "6677", "Server Port")
	method := flag.String("method", proxyMethod.HttpsProxy, "Server METHOD NORMAL,TUNNEL, SNIFF")
	protocol := flag.String("protocol", "TCP", "Connect Protocol")
	logPath := flag.String("log", "log/proxy.log", "log file path")
	flag.Parse()
	// proxyServer := ciproxy.NewProxyServe()
	proxyServe := serve.ProxyServe{
		Ip:       *ip,
		Port:     *port,
		Method:   *method,
		Protocol: *protocol,
		LogPath:  *logPath,
	}
	middleHandle.Add(func(client net.Conn, target net.Conn) {
		// Todo Some regular u want implement
	})
	proxyServe.Start()
}

```

### HttpProxy

```go

```

### HttpsProxy

```go

```

### HttpsSniffProxy

```go

```


### TunnelProxy

```go

```

## Directory(目录说明)
```
├─cert (证书目录)
│      private.pem
│      root.crt
│      root.csr
│
├─cmd (例子目录)
│  ├─customProxyServer
│  │      main.go
│  │
│  ├─generateCert
│  │      common.txt
│  │      main.go
│  │
│  ├─httpsProxyServer
│  │      main.go
│  │
│  ├─httpsSniffProxyServer
│  │      main.go
│  │
│  ├─tunnelProxyClient
│  │      main.go
│  │
│  ├─tunnelProxyServer
│  │      main.go
│  │
│  └─websocketProxyServer
│          main.go
│
├─internal (内部包)
│  ├─middleware
│  └─util  (常用工具)
│          genHttp.go   (http内容生成)
│          regex.go	    (旧版解析信息用)
│          tlsUtil.go   (tls配置生成)
│          util.go	   
│          windowApi.go (windowApi业务)
│
├─log  (默认日志生成路径)
│      proxy.log
│
├─test
|       main.go
│  ciproxy.go			
│  constants.go			(配置常量)
│  context.go			(代理上下文)
│  data.txt				(例子生成的文件)
│  go.mod		
│  LICENSE
│  proxyHandle.go		 (代理请求处理)
│  README.md
│  serve.go				 (程序主入口)
│  serveHandle.go         (监听服务处理)
│  trafficHandle.go		  (网络流浪处理)


```

## Contact(联系)

> ## ciproxy communicate
>
>> 
>> google email cilanguser@gmail.com 
>> - [bilibili](https://space.bilibili.com/433915419)
>> - [twitter]()
>> - [slack]()


## Support(支持)

> CiPorxy遵循Mit开源协议,如果你想请支持,可以点我











