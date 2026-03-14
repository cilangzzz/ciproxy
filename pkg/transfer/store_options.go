/**
  @creator: cilang
  @since: 2024
  @desc: 流量存储选项配置
**/

package transfer

// 默认配置
const (
	// DefaultMaxEntries 默认最大条目数
	DefaultMaxEntries = 10000
)

// StoreOption 存储选项函数
type StoreOption func(*TrafficStore)

// WithMaxEntries 设置最大条目数
// 设置为 0 表示无限制
func WithMaxEntries(max int) StoreOption {
	return func(s *TrafficStore) {
		s.maxEntries = max
	}
}

// WithEvictCallback 设置淘汰回调
// 当条目被淘汰时调用
func WithEvictCallback(fn func(*TrafficEntry)) StoreOption {
	return func(s *TrafficStore) {
		s.onEvict = fn
	}
}

// WithIDGenerator 设置ID生成器
func WithIDGenerator(g IDGenerator) StoreOption {
	return func(s *TrafficStore) {
		if g != nil {
			s.idGenerator = g
		}
	}
}