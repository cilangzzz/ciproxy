/**
  @creator: cilang
  @since: 2024
  @desc: CORS 跨域中间件
**/

package builtins

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/opencvlzg/ciproxy/pkg/context"
)

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
func (m *CorsMiddleware) OnRequest(c *context.Context) *http.Response {
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
func (m *CorsMiddleware) OnResponse(c *context.Context) {
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