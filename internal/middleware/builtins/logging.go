/**
  @creator: cilang
  @since: 2024
  @desc: 日志记录中间件
**/

package builtins

import (
	"log"
	"net/http"

	"github.com/opencvlzg/ciproxy/internal/context"
)

// LoggingMiddleware 日志记录中间件
type LoggingMiddleware struct {
	Logger func(format string, args ...interface{})
}

// OnRequest 记录请求日志
func (m *LoggingMiddleware) OnRequest(c *context.Context) *http.Response {
	if m.Logger == nil {
		m.Logger = log.Printf
	}

	if c.Request != nil {
		m.Logger("[REQUEST] %s %s", c.Request.Method, c.Request.URL.String())
	}

	return nil
}

// OnResponse 记录响应日志
func (m *LoggingMiddleware) OnResponse(c *context.Context) {
	if m.Logger == nil {
		m.Logger = log.Printf
	}

	m.Logger("[RESPONSE] %d", c.Response.StatusCode)
}