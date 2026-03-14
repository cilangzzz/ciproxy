/**
  @creator: cilang
  @since: 2024
  @desc: 流量重放
**/

package transfer

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// ModifyRule 请求修改规则
type ModifyRule func(req *http.Request) *http.Request

// ReplayOption 重放器选项函数
type ReplayOption func(*Replayer)

// WithTLSConfig 设置 TLS 配置
func WithTLSConfig(cfg *tls.Config) ReplayOption {
	return func(r *Replayer) {
		if transport, ok := r.client.Transport.(*http.Transport); ok {
			transport.TLSClientConfig = cfg
		}
	}
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) ReplayOption {
	return func(r *Replayer) {
		r.client.Timeout = timeout
	}
}

// WithModifyRule 添加修改规则
func WithModifyRule(rule ModifyRule) ReplayOption {
	return func(r *Replayer) {
		r.modifyRules = append(r.modifyRules, rule)
	}
}

// WithBeforeReplay 设置重放前钩子
func WithBeforeReplay(fn func(req *http.Request) error) ReplayOption {
	return func(r *Replayer) {
		r.beforeReplay = fn
	}
}

// WithAfterReplay 设置重放后钩子
func WithAfterReplay(fn func(resp *http.Response) error) ReplayOption {
	return func(r *Replayer) {
		r.afterReplay = fn
	}
}

// WithInsecureSkipVerify 设置跳过 TLS 验证
func WithInsecureSkipVerify(skip bool) ReplayOption {
	return func(r *Replayer) {
		if transport, ok := r.client.Transport.(*http.Transport); ok {
			transport.TLSClientConfig.InsecureSkipVerify = skip
		}
	}
}

// Replayer 流量重放器
type Replayer struct {
	client       *http.Client
	modifyRules  []ModifyRule
	beforeReplay func(req *http.Request) error
	afterReplay  func(resp *http.Response) error
}

// NewReplayer 创建新的重放器
func NewReplayer(opts ...ReplayOption) *Replayer {
	r := &Replayer{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			Timeout: 30 * time.Second,
		},
		modifyRules: make([]ModifyRule, 0),
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// AddModifyRule 添加修改规则
func (r *Replayer) AddModifyRule(rule ModifyRule) {
	r.modifyRules = append(r.modifyRules, rule)
}

// SetBeforeReplay 设置重放前钩子
func (r *Replayer) SetBeforeReplay(fn func(req *http.Request) error) {
	r.beforeReplay = fn
}

// SetAfterReplay 设置重放后钩子
func (r *Replayer) SetAfterReplay(fn func(resp *http.Response) error) {
	r.afterReplay = fn
}

// SetTimeout 设置超时时间
func (r *Replayer) SetTimeout(timeout time.Duration) {
	r.client.Timeout = timeout
}

// SetTLSConfig 设置 TLS 配置
func (r *Replayer) SetTLSConfig(cfg *tls.Config) {
	if transport, ok := r.client.Transport.(*http.Transport); ok {
		transport.TLSClientConfig = cfg
	}
}

// Replay 重放请求
func (r *Replayer) Replay(entry *TrafficEntry) (*CapturedResponse, error) {
	req, err := r.reconstructRequest(&entry.Request)
	if err != nil {
		return nil, err
	}

	// 应用修改规则
	for _, rule := range r.modifyRules {
		req = rule(req)
	}

	// 执行重放前钩子
	if r.beforeReplay != nil {
		if err := r.beforeReplay(req); err != nil {
			return nil, err
		}
	}

	startTime := time.Now()
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 执行重放后钩子
	if r.afterReplay != nil {
		if err := r.afterReplay(resp); err != nil {
			return nil, err
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	headers := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	return &CapturedResponse{
		ID:        entry.ID,
		Timestamp: startTime,
		Status:    resp.StatusCode,
		Headers:   headers,
		Body:      body,
	}, nil
}

// reconstructRequest 重建请求
func (r *Replayer) reconstructRequest(captured *CapturedRequest) (*http.Request, error) {
	u, err := url.Parse(captured.URL)
	if err != nil {
		// 尝试构造 URL
		u, _ = url.Parse("https://" + captured.Host + "/")
	}

	req, err := http.NewRequest(captured.Method, u.String(), bytes.NewReader(captured.Body))
	if err != nil {
		return nil, err
	}

	for k, v := range captured.Headers {
		req.Header.Set(k, v)
	}

	req.Host = captured.Host
	return req, nil
}

// ReplayByID 通过ID重放
func (r *Replayer) ReplayByID(store Storer, id string) (*CapturedResponse, error) {
	entry, ok := store.Get(id)
	if !ok {
		return nil, ErrEntryNotFound
	}
	return r.Replay(entry)
}

// ReplayAll 重放所有请求
func (r *Replayer) ReplayAll(store Storer) ([]*CapturedResponse, []error) {
	entries := store.List()
	responses := make([]*CapturedResponse, 0, len(entries))
	errs := make([]error, 0)

	for _, entry := range entries {
		resp, err := r.Replay(entry)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		responses = append(responses, resp)
	}

	return responses, errs
}

// ReplayParallel 并行重放所有请求
func (r *Replayer) ReplayParallel(store Storer, concurrency int) ([]*CapturedResponse, []error) {
	entries := store.List()
	if len(entries) == 0 {
		return nil, nil
	}

	if concurrency <= 0 {
		concurrency = 10
	}

	responses := make([]*CapturedResponse, 0, len(entries))
	errs := make([]error, 0)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 使用信号量控制并发
	sem := make(chan struct{}, concurrency)

	for _, entry := range entries {
		wg.Add(1)
		go func(e *TrafficEntry) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			resp, err := r.Replay(e)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errs = append(errs, err)
				return
			}
			responses = append(responses, resp)
		}(entry)
	}

	wg.Wait()
	return responses, errs
}

// ReplayFiltered 过滤后重放
func (r *Replayer) ReplayFiltered(store Storer, filter func(*TrafficEntry) bool) ([]*CapturedResponse, []error) {
	entries := store.Filter(filter)
	responses := make([]*CapturedResponse, 0, len(entries))
	errs := make([]error, 0)

	for _, entry := range entries {
		resp, err := r.Replay(entry)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		responses = append(responses, resp)
	}

	return responses, errs
}

// ReplayWithCallback 带回调的重放
func (r *Replayer) ReplayWithCallback(store Storer, callback func(*TrafficEntry, *CapturedResponse, error)) {
	entries := store.List()
	for _, entry := range entries {
		resp, err := r.Replay(entry)
		callback(entry, resp, err)
	}
}