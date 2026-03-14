/**
  @author: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @since: 2023/6/11
  @desc: 流量转发（使用 pkg/transfer 模块实现）
**/

// Package ciproxy 流量转发
package ciproxy

import (
	"bufio"
	"fmt"
	"github.com/opencvlzg/ciproxy/pkg/transfer"
	"github.com/opencvlzg/ciproxy/pkg/util"
	"io"
	"net/http"
)

// ========== 类型别名（向后兼容） ==========

// TrafficCryptor 流量加密器
type TrafficCryptor = transfer.TrafficCryptor

// Cryptor 加解密接口
type Cryptor = transfer.Cryptor

// ========== 包装函数（向后兼容） ==========

// NewTrafficCryptor 创建新的流量加密器
func NewTrafficCryptor(key []byte) (*TrafficCryptor, error) {
	return transfer.NewTrafficCryptor(key)
}

// CryptTraffic 流量加密
func CryptTraffic(data []byte, cryptor *TrafficCryptor) ([]byte, error) {
	return transfer.CryptTraffic(data, cryptor)
}

// DecryptTraffic 流量解密
func DecryptTraffic(data []byte, cryptor *TrafficCryptor) ([]byte, error) {
	return transfer.DecryptTraffic(data, cryptor)
}

// Transfer traffic transfer 流量Io转发
// 注意：忽略错误以保持向后兼容
func Transfer(destination io.WriteCloser, source io.ReadCloser) {
	_ = transfer.Transfer(destination, source)
}

// TeeTransfer traffic transfer 流量Io转发
// 使用 DefaultWriter 作为 tee 输出
func TeeTransfer(destination io.WriteCloser, source io.ReadCloser) {
	_ = transfer.TeeTransfer(destination, source, DefaultWriter)
}

// TunnelTransfer 加密流量转发
func TunnelTransfer(client, server io.ReadWriteCloser, cryptor *TrafficCryptor) {
	transfer.TunnelTransfer(client, server, cryptor)
}

// ========== 独特功能（无 pkg 等价实现） ==========

// TeeDoRequestTransfer traffic transfer 流量Io转发,手动处理请求
func TeeDoRequestTransfer(c *Context) {
	cReader := bufio.NewReader(c.TlsClientConn)
	request, err := http.ReadRequest(cReader)
	if err != nil {
		return
	}

	newReq, _ := util.NewRequest(request)
	client := http.Client{}
	newReq.URL.Scheme = "https"
	newReq.URL.Host = request.Host
	resp, err := client.Do(newReq)
	if err != nil {
		fmt.Println("Error do request:", err)
		return
	}
	bytes, err := util.ResponseToBytes(resp)
	if err != nil {
		return
	}

	// 将响应内容写回 TCP 连接
	c.TlsClientConn.Write(bytes)
}