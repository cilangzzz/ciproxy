/**
  @creator: cilang
  @since: 2024
  @desc: 请求头修改中间件
**/

package builtins

import (
	"net/http"

	"github.com/opencvlzg/ciproxy/pkg/context"
)

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
func (m *HeaderModifier) OnRequest(c *context.Context) *http.Response {
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
func (m *HeaderModifier) OnResponse(c *context.Context) {
	// 删除响应头
	for _, key := range m.DelResponseHeaders {
		c.Response.Header.Del(key)
	}

	// 添加/修改响应头
	for key, value := range m.ResponseHeaders {
		c.Response.Header.Set(key, value)
	}
}