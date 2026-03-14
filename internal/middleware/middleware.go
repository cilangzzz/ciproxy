/**
  @creator: cilang
  @since: 2024
  @desc: 中间件接口定义
**/

package middleware

import (
	"net/http"

	"github.com/opencvlzg/ciproxy/internal/context"
)

// Middleware 中间件接口
type Middleware interface {
	// OnRequest 请求拦截，返回 nil 表示放行，返回 response 表示阻断
	OnRequest(c *context.Context) *http.Response
	// OnResponse 响应拦截
	OnResponse(c *context.Context)
}

// MiddlewareFunc 中间件函数类型，用于将函数转换为 Middleware
type MiddlewareFunc func(c *context.Context) *http.Response

// OnRequest 实现 Middleware 接口
func (mf MiddlewareFunc) OnRequest(c *context.Context) *http.Response {
	return mf(c)
}

// OnResponse 实现 Middleware 接口
func (mf MiddlewareFunc) OnResponse(c *context.Context) {}