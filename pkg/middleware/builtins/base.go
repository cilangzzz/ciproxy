/**
  @creator: cilang
  @since: 2024
  @desc: 基础中间件实现
**/

package builtins

import (
	"net/http"

	"github.com/opencvlzg/ciproxy/pkg/context"
)

// BaseMiddleware 基础中间件，提供默认实现
type BaseMiddleware struct {
	RequestHandler  func(c *context.Context) *http.Response
	ResponseHandler func(c *context.Context)
}

// OnRequest 请求处理
func (m *BaseMiddleware) OnRequest(c *context.Context) *http.Response {
	if m.RequestHandler != nil {
		return m.RequestHandler(c)
	}
	return nil
}

// OnResponse 响应处理
func (m *BaseMiddleware) OnResponse(c *context.Context) {
	if m.ResponseHandler != nil {
		m.ResponseHandler(c)
	}
}