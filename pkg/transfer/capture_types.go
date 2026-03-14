/**
  @creator: cilang
  @since: 2024
  @desc: 流量捕获类型定义
**/

package transfer

import "time"

// CapturedRequest 捕获的HTTP请求
type CapturedRequest struct {
	ID        string            `json:"id"`
	Timestamp time.Time         `json:"timestamp"`
	Method    string            `json:"method"`
	URL       string            `json:"url"`
	Host      string            `json:"host"`
	Headers   map[string]string `json:"headers"`
	Body      []byte            `json:"body"`
	Protocol  string            `json:"protocol"`
}

// CapturedResponse 捕获的HTTP响应
type CapturedResponse struct {
	ID        string            `json:"id"`
	Timestamp time.Time         `json:"timestamp"`
	Status    int               `json:"status"`
	Headers   map[string]string `json:"headers"`
	Body      []byte            `json:"body"`
}

// TrafficEntry 完整的请求/响应对
type TrafficEntry struct {
	ID       string           `json:"id"`
	Request  CapturedRequest  `json:"request"`
	Response CapturedResponse `json:"response"`
	Duration time.Duration    `json:"duration"`
	Tags     []string         `json:"tags,omitempty"`
}

// Reset 重置 TrafficEntry 字段（用于对象池）
func (e *TrafficEntry) Reset() {
	e.ID = ""
	e.Request = CapturedRequest{}
	e.Response = CapturedResponse{}
	e.Duration = 0
	e.Tags = e.Tags[:0]
}

// Reset 重置 CapturedRequest 字段（用于对象池）
func (r *CapturedRequest) Reset() {
	r.ID = ""
	r.Timestamp = time.Time{}
	r.Method = ""
	r.URL = ""
	r.Host = ""
	r.Headers = make(map[string]string)
	r.Body = r.Body[:0]
	r.Protocol = ""
}

// Reset 重置 CapturedResponse 字段（用于对象池）
func (r *CapturedResponse) Reset() {
	r.ID = ""
	r.Timestamp = time.Time{}
	r.Status = 0
	r.Headers = make(map[string]string)
	r.Body = r.Body[:0]
}