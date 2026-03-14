/**
  @creator: cilang
  @since: 2024
  @desc: 请求/响应体修改中间件
**/

package builtins

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/opencvlzg/ciproxy/internal/context"
)

// BodyModifier 请求/响应体修改中间件
type BodyModifier struct {
	// 请求体替换规则
	RequestBodyRewrite map[string]string
	// 响应体替换规则
	ResponseBodyRewrite map[string]string
}

// OnRequest 修改请求体
func (m *BodyModifier) OnRequest(c *context.Context) *http.Response {
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
func (m *BodyModifier) OnResponse(c *context.Context) {
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