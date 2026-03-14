/**
  @creator: cilang
  @since: 2024
  @desc: MITM 拦截器核心模块
**/

package mitm

import (
	"crypto/tls"
	"github.com/opencvlzg/ciproxy/traffic"
	"net"
	"net/http"
	"sync"
	"time"
)

// Context MITM 上下文接口
// 避免循环导入，由 ciproxy.Context 实现
type Context interface {
	// 请求相关
	GetRequest() *http.Request
	SetRequest(req *http.Request)
	GetRequestBody() []byte
	SetRequestBody(body []byte)
	GetRequestHeader(key string) string
	SetRequestHeader(key, value string)
	DelRequestHeader(key string)

	// 响应相关
	GetResponse() *http.Response
	SetResponse(resp *http.Response)
	GetResponseBody() []byte
	SetResponseBody(body []byte)
	GetResponseHeader(key string) string
	SetResponseHeader(key, value string)
	DelResponseHeader(key string)

	// 流量捕获
	GetCapturedRequest() interface{}
	SetCapturedRequest(req interface{})
	GetCapturedResponse() interface{}
	SetCapturedResponse(resp interface{})
	GetCaptureID() string
	SetCaptureID(id string)

	// 时间相关
	GetStartTime() time.Time
	SetStartTime(t time.Time)

	// 标签
	AddTag(tag string)
	GetTags() []string

	// 中断
	Abort()
	IsAborted() bool
}

// Middleware 中间件接口
type Middleware interface {
	// OnRequest 请求拦截，返回 nil 表示放行，返回 response 表示阻断
	OnRequest(c Context) *http.Response
	// OnResponse 响应拦截
	OnResponse(c Context)
}

// MiddlewareFunc 中间件函数类型，用于将函数转换为 Middleware
type MiddlewareFunc func(c Context) *http.Response

// OnRequest 实现 Middleware 接口
func (mf MiddlewareFunc) OnRequest(c Context) *http.Response {
	return mf(c)
}

// OnResponse 实现 Middleware 接口
func (mf MiddlewareFunc) OnResponse(c Context) {}

// ResponseMiddleware 响应中间件函数类型
type ResponseMiddleware func(c Context)

// Interceptor MITM 拦截器
type Interceptor struct {
	mu           sync.RWMutex
	middlewares  []Middleware
	trafficStore *traffic.TrafficStore
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
		middlewares:          make([]Middleware, 0),
		trafficStore:         traffic.GetTrafficStore(),
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
func (i *Interceptor) Use(middleware Middleware) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.middlewares = append(i.middlewares, middleware)
}

// UseFunc 添加中间件函数
func (i *Interceptor) UseFunc(fn func(c Context) *http.Response) {
	i.Use(MiddlewareFunc(fn))
}

// GetMiddlewares 获取中间件列表
func (i *Interceptor) GetMiddlewares() []Middleware {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.middlewares
}

// HandleRequest 处理请求拦截
func (i *Interceptor) HandleRequest(c Context) *http.Response {
	for _, mw := range i.GetMiddlewares() {
		if resp := mw.OnRequest(c); resp != nil {
			return resp
		}
	}
	return nil
}

// HandleResponse 处理响应拦截
func (i *Interceptor) HandleResponse(c Context) {
	for _, mw := range i.GetMiddlewares() {
		mw.OnResponse(c)
	}
}

// CaptureRequest 捕获请求
func (i *Interceptor) CaptureRequest(c Context) string {
	if !i.EnableTrafficCapture || c.GetRequest() == nil {
		return ""
	}

	captured := traffic.CaptureFromRequest(c.GetRequest())
	c.SetCapturedRequest(captured)
	c.SetCaptureID(captured.ID)

	if i.trafficStore != nil {
		return i.trafficStore.Capture(captured)
	}
	return captured.ID
}

// CaptureResponse 捕获响应
func (i *Interceptor) CaptureResponse(c Context) {
	if !i.EnableTrafficCapture || c.GetCaptureID() == "" {
		return
	}

	resp := c.GetResponse()
	if resp == nil {
		return
	}

	captured := &traffic.CapturedResponse{
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
func (i *Interceptor) GetTrafficStore() *traffic.TrafficStore {
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
func (i *Interceptor) HandleConnection(clientConn, serverTLS net.Conn, proto string, ctx Context) error {
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
