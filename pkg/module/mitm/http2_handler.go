/**
  @creator: cilang
  @since: 2024
  @desc: HTTP/2 协议处理器
**/

package mitm

import (
	"bytes"
	"golang.org/x/net/http2"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

// HTTP2Handler HTTP/2 协议处理器
type HTTP2Handler struct {
	interceptor *Interceptor
	server      *http2.Server
}

// NewHTTP2Handler 创建新的 HTTP/2 处理器
func NewHTTP2Handler(interceptor *Interceptor) *HTTP2Handler {
	return &HTTP2Handler{
		interceptor: interceptor,
		server: &http2.Server{
			IdleTimeout: 30 * time.Second,
		},
	}
}

// HandleConnection 处理 HTTP/2 连接
func (h *HTTP2Handler) HandleConnection(clientConn, serverConn net.Conn, ctx Context) error {
	// 创建到目标服务器的 HTTP/2 客户端传输
	clientTransport := &http2.Transport{
		AllowHTTP: false,
	}

	// 使用连接池管理服务器连接
	connPool := &http2ConnPool{
		conn: serverConn,
	}

	// HTTP/2 服务器处理
	opts := &http2.ServeConnOpts{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 创建简单的上下文用于每个请求
			reqCtx := &simpleContext{
				request:   r,
				startTime: time.Now(),
				tags:      make([]string, 0),
			}

			// 捕获请求
			h.interceptor.CaptureRequest(reqCtx)

			// 请求拦截
			if blockResp := h.interceptor.HandleRequest(reqCtx); blockResp != nil {
				// 阻断请求
				h.copyResponse(w, blockResp)
				return
			}

			// 转发请求到目标服务器
			resp, err := h.forwardHTTP2Request(clientTransport, connPool, r)
			if err != nil {
				log.Println("forward http2 request error:", err)
				http.Error(w, "Bad Gateway", http.StatusBadGateway)
				return
			}
			defer resp.Body.Close()

			reqCtx.SetResponse(resp)

			// 响应拦截
			h.interceptor.HandleResponse(reqCtx)

			// 捕获响应
			h.interceptor.CaptureResponse(reqCtx)

			// 复制响应头
			for k, v := range resp.Header {
				w.Header()[k] = v
			}

			// 设置状态码
			w.WriteHeader(resp.StatusCode)

			// 复制响应体
			io.Copy(w, resp.Body)
		}),
	}

	h.server.ServeConn(clientConn, opts)
	return nil
}

// forwardHTTP2Request 转发 HTTP/2 请求
func (h *HTTP2Handler) forwardHTTP2Request(transport *http2.Transport, pool *http2ConnPool, r *http.Request) (*http.Response, error) {
	// 使用连接池获取连接
	conn := pool.Get()
	if conn == nil {
		return nil, io.ErrUnexpectedEOF
	}

	// 创建转发请求
	forwardReq := &http.Request{
		Method:        r.Method,
		URL:           r.URL,
		Host:          r.Host,
		Header:        r.Header,
		Body:          r.Body,
		ContentLength: r.ContentLength,
	}

	// 发送请求
	return transport.RoundTrip(forwardReq)
}

// copyResponse 复制响应
func (h *HTTP2Handler) copyResponse(w http.ResponseWriter, resp *http.Response) {
	// 复制响应头
	for k, v := range resp.Header {
		w.Header()[k] = v
	}

	// 设置状态码
	w.WriteHeader(resp.StatusCode)

	// 复制响应体
	if resp.Body != nil {
		defer resp.Body.Close()
		io.Copy(w, resp.Body)
	}
}

// http2ConnPool HTTP/2 连接池
type http2ConnPool struct {
	mu   sync.Mutex
	conn net.Conn
}

// Get 获取连接
func (p *http2ConnPool) Get() net.Conn {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.conn
}

// Put 放回连接
func (p *http2ConnPool) Put(conn net.Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.conn = conn
}

