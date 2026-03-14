/**
  @creator: cilang
  @since: 2023/12/21
  @desc: 代理处理器
**/

package handler

import (
	"bufio"
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/opencvlzg/ciproxy/internal/context"
	"github.com/opencvlzg/ciproxy/internal/transfer"
	"github.com/opencvlzg/ciproxy/internal/util"
	"github.com/opencvlzg/ciproxy/pkg/mitm"
)

// DefaultOutTime 默认超时时间
var DefaultOutTime = 30e9

// DefaultWriter 默认写入器
var DefaultWriter io.Writer

// 全局拦截器实例
var (
	globalInterceptor     *mitm.Interceptor
	interceptorOnce       sync.Once
	interceptorConfig     *mitm.InterceptorConfig
	interceptorConfigLock sync.RWMutex
)

// SetInterceptorConfig 设置拦截器配置
func SetInterceptorConfig(config *mitm.InterceptorConfig) {
	interceptorConfigLock.Lock()
	defer interceptorConfigLock.Unlock()
	interceptorConfig = config
}

// GetInterceptor 获取全局拦截器实例
func GetInterceptor() *mitm.Interceptor {
	interceptorOnce.Do(func() {
		config := &mitm.InterceptorConfig{
			EnableTrafficCapture: false,
			EnableHTTP2:          true,
		}
		if interceptorConfig != nil {
			config = interceptorConfig
		}
		globalInterceptor = mitm.NewInterceptor(config)
	})
	return globalInterceptor
}

// proxyTransfer 转发流量 内部使用
func proxyTransfer(c net.Conn, s net.Conn) {
	go transfer.Transfer(c, s)
	go transfer.Transfer(s, c)
}

// proxyLogTransfer 转发流量 同时输出 内部使用
func proxyLogTransfer(c net.Conn, s net.Conn) {
	go transfer.TeeTransfer(c, s, DefaultWriter)
	go transfer.TeeTransfer(s, c, DefaultWriter)
}

// HttpProxyHandle Http处理
func HttpProxyHandle(c *context.Context) {
	buf := bufio.NewReader(c.ClientConn)
	request, err := http.ReadRequest(buf)
	if err != nil {
		return
	}
	c.ServerConn, err = net.DialTimeout("tcp", request.Host, DefaultOutTime)
	if err != nil {
		log.Println("remote host connect failed"+request.Host, err)
		return
	}
	proxyTransfer(c.ClientConn, c.ServerConn)
}

// HttpsProxyHandle Https处理
func HttpsProxyHandle(c *context.Context) {
	buf := bufio.NewReader(c.ClientConn)
	request, err := http.ReadRequest(buf)
	if err != nil {
		return
	}
	if !strings.HasSuffix(request.Host, ":443") {
		request.Host += ":443"
	}
	s, err := net.DialTimeout("tcp", request.Host, DefaultOutTime)
	if err != nil {
		log.Println("remote host connect failed"+request.Host, err)
		return
	}
	switch request.Method {
	case "CONNECT":
		_, err := c.ClientConn.Write([]byte("HTTP/1.1 200 Connection Established \r\n\r\n"))
		if err != nil {
			log.Println("write hello failed"+request.Host+request.Method, err)
			return
		}
	default:

	}
	proxyTransfer(c.ClientConn, s)
}

// HttpsSniffProxyHandle https中间人处理
func HttpsSniffProxyHandle(c *context.Context) {
	cReader := bufio.NewReader(c.ClientConn)
	request, err := http.ReadRequest(cReader)
	if err != nil {
		return
	}
	tlsCnf, err := util.GenerateTlsConfig(request.Host)
	if err != nil {
		return
	}
	if !strings.Contains(request.Host, ":443") {
		request.Host += ":443"
	}
	tlsS, err := tls.Dial("tcp", request.Host, tlsCnf)
	if err != nil {
		log.Println("remote host connect failed", err)
		return
	}
	_, err = c.ClientConn.Write([]byte("HTTP/1.1 200 Connection Established \r\n\r\n"))
	if err != nil {
		log.Println("write hello failed"+request.Host+request.Method, err)
		return
	}
	tlsC, err := upgradeTls(c.ClientConn, tlsCnf)
	if err != nil {
		log.Println("upgrade tls failed", err)
		closeConn(tlsC)
		closeConn(tlsS)
		return
	}
	proxyTransfer(tlsC, tlsS)
}

