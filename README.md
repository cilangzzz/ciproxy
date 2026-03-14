<div align=center></br></br></br>
<center> <img src="http://www.cilang.buzz/static/file/logo.jpg" zoom="100%"  width="300" height="300" /></center>

#  <center>  CiProxy </center>

###### <center>Introduce</center>

<center>ciproxy 是一个基于 TCP 实现的 Go 语言代理框架，支持 HTTP/HTTPS 代理、MITM 拦截、流量捕获与重放等功能</center>

<center>A Go proxy framework based on TCP, supporting HTTP/HTTPS proxy, MITM interception, traffic capture and replay</center>

</div>

![golang](https://img.shields.io/badge/golang-blue?logo=go) ![request](https://img.shields.io/badge/proxy-blue?logo=proxy)  ![open sourse (shields.io)](https://img.shields.io/badge/open%20sourse-darkgreen?logo=opensourceinitiative)   ![github (shields.io)](https://img.shields.io/badge/github-grey?logo=github) ![gitee (shields.io)](https://img.shields.io/badge/gitee-orange?logo=gitee) ![git (shields.io)](https://img.shields.io/badge/git-lightblue?logo=git) ![Mit: license (shields.io)](https://img.shields.io/badge/Mit-license-blue?logo=bookstack) ![img](https://komarev.com/ghpvc/?username=cilproxy&&style=flat-square)

## Features

- 纯原生 Go 实现，无第三方依赖
- 参考 Gin/Echo 框架设计的中间件系统
- 支持多种代理模式：HTTP、HTTPS、HTTPS Sniff、TCP Tunnel、WebSocket
- MITM 拦截：完整支持 HTTP/1.1 和 HTTP/2 流量拦截
- 流量捕获与重放：支持请求/响应的捕获、存储和重放
- 链式调用 API：优雅的服务器配置方式
- 优雅关闭：支持信号处理和上下文取消
- 配置系统：支持 JSON 配置文件

## Proxy Methods

| Method | Description |
|--------|-------------|
| `HttpProxy` | HTTP 代理 |
| `HttpsProxy` | HTTPS 代理 (CONNECT 隧道) |
| `HttpsSniffProxy` | HTTPS 嗅探代理 |
| `HttpsSniffDetailProxy` | HTTPS 详细嗅探代理 |
| `HttpInterceptProxy` | HTTPS MITM 拦截代理 (支持 HTTP/1.1 & HTTP/2) |
| `TcpTunnelProxy` | TCP 隧道代理 |
| `WebsocketProxy` | WebSocket 代理 |
| `DefaultProxy` | 自定义代理模式 |

## Install

```bash
go get github.com/opencvlzg/ciproxy
```

## Quick Start

### 方式一：链式调用 (推荐)

```go
package main

import (
    "github.com/opencvlzg/ciproxy"
)

func main() {
    // 使用链式调用创建服务器
    proxy := ciproxy.New().
        SetHost("127.0.0.1", "8080").
        SetMethod(ciproxy.HttpsProxy).
        Use(loggingMiddleware).
        Handle(customHandler)

    // 启动服务器
    proxy.Run()
}

func loggingMiddleware(ctx *ciproxy.Context) {
    ctx.Next()
}

func customHandler(ctx *ciproxy.Context) {
    // 自定义处理逻辑
    ctx.Next()
}
```

### 方式二：配置选项

```go
package main

import (
    "github.com/opencvlzg/ciproxy"
)

func main() {
    proxy := ciproxy.New(
        ciproxy.WithHost("127.0.0.1", "8080"),
        ciproxy.WithMethod(ciproxy.HttpInterceptProxy),
        ciproxy.WithMaxConnections(10000),
    )

    proxy.Use(loggingMiddleware)
    proxy.Run()
}
```

### 方式三：向后兼容 API

```go
package main

import (
    "flag"
    "github.com/opencvlzg/ciproxy"
)

func main() {
    ip := flag.String("ip", "127.0.0.1", "Server Ip Address")
    port := flag.String("port", "8080", "Server Port")
    method := flag.String("method", ciproxy.HttpsProxy, "Proxy Method")
    flag.Parse()

    proxyServe := ciproxy.ProxyServe{
        Ip:       *ip,
        Port:     *port,
        Method:   *method,
        Protocol: "TCP",
    }
    proxyServe.AddMiddleware(loggingMiddleware)
    proxyServe.Start()
}
```

## MITM 拦截

```go
package main

import (
    "github.com/opencvlzg/ciproxy"
    "github.com/opencvlzg/ciproxy/pkg/mitm"
    "github.com/opencvlzg/ciproxy/pkg/middleware"
)

func main() {
    // 创建 MITM 拦截器
    interceptor := mitm.NewInterceptor(&mitm.InterceptorConfig{
        EnableTrafficCapture: true,
        EnableHTTP2:          true,
    })

    // 添加拦截中间件
    interceptor.UseFunc(func(ctx *ciproxy.Context) *http.Response {
        // 修改请求或返回自定义响应
        if ctx.GetRequestHeader("X-Custom") == "" {
            ctx.SetRequestHeader("X-Custom", "injected")
        }
        return nil // 返回 nil 放行，返回 response 则阻断请求
    })

    // 启动代理服务器
    proxy := ciproxy.New().SetMethod(ciproxy.HttpInterceptProxy)
    proxy.Run()
}
```

## 流量捕获与重放

```go
package main

import (
    "github.com/opencvlzg/ciproxy/pkg/transfer"
)

func main() {
    store := transfer.GetTrafficStore()

    // 捕获请求
    capturedReq := transfer.CaptureFromRequest(req)
    id := store.Capture(capturedReq)

    // 获取流量
    entry, ok := store.Get(id)

    // 列出所有流量
    entries := store.List()

    // 导出/导入
    data, _ := store.Export()
    store.Import(data)

    // 重放请求
    replayer := transfer.NewReplayer()
    resp, _ := replayer.Replay(entry)
}
```

## 内置中间件

```go
import "github.com/opencvlzg/ciproxy/pkg/middleware/builtins"

// 日志中间件
proxy.Use(builtins.Logging())

// CORS 中间件
proxy.Use(builtins.CORS())

// 请求头修改
proxy.Use(builtins.SetHeader("X-Proxy", "CiProxy"))

// 请求体修改
proxy.Use(builtins.SetBody([]byte("modified body")))

// 请求阻断
proxy.Use(builtins.BlockWithStatus(403, "Forbidden"))
```

## 配置文件

```json
{
  "ip": "127.0.0.1",
  "port": "8080",
  "protocol": "TCP",
  "method": "HttpInterceptProxy",
  "logPath": "log/proxy.log",
  "tls": {
    "certPath": "",
    "keyPath": ""
  },
  "timeout": {
    "connect": "10s",
    "read": "30s",
    "write": "30s"
  },
  "features": {
    "enableTrafficCapture": true,
    "enableTrafficReplay": true
  }
}
```

> **TLS 证书说明**：默认使用内嵌证书，零配置开箱即用。如需使用自定义证书，设置 `certPath` 和 `keyPath` 即可。

加载配置：

```go
config, _ := ciproxy.LoadConfig("config.json")
proxy := ciproxy.New(ciproxy.WithConfig(config))
proxy.Run()
```

## Directory

```
├── cmd/                     # 示例程序
│   ├── custom_proxy_server/ # 自定义代理服务器
│   ├── generate_cert/       # 证书生成工具
│   ├── https_proxy_server/  # HTTPS 代理服务器
│   ├── https_sniff_proxy_server/ # HTTPS 嗅探代理
│   ├── tunnel_proxy_client/ # 隧道代理客户端
│   ├── tunnel_proxy_server/ # 隧道代理服务器
│   └── websocket_proxy_server/ # WebSocket 代理
├── pkg/                     # 核心包
│   ├── context/             # 请求上下文 (参考 Gin 设计)
│   ├── middleware/          # 中间件系统
│   │   └── builtins/        # 内置中间件
│   ├── module/              # 功能模块
│   │   └── mitm/            # MITM 拦截器
│   ├── transfer/            # 流量捕获与重放
│   └── util/                # 工具函数
├── ciproxy.go               # 入口文件
├── config.go                # 配置管理
├── constants.go             # 常量定义
├── errors.go                # 错误处理
├── logger.go                # 日志系统
├── proxy_handle.go          # 代理处理器
├── serve.go                 # 服务器核心实现
├── serve_handle.go          # 服务处理
├── traffic_handle.go        # 流量处理
└── go.mod
```

## API Reference

### ProxyServe

| Method | Description |
|--------|-------------|
| `New(opts...ServerOption)` | 创建服务器实例 |
| `Use(middleware...ProxyHandle)` | 添加中间件 |
| `Handle(handler ProxyHandle)` | 添加处理器 |
| `SetMethod(method string)` | 设置代理模式 |
| `SetHost(ip, port string)` | 设置监听地址 |
| `Start() error` | 启动服务器 |
| `Run() error` | 启动并阻塞 |
| `Shutdown(ctx context.Context) error` | 优雅关闭 |
| `Stats() ServerStats` | 获取统计信息 |

### Context

| Method | Description |
|--------|-------------|
| `Next()` | 执行下一个处理器 |
| `Abort()` | 中断处理器链 |
| `GetRequest() *http.Request` | 获取请求 |
| `GetResponse() *http.Response` | 获取响应 |
| `GetRequestBody() []byte` | 获取请求体 |
| `SetRequestBody([]byte)` | 设置请求体 |
| `GetRequestHeader(key) string` | 获取请求头 |
| `SetRequestHeader(key, value)` | 设置请求头 |
| `BlockWithStatus(code, body)` | 阻断并返回响应 |

## Todo

+ [ ] 完善文档和注释
+ [ ] 实现代理切换控制台
+ [ ] 添加更多使用示例
+ [ ] 性能优化
+ [ ] 单元测试覆盖

## Contact

> google email: cilanguser@gmail.com
> - [bilibili](https://space.bilibili.com/433915419)

## License

CiProxy 遵循 MIT 开源协议