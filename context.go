/**
  @creator: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @github: https://github.com/OpencvLZG
  @since: 2024/1/10
  @desc: // reference from gin context design
**/

package ciproxy

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// ProxyHandle define the proxyHandle to handle proxyRequest and used by middleware
type ProxyHandle func(ctx *Context)

// ProxyHandlersChain define proxyHandle slice
type ProxyHandlersChain []ProxyHandle

// Context reference from go-proxy and gin
type Context struct {
	// ConnStatus conn status connecting closed nil
	ConnStatus string
	// ClientConn client net conn
	ClientConn net.Conn
	// TlsClientConn client tls net conn
	TlsClientConn net.Conn
	Request       *http.Request

	index    int
	handlers ProxyHandlersChain

	// protect middleware context
	mu sync.RWMutex

	// ServerConn server net conn
	ServerConn net.Conn
	// TlsServerConn server tls net conn
	TlsServerConn net.Conn
	Response      http.Response

	// Capture 流量捕获相关
	CapturedRequest  interface{} // 捕获的请求
	CapturedResponse interface{} // 捕获的响应
	CaptureID        string      // 捕获ID

	// Error 错误处理
	err           error
	errorCallback func(ctx *Context, err error)

	// Metadata 元数据
	StartTime time.Time
	Tags      []string
}

// SetClientConn set client conn
func (c *Context) SetClientConn(ClientConn net.Conn) {
	c.ClientConn = ClientConn
}

// SetRequest set the request to the context.Request
func (c *Context) SetRequest(req *http.Request) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Request = req
}

// SetServerConn set the server conn
func (c *Context) SetServerConn(ServerConn net.Conn) {
	c.ServerConn = ServerConn
}

// SetResponse set the response to the context.Response
func (c *Context) SetResponse(resp *http.Response) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Response = *resp
}

// IsAbort return abort label
func (c *Context) IsAbort() bool {
	return c.index == -1
}

// Abort set the abort index
func (c *Context) Abort() {
	c.index = -1
}

// Next set to next handle
func (c *Context) Next() {
	//c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// reset reset the context
func (c *Context) reset() {
	c.index = 0
	c.ConnStatus = "closed"
	c.ClientConn = nil
	c.TlsClientConn = nil
	c.TlsServerConn = nil
	c.Request = nil
	//c.handlers = nil
	c.ServerConn = nil
	c.Response = http.Response{}
	c.CapturedRequest = nil
	c.CapturedResponse = nil
	c.CaptureID = ""
	c.err = nil
	c.errorCallback = nil
	c.StartTime = time.Time{}
	c.Tags = c.Tags[:0]
}

// SetError 设置错误
func (c *Context) SetError(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.err = err
}

// Error 获取错误
func (c *Context) Error() error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.err
}

// SetErrorCallback 设置错误回调
func (c *Context) SetErrorCallback(cb func(ctx *Context, err error)) {
	c.errorCallback = cb
}

// HandleError 处理错误
func (c *Context) HandleError(err error) {
	c.SetError(err)
	if c.errorCallback != nil {
		c.errorCallback(c, err)
	}
}

// GetDuration 获取持续时间
func (c *Context) GetDuration() time.Duration {
	return time.Since(c.StartTime)
}

// AddTag 添加标签
func (c *Context) AddTag(tag string) {
	c.Tags = append(c.Tags, tag)
}

//func (c *Context)GetTransport()    *http.Transport{
//	return &http.Transport{DialTLS: c.TlsServerConn.RemoteAddr().String()}
//}