// HttpsSniffDetailProxyHandle https中间人处理
func HttpsSniffDetailProxyHandle(c *context.Context) {
	cReader := bufio.NewReader(c.ClientConn)
	request, err := http.ReadRequest(cReader)
	if err != nil {
		return
	}
	tlsCnf, err := util.GenerateTlsConfig(request.Host)
	if err != nil {
		return
	}
	if !strings.Contains(request.Host, ":443") {
		request.Host += ":443"
	}
	tlsS, err := tls.Dial("tcp", request.Host, tlsCnf)
	if err != nil {
		log.Println("remote host connect failed", err)
		return
	}
	c.TlsServerConn = tlsS
	_, err = c.ClientConn.Write([]byte("HTTP/1.1 200 Connection Established \r\n\r\n"))
	if err != nil {
		log.Println("write hello failed"+request.Host+request.Method, err)
		return
	}
	tlsC, err := upgradeTls(c.ClientConn, tlsCnf)
	if err != nil {
		log.Println("upgrade tls failed", err)
		return
	}
	c.TlsClientConn = tlsC

	go transfer.TeeDoRequestTransfer(c)
}

// TunnelProxyHandle 加密代理
func TunnelProxyHandle(c *context.Context) {
	buf := bufio.NewReader(c.ClientConn)
	request, err := http.ReadRequest(buf)
	if err != nil {
		return
	}
	if !strings.HasSuffix(request.Host, ":443") {
		request.Host += ":443"
	}
	s, err := net.DialTimeout("tcp", request.Host, DefaultOutTime)
	if err != nil {
		log.Println("remote host connect failed"+request.Host, err)
		return
	}
	switch request.Method {
	case "CONNECT":
		_, err := c.ClientConn.Write([]byte("HTTP/1.1 200 Connection Established \r\n\r\n"))
		if err != nil {
			log.Println("write hello failed"+request.Host+request.Method, err)
			return
		}
	default:

	}
	transfer.Transfer(c.ClientConn, s)
}

// isWebSocketUpgrade 判断是否为websocket链接
func isWebSocketUpgrade(req *http.Request) bool {
	upgradeHeader := req.Header.Get("Upgrade")
	connectionHeader := req.Header.Get("Connection")

	return strings.ToLower(upgradeHeader) == "websocket" && strings.Contains(strings.ToLower(connectionHeader), "upgrade")
}

// WebsocketProxyHandle websocket 代理
func WebsocketProxyHandle(c *context.Context) {
	buf := bufio.NewReader(c.ClientConn)
	request, err := http.ReadRequest(buf)
	if err != nil {
		return
	}
	if !strings.HasSuffix(request.Host, ":443") {
		request.Host += ":443"
	}
	s, err := net.DialTimeout("tcp", request.Host, DefaultOutTime)
	if err != nil {
		log.Println("remote host connect failed"+request.Host, err)
		return
	}
	if isWebSocketUpgrade(request) {
		_, err = s.Write([]byte("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\n\r\n"))
		if err != nil {
			log.Println("upgrade failed"+request.Host, err)
			return
		}
	} else {
		return
	}
	transfer.Transfer(c.ClientConn, s)
}

// closeConn 关闭连接
func closeConn(c net.Conn) {
	err := c.Close()
	if err != nil {
		log.Println("close conn failed", err)
		return
	}
}

// upgradeTls 从tcp升级到tls连接
func upgradeTls(c net.Conn, conf *tls.Config) (net.Conn, error) {
	tlsC := tls.Server(c, conf)
	err := tlsC.Handshake()
	if err != nil {
		log.Println("tls handshake failed", err)
		return nil, err
	}

	return tlsC, nil
}

// HttpInterceptProxyHandle 完整的 HTTPS MITM 拦截处理
// 支持请求/响应拦截、修改、流量捕获、HTTP/1.1 和 HTTP/2
func HttpInterceptProxyHandle(c *context.Context) {
	// 1. 读取 CONNECT 请求
	cReader := bufio.NewReader(c.ClientConn)
	request, err := http.ReadRequest(cReader)
	if err != nil {
		return
	}

	// 2. 获取拦截器
	interceptor := GetInterceptor()

	// 3. 使用 ALPN 选择器处理连接
	selector := mitm.NewALPNSelector(interceptor)
	err = selector.HandleMITMConnection(c.ClientConn, request.Host, c)
	if err != nil {
		log.Println("mitm handle error:", err)
	}
}

// AddMITMMiddleware 添加 MITM 中间件
func AddMITMMiddleware(mw *mitm.Interceptor) {
	GetInterceptor().Use(nil) // placeholder
}