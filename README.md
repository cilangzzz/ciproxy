

<div align=center></br></br></br>

<center> <img src="https://thirdqq.qlogo.cn/g?b=sdk&k=iaNcdgTAPWOS0JJseiafW1Dw&kti=ZIsqGgAAAAI&s=40&t=1638804590" style="zoom:300%;" /></center>

#  <center>  CiProxy </center>

###### <center>Introduce</center>

<center>use tcp implement, not http server, base on golang</center>
</div>

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

## HttpProxy

```go

```

## HttpsProxy

```go

```

## HttpsSniffProxy

```go

```


## TunnelProxy

```go

```


## Contact

> ## ciproxy communicate
>
>> 
>> google email cilanguser@gmail.com 
>> - [bilibili](https://space.bilibili.com/433915419)
>> - [twitter]()
>> - [slack]()















