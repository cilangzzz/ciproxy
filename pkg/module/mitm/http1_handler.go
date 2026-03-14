/**
  @creator: cilang
  @since: 2024
  @desc: HTTP/1.1 协议处理器
**/

package mitm

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HTTP1Handler HTTP/1.1 协议处理器
type HTTP1Handler struct {
	interceptor *Interceptor
}

// NewHTTP1Handler 创建新的 HTTP/1.1 处理器
func NewHTTP1Handler(interceptor *Interceptor) *HTTP1Handler {
	return &HTTP1Handler{
		interceptor: interceptor,
	}
}

// HandleConnection 处理 HTTP/1.1 连接
func (h *HTTP1Handler) HandleConnection(clientConn, serverConn net.Conn, ctx Context) error {
	reader := bufio.NewReader(clientConn)

	// 支持 keep-alive 循环处理
	for {
		ctx.SetStartTime(time.Now())

		// 1. 解析请求
		req, err := http.ReadRequest(reader)
		if err != nil {
			if err != io.EOF {
				log.Println("read request error:", err)
			}
			return err
		}

		ctx.SetRequest(req)

		// 2. 捕获请求
		h.interceptor.CaptureRequest(ctx)

		// 3. 请求拦截
		if blockResp := h.interceptor.HandleRequest(ctx); blockResp != nil {
			// 阻断请求，直接返回响应
			return h.writeResponse(clientConn, blockResp)
		}

		// 4. 转发请求到目标服务器
		resp, err := h.forwardRequest(serverConn, ctx)
		if err != nil {
			log.Println("forward request error:", err)
			return err
		}

		ctx.SetResponse(resp)

		// 5. 响应拦截
		h.interceptor.HandleResponse(ctx)

		// 6. 捕获响应
		h.interceptor.CaptureResponse(ctx)

		// 7. 返回响应给客户端
		if err := h.writeResponse(clientConn, ctx.GetResponse()); err != nil {
			log.Println("write response error:", err)
			return err
		}

		// 8. 检查是否 keep-alive
		if !h.shouldKeepAlive(req, ctx.GetResponse()) {
			break
		}
	}

	return nil
}

// forwardRequest 转发请求到目标服务器
func (h *HTTP1Handler) forwardRequest(serverConn net.Conn, ctx Context) (*http.Response, error) {
	req := ctx.GetRequest()

	// 构建完整 URL
	scheme := "https"
	host := req.Host
	if !strings.Contains(host, ":") {
		host += ":443"
	}

	// 创建转发请求
	u, _ := url.Parse(scheme + "://" + host + req.URL.String())

	forwardReq := &http.Request{
		Method:        req.Method,
		URL:           u,
		Host:          req.Host,
		Header:        req.Header,
		Body:          req.Body,
		ContentLength: req.ContentLength,
		Close:         req.Close,
	}

	// 写入请求到服务器连接
	var buf bytes.Buffer
	forwardReq.Write(&buf)
	_, err := serverConn.Write(buf.Bytes())
	if err != nil {
		return nil, err
	}

	// 读取响应
	reader := bufio.NewReader(serverConn)
	resp, err := http.ReadResponse(reader, forwardReq)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// writeResponse 写入响应到客户端
func (h *HTTP1Handler) writeResponse(clientConn net.Conn, resp *http.Response) error {
	var buf bytes.Buffer
	resp.Write(&buf)
	_, err := clientConn.Write(buf.Bytes())
	return err
}

// shouldKeepAlive 检查是否保持连接
func (h *HTTP1Handler) shouldKeepAlive(req *http.Request, resp *http.Response) bool {
	// 检查请求 Connection 头
	reqConn := req.Header.Get("Connection")
	if strings.ToLower(reqConn) == "close" {
		return false
	}

	// 检查响应 Connection 头
	respConn := resp.Header.Get("Connection")
	if strings.ToLower(respConn) == "close" {
		return false
	}

	// HTTP/1.1 默认 keep-alive
	return req.ProtoAtLeast(1, 1)
}
