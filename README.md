

<div align=center></br></br></br>
<center> <img src="http://www.cilang.buzz/static/file/logo.jpg" zoom="100%"  width="300" height="300" /></center>

#  <center>  CiProxy </center>

###### <center>Introduce</center>

<center>use tcp implement, not http server, base on golang</center>

</div>


![golang](https://img.shields.io/badge/golang-blue?logo=go) ![request](https://img.shields.io/badge/proxy-blue?logo=proxy)  ![open sourse (shields.io)](https://img.shields.io/badge/open%20sourse-darkgreen?logo=opensourceinitiative)   ![github (shields.io)](https://img.shields.io/badge/github-grey?logo=github) ![gitee (shields.io)](https://img.shields.io/badge/gitee-orange?logo=gitee) ![git (shields.io)](https://img.shields.io/badge/git-lightblue?logo=git) ![Mit: license (shields.io)](https://img.shields.io/badge/Mit-license-blue?logo=bookstack) ![img](https://komarev.com/ghpvc/?username=cilproxy&&style=flat-square)


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











