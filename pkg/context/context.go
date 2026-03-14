/**
  @creator: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @github: https://github.com/OpencvLZG
  @since: 2024/1/10
  @desc: // reference from gin context design
**/

package context

import (
	"bytes"
	"io"
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
	Handlers ProxyHandlersChain

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
	s := len(c.Handlers)
	for ; c.index < s; c.index++ {
		c.Handlers[c.index](c)
	}
}

// Reset reset the context
func (c *Context) Reset() {
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

// GetRequestBody 获取请求体
func (c *Context) GetRequestBody() []byte {
	if c.Request == nil || c.Request.Body == nil {
		return nil
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil
	}
	c.Request.Body.Close()
	c.Request.Body = io.NopCloser(bytes.NewReader(body))
	return body
}

// SetRequestBody 设置请求体
func (c *Context) SetRequestBody(body []byte) {
	if c.Request == nil {
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))
	c.Request.ContentLength = int64(len(body))
}

// GetResponseBody 获取响应体
func (c *Context) GetResponseBody() []byte {
	if c.Response.Body == nil {
		return nil
	}
	body, err := io.ReadAll(c.Response.Body)
	if err != nil {
		return nil
	}
	c.Response.Body.Close()
	c.Response.Body = io.NopCloser(bytes.NewReader(body))
	return body
}

// SetResponseBody 设置响应体
func (c *Context) SetResponseBody(body []byte) {
	c.Response.Body = io.NopCloser(bytes.NewReader(body))
	c.Response.ContentLength = int64(len(body))
}

// GetRequestHeader 获取请求头
func (c *Context) GetRequestHeader(key string) string {
	if c.Request == nil {
		return ""
	}
	return c.Request.Header.Get(key)
}

// SetRequestHeader 设置请求头
func (c *Context) SetRequestHeader(key, value string) {
	if c.Request == nil {
		return
	}
	c.Request.Header.Set(key, value)
}

// DelRequestHeader 删除请求头
func (c *Context) DelRequestHeader(key string) {
	if c.Request == nil {
		return
	}
	c.Request.Header.Del(key)
}

// GetResponseHeader 获取响应头
func (c *Context) GetResponseHeader(key string) string {
	return c.Response.Header.Get(key)
}

// SetResponseHeader 设置响应头
func (c *Context) SetResponseHeader(key, value string) {
	c.Response.Header.Set(key, value)
}

// DelResponseHeader 删除响应头
func (c *Context) DelResponseHeader(key string) {
	c.Response.Header.Del(key)
}

// BlockWithResponse 阻断请求并返回自定义响应
func (c *Context) BlockWithResponse(resp *http.Response) {
	c.Response = *resp
	c.Abort()
}

// BlockWithStatus 阻断请求并返回简单响应
func (c *Context) BlockWithStatus(code int, body string) {
	c.Response = http.Response{
		StatusCode:    code,
		Header:        make(http.Header),
		Body:          io.NopCloser(bytes.NewReader([]byte(body))),
		ContentLength: int64(len(body)),
	}
	c.Abort()
}

// GetRequest 获取请求
func (c *Context) GetRequest() *http.Request {
	return c.Request
}

// GetResponse 获取响应
func (c *Context) GetResponse() *http.Response {
	return &c.Response
}

// GetCapturedRequest 获取捕获的请求
func (c *Context) GetCapturedRequest() interface{} {
	return c.CapturedRequest
}

// SetCapturedRequest 设置捕获的请求
func (c *Context) SetCapturedRequest(req interface{}) {
	c.CapturedRequest = req
}

// GetCapturedResponse 获取捕获的响应
func (c *Context) GetCapturedResponse() interface{} {
	return c.CapturedResponse
}

// SetCapturedResponse 设置捕获的响应
func (c *Context) SetCapturedResponse(resp interface{}) {
	c.CapturedResponse = resp
}

// GetCaptureID 获取捕获ID
func (c *Context) GetCaptureID() string {
	return c.CaptureID
}

// SetCaptureID 设置捕获ID
func (c *Context) SetCaptureID(id string) {
	c.CaptureID = id
}

// GetStartTime 获取开始时间
func (c *Context) GetStartTime() time.Time {
	return c.StartTime
}

// SetStartTime 设置开始时间
func (c *Context) SetStartTime(t time.Time) {
	c.StartTime = t
}

// GetTags 获取标签
func (c *Context) GetTags() []string {
	return c.Tags
}

// IsAborted 是否已中断
func (c *Context) IsAborted() bool {
	return c.index == -1
}