package transfer

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

// ErrEntryNotFound 条目未找到错误
var ErrEntryNotFound = errors.New("traffic entry not found")

// ModifyRule 请求修改规则
type ModifyRule func(req *http.Request) *http.Request

// Replayer 流量重放器
type Replayer struct {
	client       *http.Client
	modifyRules  []ModifyRule
	beforeReplay func(req *http.Request) error
	afterReplay  func(resp *http.Response) error
}

// NewReplayer 创建新的重放器
func NewReplayer() *Replayer {
	return &Replayer{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			Timeout: 30 * time.Second,
		},
	}
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

// Replay 重放请求
func (r *Replayer) Replay(entry *TrafficEntry) (*CapturedResponse, error) {
	req, err := r.reconstructRequest(&entry.Request)
	if err != nil {
		return nil, err
	}

	for _, rule := range r.modifyRules {
		req = rule(req)
	}

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

func (r *Replayer) reconstructRequest(captured *CapturedRequest) (*http.Request, error) {
	u, err := url.Parse(captured.URL)
	if err != nil {
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
func (r *Replayer) ReplayByID(store *TrafficStore, id string) (*CapturedResponse, error) {
	entry, ok := store.Get(id)
	if !ok {
		return nil, ErrEntryNotFound
	}
	return r.Replay(entry)
}

// ReplayAll 重放所有请求
func (r *Replayer) ReplayAll(store *TrafficStore) ([]*CapturedResponse, []error) {
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