// Write 实现 io.Writer 接口
func (p *http2ConnPool) Write(data []byte) (int, error) {
	conn := p.Get()
	if conn == nil {
		return 0, io.ErrUnexpectedEOF
	}
	return conn.Write(data)
}

// Read 实现 io.Reader 接口
func (p *http2ConnPool) Read(data []byte) (int, error) {
	conn := p.Get()
	if conn == nil {
		return 0, io.ErrUnexpectedEOF
	}
	return conn.Read(data)
}

// simpleContext 简单的上下文实现，用于 HTTP/2 请求处理
type simpleContext struct {
	request          *http.Request
	response         http.Response
	startTime        time.Time
	tags             []string
	capturedRequest  interface{}
	capturedResponse interface{}
	captureID        string
	aborted          bool
}

func (c *simpleContext) GetRequest() *http.Request {
	return c.request
}

func (c *simpleContext) SetRequest(req *http.Request) {
	c.request = req
}

func (c *simpleContext) GetRequestBody() []byte {
	if c.request == nil || c.request.Body == nil {
		return nil
	}
	body, _ := io.ReadAll(c.request.Body)
	c.request.Body = io.NopCloser(bytes.NewReader(body))
	return body
}

func (c *simpleContext) SetRequestBody(body []byte) {
	if c.request != nil {
		c.request.Body = io.NopCloser(bytes.NewReader(body))
		c.request.ContentLength = int64(len(body))
	}
}

func (c *simpleContext) GetRequestHeader(key string) string {
	if c.request == nil {
		return ""
	}
	return c.request.Header.Get(key)
}

func (c *simpleContext) SetRequestHeader(key, value string) {
	if c.request != nil {
		c.request.Header.Set(key, value)
	}
}

func (c *simpleContext) DelRequestHeader(key string) {
	if c.request != nil {
		c.request.Header.Del(key)
	}
}

func (c *simpleContext) GetResponse() *http.Response {
	return &c.response
}

func (c *simpleContext) SetResponse(resp *http.Response) {
	if resp != nil {
		c.response = *resp
	}
}

func (c *simpleContext) GetResponseBody() []byte {
	if c.response.Body == nil {
		return nil
	}
	body, _ := io.ReadAll(c.response.Body)
	c.response.Body = io.NopCloser(bytes.NewReader(body))
	return body
}

func (c *simpleContext) SetResponseBody(body []byte) {
	c.response.Body = io.NopCloser(bytes.NewReader(body))
	c.response.ContentLength = int64(len(body))
}

func (c *simpleContext) GetResponseHeader(key string) string {
	return c.response.Header.Get(key)
}

func (c *simpleContext) SetResponseHeader(key, value string) {
	c.response.Header.Set(key, value)
}

func (c *simpleContext) DelResponseHeader(key string) {
	c.response.Header.Del(key)
}

func (c *simpleContext) GetCapturedRequest() interface{} {
	return c.capturedRequest
}

func (c *simpleContext) SetCapturedRequest(req interface{}) {
	c.capturedRequest = req
}

func (c *simpleContext) GetCapturedResponse() interface{} {
	return c.capturedResponse
}

func (c *simpleContext) SetCapturedResponse(resp interface{}) {
	c.capturedResponse = resp
}

func (c *simpleContext) GetCaptureID() string {
	return c.captureID
}

func (c *simpleContext) SetCaptureID(id string) {
	c.captureID = id
}

func (c *simpleContext) GetStartTime() time.Time {
	return c.startTime
}

func (c *simpleContext) SetStartTime(t time.Time) {
	c.startTime = t
}

func (c *simpleContext) AddTag(tag string) {
	c.tags = append(c.tags, tag)
}

func (c *simpleContext) GetTags() []string {
	return c.tags
}

func (c *simpleContext) Abort() {
	c.aborted = true
}

func (c *simpleContext) IsAborted() bool {
	return c.aborted
}
