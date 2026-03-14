/**
  @creator: cilang
  @since: 2024
  @desc: 请求阻断中间件
**/

package builtins

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/opencvlzg/ciproxy/internal/context"
)

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
func (m *BlockMiddleware) OnRequest(c *context.Context) *http.Response {
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
func (m *BlockMiddleware) OnResponse(c *context.Context) {}
