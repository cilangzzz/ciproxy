/**
  @creator: cilang
  @since: 2024
  @desc: 流量捕获核心逻辑
**/

package transfer

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

// Captor 流量捕获器接口
type Captor interface {
	CaptureFromRequest(req *http.Request) (*CapturedRequest, error)
	CaptureFromResponse(resp *http.Response, id string) (*CapturedResponse, error)
}

// DefaultCaptor 默认捕获器
type DefaultCaptor struct {
	idGenerator IDGenerator
}

// NewCaptor 创建捕获器
func NewCaptor() *DefaultCaptor {
	return &DefaultCaptor{
		idGenerator: defaultIDGenerator,
	}
}

// NewCaptorWithGenerator 创建带自定义ID生成器的捕获器
func NewCaptorWithGenerator(g IDGenerator) *DefaultCaptor {
	return &DefaultCaptor{
		idGenerator: g,
	}
}

// CaptureFromRequest 从 http.Request 捕获
func (c *DefaultCaptor) CaptureFromRequest(req *http.Request) (*CapturedRequest, error) {
	headers := make(map[string]string)
	for k, v := range req.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	var body []byte
	if req.Body != nil {
		var err error
		body, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		// 重新包装 Body 使其可重复读取
		req.Body = io.NopCloser(bytes.NewReader(body))
	}

	return &CapturedRequest{
		ID:        c.idGenerator.Generate(),
		Timestamp: time.Now(),
		Method:    req.Method,
		URL:       req.URL.String(),
		Host:      req.Host,
		Headers:   headers,
		Body:      body,
		Protocol:  req.Proto,
	}, nil
}

// CaptureFromResponse 从 http.Response 捕获
func (c *DefaultCaptor) CaptureFromResponse(resp *http.Response, id string) (*CapturedResponse, error) {
	headers := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	var body []byte
	if resp.Body != nil {
		var err error
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		// 重新包装 Body 使其可重复读取
		resp.Body = io.NopCloser(bytes.NewReader(body))
	}

	return &CapturedResponse{
		ID:        id,
		Timestamp: time.Now(),
		Status:    resp.StatusCode,
		Headers:   headers,
		Body:      body,
	}, nil
}

// 全局捕获器
var defaultCaptor = NewCaptor()

// CaptureFromRequest 全局捕获方法（向后兼容）
// 返回捕获的请求，忽略错误（向后兼容旧行为）
func CaptureFromRequest(req *http.Request) *CapturedRequest {
	captured, _ := defaultCaptor.CaptureFromRequest(req)
	return captured
}

// CaptureFromResponse 全局捕获响应方法
func CaptureFromResponse(resp *http.Response, id string) *CapturedResponse {
	captured, _ := defaultCaptor.CaptureFromResponse(resp, id)
	return captured
}

// SetDefaultCaptor 设置默认捕获器
func SetDefaultCaptor(c Captor) {
	if c != nil {
		if dc, ok := c.(*DefaultCaptor); ok {
			defaultCaptor = dc
		}
	}
}

// CaptureRequestToStore 捕获请求并存储到指定存储器
func CaptureRequestToStore(store Storer, req *http.Request) (string, *CapturedRequest, error) {
	captured, err := defaultCaptor.CaptureFromRequest(req)
	if err != nil {
		return "", nil, err
	}
	id := store.Capture(captured)
	return id, captured, nil
}

// CaptureResponseToStore 捕获响应并存储到指定存储器
func CaptureResponseToStore(store Storer, id string, resp *http.Response) error {
	captured, err := defaultCaptor.CaptureFromResponse(resp, id)
	if err != nil {
		return err
	}
	store.SetResponse(id, captured)
	return nil
}