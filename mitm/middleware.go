/**
  @creator: cilang
  @since: 2024
  @desc: 中间件接口和常用中间件实现
**/

package mitm

import (
	"bytes"
	"github.com/opencvlzg/ciproxy"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// BaseMiddleware 基础中间件，提供默认实现
type BaseMiddleware struct {
	RequestHandler  func(c *ciproxy.Context) *http.Response
	ResponseHandler func(c *ciproxy.Context)
}

// OnRequest 请求处理
func (m *BaseMiddleware) OnRequest(c *ciproxy.Context) *http.Response {
	if m.RequestHandler != nil {
		return m.RequestHandler(c)
	}
	return nil
}

// OnResponse 响应处理
func (m *BaseMiddleware) OnResponse(c *ciproxy.Context) {
	if m.ResponseHandler != nil {
		m.ResponseHandler(c)
	}
}

// HeaderModifier 请求头修改中间件
type HeaderModifier struct {
	// 要添加/修改的请求头
	RequestHeaders map[string]string
	// 要删除的请求头
	DelRequestHeaders []string
	// 要添加/修改的响应头
	ResponseHeaders map[string]string
	// 要删除的响应头
	DelResponseHeaders []string
}

// OnRequest 修改请求头
func (m *HeaderModifier) OnRequest(c *ciproxy.Context) *http.Response {
	if c.Request == nil {
		return nil
	}

	// 删除请求头
	for _, key := range m.DelRequestHeaders {
		c.Request.Header.Del(key)
	}

	// 添加/修改请求头
	for key, value := range m.RequestHeaders {
		c.Request.Header.Set(key, value)
	}

	return nil
}

// OnResponse 修改响应头
func (m *HeaderModifier) OnResponse(c *ciproxy.Context) {
	// 删除响应头
	for _, key := range m.DelResponseHeaders {
		c.Response.Header.Del(key)
	}

	// 添加/修改响应头
	for key, value := range m.ResponseHeaders {
		c.Response.Header.Set(key, value)
	}
}

// BodyModifier 请求/响应体修改中间件
type BodyModifier struct {
	// 请求体替换规则
	RequestBodyRewrite map[string]string
	// 响应体替换规则
	ResponseBodyRewrite map[string]string
}

// OnRequest 修改请求体
func (m *BodyModifier) OnRequest(c *ciproxy.Context) *http.Response {
	if c.Request == nil || c.Request.Body == nil {
		return nil
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil
	}
	c.Request.Body.Close()

	// 执行替换
	modified := string(body)
	for old, new := range m.RequestBodyRewrite {
		modified = strings.ReplaceAll(modified, old, new)
	}

	c.Request.Body = io.NopCloser(bytes.NewReader([]byte(modified)))
	c.Request.ContentLength = int64(len(modified))

	return nil
}

// OnResponse 修改响应体
func (m *BodyModifier) OnResponse(c *ciproxy.Context) {
	if c.Response.Body == nil {
		return
	}

	body, err := io.ReadAll(c.Response.Body)
	if err != nil {
		return
	}
	c.Response.Body.Close()

	// 执行替换
	modified := string(body)
	for old, new := range m.ResponseBodyRewrite {
		modified = strings.ReplaceAll(modified, old, new)
	}

	c.Response.Body = io.NopCloser(bytes.NewReader([]byte(modified)))
	c.Response.ContentLength = int64(len(modified))
}

// BlockMiddleware 请求阻断中间件
type BlockMiddleware struct {
	// URL 匹配模式（正则表达式）
	URLPatterns []string
	// 阻断时的响应状态码
	StatusCode int
	// 阻断时的响应体
	ResponseBody string
	// 阻断时的响应头
	ResponseHeaders map[string]string

	compiledPatterns []*regexp.Regexp
}

// OnRequest 检查并阻断请求
func (m *BlockMiddleware) OnRequest(c *ciproxy.Context) *http.Response {
	if c.Request == nil {
		return nil
	}

	// 编译正则表达式
	if m.compiledPatterns == nil {
		m.compiledPatterns = make([]*regexp.Regexp, 0, len(m.URLPatterns))
		for _, pattern := range m.URLPatterns {
			re, err := regexp.Compile(pattern)
			if err != nil {
				log.Printf("invalid pattern %s: %v", pattern, err)
				continue
			}
			m.compiledPatterns = append(m.compiledPatterns, re)
		}
	}

	// 检查 URL 是否匹配
	url := c.Request.URL.String()
	for _, re := range m.compiledPatterns {
		if re.MatchString(url) {
			// 创建阻断响应
			resp := &http.Response{
				StatusCode: m.StatusCode,
				Header:     make(http.Header),
				Body:       io.NopCloser(bytes.NewReader([]byte(m.ResponseBody))),
			}

			// 设置响应头
			for key, value := range m.ResponseHeaders {
				resp.Header.Set(key, value)
			}

			return resp
		}
	}

	return nil
}

// OnResponse 空实现
func (m *BlockMiddleware) OnResponse(c *ciproxy.Context) {}

// LoggingMiddleware 日志记录中间件
type LoggingMiddleware struct {
	Logger func(format string, args ...interface{})
}

// OnRequest 记录请求日志
func (m *LoggingMiddleware) OnRequest(c *ciproxy.Context) *http.Response {
	if m.Logger == nil {
		m.Logger = log.Printf
	}

	if c.Request != nil {
		m.Logger("[REQUEST] %s %s", c.Request.Method, c.Request.URL.String())
	}

	return nil
}

// OnResponse 记录响应日志
func (m *LoggingMiddleware) OnResponse(c *ciproxy.Context) {
	if m.Logger == nil {
		m.Logger = log.Printf
	}

	m.Logger("[RESPONSE] %d", c.Response.StatusCode)
}

// CorsMiddleware CORS 跨域中间件
type CorsMiddleware struct {
	// 允许的来源
	AllowOrigins []string
	// 允许的方法
	AllowMethods []string
	// 允许的头
	AllowHeaders []string
	// 是否允许凭证
	AllowCredentials bool
}

// OnRequest 处理预检请求
func (m *CorsMiddleware) OnRequest(c *ciproxy.Context) *http.Response {
	if c.Request == nil {
		return nil
	}

	// 处理 OPTIONS 预检请求
	if c.Request.Method == "OPTIONS" {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		}

		m.setCorsHeaders(resp)

		return resp
	}

	return nil
}

// OnResponse 添加 CORS 头
func (m *CorsMiddleware) OnResponse(c *ciproxy.Context) {
	m.setCorsHeaders(&c.Response)
}

// setCorsHeaders 设置 CORS 头
func (m *CorsMiddleware) setCorsHeaders(resp *http.Response) {
	if len(m.AllowOrigins) > 0 {
		resp.Header.Set("Access-Control-Allow-Origin", strings.Join(m.AllowOrigins, ", "))
	}
	if len(m.AllowMethods) > 0 {
		resp.Header.Set("Access-Control-Allow-Methods", strings.Join(m.AllowMethods, ", "))
	}
	if len(m.AllowHeaders) > 0 {
		resp.Header.Set("Access-Control-Allow-Headers", strings.Join(m.AllowHeaders, ", "))
	}
	if m.AllowCredentials {
		resp.Header.Set("Access-Control-Allow-Credentials", "true")
	}
}
