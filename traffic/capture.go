package traffic

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"
)

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

// TrafficStore 流量存储
type TrafficStore struct {
	mu      sync.RWMutex
	entries map[string]*TrafficEntry
	order   []string
}

var globalStore *TrafficStore
var storeOnce sync.Once

// GetTrafficStore 获取全局流量存储
func GetTrafficStore() *TrafficStore {
	storeOnce.Do(func() {
		globalStore = &TrafficStore{
			entries: make(map[string]*TrafficEntry),
			order:   make([]string, 0),
		}
	})
	return globalStore
}

// Capture 捕获请求
func (s *TrafficStore) Capture(req *CapturedRequest) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := &TrafficEntry{
		ID:      generateID(),
		Request: *req,
	}
	s.entries[entry.ID] = entry
	s.order = append(s.order, entry.ID)
	return entry.ID
}

// SetResponse 设置响应
func (s *TrafficStore) SetResponse(id string, resp *CapturedResponse) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entry, ok := s.entries[id]; ok {
		entry.Response = *resp
		entry.Duration = resp.Timestamp.Sub(entry.Request.Timestamp)
	}
}

// Get 获取流量条目
func (s *TrafficStore) Get(id string) (*TrafficEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.entries[id]
	return entry, ok
}

// List 列出所有流量条目
func (s *TrafficStore) List() []*TrafficEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries := make([]*TrafficEntry, 0, len(s.order))
	for _, id := range s.order {
		if entry, ok := s.entries[id]; ok {
			entries = append(entries, entry)
		}
	}
	return entries
}

// Delete 删除流量条目
func (s *TrafficStore) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.entries, id)
	for i, entryID := range s.order {
		if entryID == id {
			s.order = append(s.order[:i], s.order[i+1:]...)
			break
		}
	}
}

// Clear 清空所有流量
func (s *TrafficStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.entries = make(map[string]*TrafficEntry)
	s.order = make([]string, 0)
}

// Export 导出流量
func (s *TrafficStore) Export() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return json.MarshalIndent(s.List(), "", "  ")
}

// Import 导入流量
func (s *TrafficStore) Import(data []byte) error {
	var entries []*TrafficEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, entry := range entries {
		s.entries[entry.ID] = entry
		s.order = append(s.order, entry.ID)
	}
	return nil
}

// CaptureFromRequest 从http.Request捕获
func CaptureFromRequest(req *http.Request) *CapturedRequest {
	headers := make(map[string]string)
	for k, v := range req.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	var body []byte
	if req.Body != nil {
		body, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewReader(body))
	}

	return &CapturedRequest{
		ID:        generateID(),
		Timestamp: time.Now(),
		Method:    req.Method,
		URL:       req.URL.String(),
		Host:      req.Host,
		Headers:   headers,
		Body:      body,
		Protocol:  req.Proto,
	}
}

func generateID() string {
	return time.Now().Format("20060102150405.999999999")
}
