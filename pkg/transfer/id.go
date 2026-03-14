/**
  @creator: cilang
  @since: 2024
  @desc: 唯一ID生成器
**/

package transfer

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"
)

// IDGenerator ID生成器接口
type IDGenerator interface {
	Generate() string
}

// TimeBasedIDGenerator 基于时间的ID生成器
// 使用时间戳 + 原子计数器 + 随机数确保唯一性
type TimeBasedIDGenerator struct {
	counter uint64
}

// NewTimeBasedIDGenerator 创建时间ID生成器
func NewTimeBasedIDGenerator() *TimeBasedIDGenerator {
	return &TimeBasedIDGenerator{}
}

// Generate 生成唯一ID
func (g *TimeBasedIDGenerator) Generate() string {
	ts := time.Now().UnixNano()
	counter := atomic.AddUint64(&g.counter, 1)
	random := make([]byte, 4)
	rand.Read(random)
	return fmt.Sprintf("%d-%d-%s", ts, counter, hex.EncodeToString(random))
}

// 全局ID生成器
var defaultIDGenerator IDGenerator = NewTimeBasedIDGenerator()

// generateID 生成唯一ID（内部使用）
func generateID() string {
	return defaultIDGenerator.Generate()
}

// SetIDGenerator 设置自定义ID生成器
func SetIDGenerator(g IDGenerator) {
	if g != nil {
		defaultIDGenerator = g
	}
}

// GetIDGenerator 获取当前ID生成器
func GetIDGenerator() IDGenerator {
	return defaultIDGenerator
}