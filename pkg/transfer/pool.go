/**
  @creator: cilang
  @since: 2024
  @desc: 对象池
**/

package transfer

import "sync"

// TrafficEntryPool TrafficEntry 对象池
var TrafficEntryPool = sync.Pool{
	New: func() interface{} {
		return &TrafficEntry{}
	},
}

// CapturedRequestPool CapturedRequest 对象池
var CapturedRequestPool = sync.Pool{
	New: func() interface{} {
		return &CapturedRequest{
			Headers: make(map[string]string),
		}
	},
}

// CapturedResponsePool CapturedResponse 对象池
var CapturedResponsePool = sync.Pool{
	New: func() interface{} {
		return &CapturedResponse{
			Headers: make(map[string]string),
		}
	},
}

// GetTrafficEntry 从池获取 TrafficEntry
func GetTrafficEntry() *TrafficEntry {
	return TrafficEntryPool.Get().(*TrafficEntry)
}

// PutTrafficEntry 归还 TrafficEntry
func PutTrafficEntry(entry *TrafficEntry) {
	entry.Reset()
	TrafficEntryPool.Put(entry)
}

// GetCapturedRequest 从池获取 CapturedRequest
func GetCapturedRequest() *CapturedRequest {
	return CapturedRequestPool.Get().(*CapturedRequest)
}

// PutCapturedRequest 归还 CapturedRequest
func PutCapturedRequest(req *CapturedRequest) {
	req.Reset()
	CapturedRequestPool.Put(req)
}

// GetCapturedResponse 从池获取 CapturedResponse
func GetCapturedResponse() *CapturedResponse {
	return CapturedResponsePool.Get().(*CapturedResponse)
}

// PutCapturedResponse 归还 CapturedResponse
func PutCapturedResponse(resp *CapturedResponse) {
	resp.Reset()
	CapturedResponsePool.Put(resp)
}

// BufferPool 通用缓冲区池
var BufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 0, 4096)
	},
}

// GetBuffer 从池获取缓冲区
func GetBuffer() []byte {
	return BufferPool.Get().([]byte)[:0]
}

// PutBuffer 归还缓冲区
func PutBuffer(buf []byte) {
	BufferPool.Put(buf[:0])
}

// GetBufferWithCap 从池获取指定容量的缓冲区
func GetBufferWithCap(capacity int) []byte {
	buf := BufferPool.Get().([]byte)
	if cap(buf) < capacity {
		return make([]byte, 0, capacity)
	}
	return buf[:0]
}