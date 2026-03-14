/**
  @creator: cilang
  @since: 2024
  @desc: MITM 拦截器核心模块
**/

package mitm

import (
	"crypto/tls"
	"github.com/opencvlzg/ciproxy/internal/transfer"
	"net"
	"net/http"
	"sync"

	"github.com/opencvlzg/ciproxy/internal/context"
	"github.com/opencvlzg/ciproxy/internal/middleware"
)

// Interceptor MITM 拦截器
type Interceptor struct {
	mu           sync.RWMutex
	middlewares  []middleware.Middleware
	trafficStore *transfer.TrafficStore
	http1Handler *HTTP1Handler
	http2Handler *HTTP2Handler
	alpnSelector *ALPNSelector

	// 配置选项
	EnableTrafficCapture bool
	EnableHTTP2          bool
}

// InterceptorConfig 拦截器配置
type InterceptorConfig struct {
	EnableTrafficCapture bool
	EnableHTTP2          bool
}

// NewInterceptor 创建新的拦截器
func NewInterceptor(config *InterceptorConfig) *Interceptor {
	i := &Interceptor{
		middlewares:          make([]middleware.Middleware, 0),
		trafficStore:         transfer.GetTrafficStore(),
		EnableTrafficCapture: config.EnableTrafficCapture,
		EnableHTTP2:          config.EnableHTTP2,
	}

	// 创建协议处理器
	i.http1Handler = NewHTTP1Handler(i)
	i.http2Handler = NewHTTP2Handler(i)
	i.alpnSelector = NewALPNSelector(i)

	return i
}

// Use 添加中间件
func (i *Interceptor) Use(m middleware.Middleware) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.middlewares = append(i.middlewares, m)
}

// UseFunc 添加中间件函数
func (i *Interceptor) UseFunc(fn func(c *context.Context) *http.Response) {
	i.Use(middleware.MiddlewareFunc(fn))
}

// GetMiddlewares 获取中间件列表
func (i *Interceptor) GetMiddlewares() []middleware.Middleware {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.middlewares
}

// HandleRequest 处理请求拦截
func (i *Interceptor) HandleRequest(c *context.Context) *http.Response {
	for _, mw := range i.GetMiddlewares() {
		if resp := mw.OnRequest(c); resp != nil {
			return resp
		}
	}
	return nil
}

// HandleResponse 处理响应拦截
func (i *Interceptor) HandleResponse(c *context.Context) {
	for _, mw := range i.GetMiddlewares() {
		mw.OnResponse(c)
	}
}

// CaptureRequest 捕获请求
func (i *Interceptor) CaptureRequest(c *context.Context) string {
	if !i.EnableTrafficCapture || c.GetRequest() == nil {
		return ""
	}

	captured := transfer.CaptureFromRequest(c.GetRequest())
	c.SetCapturedRequest(captured)
	c.SetCaptureID(captured.ID)

	if i.trafficStore != nil {
		return i.trafficStore.Capture(captured)
	}
	return captured.ID
}

// CaptureResponse 捕获响应
func (i *Interceptor) CaptureResponse(c *context.Context) {
	if !i.EnableTrafficCapture || c.GetCaptureID() == "" {
		return
	}

	resp := c.GetResponse()
	if resp == nil {
		return
	}

	captured := &transfer.CapturedResponse{
		ID:        c.GetCaptureID(),
		Timestamp: c.GetStartTime(),
		Status:    resp.StatusCode,
		Headers:   make(map[string]string),
	}

	for k, v := range resp.Header {
		if len(v) > 0 {
			captured.Headers[k] = v[0]
		}
	}

	c.SetCapturedResponse(captured)

	if i.trafficStore != nil {
		i.trafficStore.SetResponse(c.GetCaptureID(), captured)
	}
}

// GetTrafficStore 获取流量存储
func (i *Interceptor) GetTrafficStore() *transfer.TrafficStore {
	return i.trafficStore
}

// HandleConnect 处理 CONNECT 请求，建立 TLS 连接
func (i *Interceptor) HandleConnect(clientConn net.Conn, host string) (clientTLS, serverTLS net.Conn, proto string, err error) {
	// 生成 TLS 配置
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"h2", "http/1.1"},
	}

	// 连接目标服务器
	serverTLS, err = tls.Dial("tcp", host, tlsConfig)
	if err != nil {
		return nil, nil, "", err
	}

	// 获取协商的协议
	proto = serverTLS.ConnectionState().NegotiatedProtocol
	if proto == "" {
		proto = "http/1.1"
	}

	// 发送 200 Connection Established
	_, err = clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	if err != nil {
		serverTLS.Close()
		return nil, nil, "", err
	}

	return clientConn, serverTLS, proto, nil
}

// HandleConnection 根据协议处理连接
func (i *Interceptor) HandleConnection(clientConn, serverTLS net.Conn, proto string, ctx *context.Context) error {
	switch proto {
	case "h2":
		if i.EnableHTTP2 && i.http2Handler != nil {
			return i.http2Handler.HandleConnection(clientConn, serverTLS, ctx)
		}
		fallthrough
	default:
		return i.http1Handler.HandleConnection(clientConn, serverTLS, ctx)
	}
}
