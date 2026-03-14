/**
  @creator: cilang
  @since: 2024
  @desc: 流量存储管理
**/

package transfer

import (
	"encoding/json"
	"sync"
	"time"
)

// Storer 存储接口（遵循 -er 后缀规范）
type Storer interface {
	Capture(req *CapturedRequest) string
	SetResponse(id string, resp *CapturedResponse)
	Get(id string) (*TrafficEntry, bool)
	List() []*TrafficEntry
	Delete(id string) error
	Clear()
	Export() ([]byte, error)
	Import(data []byte) error
	Len() int
	Filter(fn func(*TrafficEntry) bool) []*TrafficEntry
}

// TrafficStore 流量存储实现
type TrafficStore struct {
	mu          sync.RWMutex
	entries     map[string]*TrafficEntry
	order       []string
	maxEntries  int                   // 最大条目数，0 表示无限制
	onEvict     func(*TrafficEntry)   // 淘汰回调
	idGenerator IDGenerator           // ID生成器
}

// 全局实例
var (
	globalStore *TrafficStore
	storeOnce   sync.Once
)

// GetTrafficStore 获取全局流量存储（单例模式）
func GetTrafficStore() *TrafficStore {
	storeOnce.Do(func() {
		globalStore = NewTrafficStore()
	})
	return globalStore
}

// NewTrafficStore 创建新的流量存储
func NewTrafficStore(opts ...StoreOption) *TrafficStore {
	s := &TrafficStore{
		entries:     make(map[string]*TrafficEntry),
		order:       make([]string, 0),
		maxEntries:  DefaultMaxEntries,
		idGenerator: defaultIDGenerator,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Capture 捕获请求
func (s *TrafficStore) Capture(req *CapturedRequest) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 如果请求没有ID，生成一个
	if req.ID == "" {
		req.ID = s.idGenerator.Generate()
	}

	// 检查容量限制，执行淘汰策略
	if s.maxEntries > 0 && len(s.entries) >= s.maxEntries {
		s.evictOldest()
	}

	entry := &TrafficEntry{
		ID:      req.ID,
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
func (s *TrafficStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.entries[id]; !ok {
		return ErrEntryNotFound
	}

	delete(s.entries, id)
	for i, entryID := range s.order {
		if entryID == id {
			s.order = append(s.order[:i], s.order[i+1:]...)
			break
		}
	}
	return nil
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

// Len 获取条目数量
func (s *TrafficStore) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.entries)
}

// evictOldest 淘汰最老的条目（FIFO策略）
func (s *TrafficStore) evictOldest() {
	if len(s.order) == 0 {
		return
	}
	oldestID := s.order[0]
	if s.onEvict != nil {
		if entry, ok := s.entries[oldestID]; ok {
			s.onEvict(entry)
		}
	}
	delete(s.entries, oldestID)
	s.order = s.order[1:]
}

// SetMaxEntries 设置最大条目数
func (s *TrafficStore) SetMaxEntries(max int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.maxEntries = max
	// 如果当前条目数超过新的限制，淘汰多余的条目
	for s.maxEntries > 0 && len(s.entries) > s.maxEntries {
		s.evictOldest()
	}
}

// GetByID 根据ID获取流量条目（Get的别名，便于链式调用）
func (s *TrafficStore) GetByID(id string) (*TrafficEntry, bool) {
	return s.Get(id)
}

// Exists 检查条目是否存在
func (s *TrafficStore) Exists(id string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.entries[id]
	return ok
}

// SetTags 设置条目标签
func (s *TrafficStore) SetTags(id string, tags []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if entry, ok := s.entries[id]; ok {
		entry.Tags = tags
	}
}

// AddTag 添加标签
func (s *TrafficStore) AddTag(id string, tag string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if entry, ok := s.entries[id]; ok {
		entry.Tags = append(entry.Tags, tag)
	}
}

// GetByTag 根据标签获取条目
func (s *TrafficStore) GetByTag(tag string) []*TrafficEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*TrafficEntry
	for _, id := range s.order {
		if entry, ok := s.entries[id]; ok {
			for _, t := range entry.Tags {
				if t == tag {
					result = append(result, entry)
					break
				}
			}
		}
	}
	return result
}

// Filter 过滤条目
func (s *TrafficStore) Filter(fn func(*TrafficEntry) bool) []*TrafficEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*TrafficEntry
	for _, id := range s.order {
		if entry, ok := s.entries[id]; ok && fn(entry) {
			result = append(result, entry)
		}
	}
	return result
}

// Range 遍历所有条目
func (s *TrafficStore) Range(fn func(*TrafficEntry) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, id := range s.order {
		if entry, ok := s.entries[id]; ok {
			if !fn(entry) {
				break
			}
		}
	}
}

// Stats 返回存储统计信息
func (s *TrafficStore) Stats() StoreStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var totalDuration time.Duration
	for _, entry := range s.entries {
		totalDuration += entry.Duration
	}

	var avgDuration time.Duration
	if len(s.entries) > 0 {
		avgDuration = totalDuration / time.Duration(len(s.entries))
	}

	return StoreStats{
		TotalEntries: len(s.entries),
		MaxEntries:   s.maxEntries,
		AvgDuration:  avgDuration,
	}
}

// StoreStats 存储统计信息
type StoreStats struct {
	TotalEntries int
	MaxEntries   int
	AvgDuration  time.Duration
}