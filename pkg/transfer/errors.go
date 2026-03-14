/**
  @creator: cilang
  @since: 2024
  @desc: 流量模块错误定义
**/

package transfer

import "errors"

// 错误变量，遵循 Err 前缀规范
var (
	// ErrEntryNotFound 流量条目未找到
	ErrEntryNotFound = errors.New("traffic entry not found")
	// ErrStoreFull 流量存储已满
	ErrStoreFull = errors.New("traffic store is full")
	// ErrInvalidKey 无效的加密密钥
	ErrInvalidKey = errors.New("invalid encryption key: must be 16, 24 or 32 bytes")
	// ErrEncryptFailed 加密失败
	ErrEncryptFailed = errors.New("encryption failed")
	// ErrDecryptFailed 解密失败
	ErrDecryptFailed = errors.New("decryption failed")
	// ErrCiphertextTooShort 密文太短
	ErrCiphertextTooShort = errors.New("ciphertext too short")
	// ErrInvalidRequest 无效的请求
	ErrInvalidRequest = errors.New("invalid request")
	// ErrTransferFailed 转发失败
	ErrTransferFailed = errors.New("transfer failed")
	// ErrConnectionClosed 连接已关闭
	ErrConnectionClosed = errors.New("connection closed")
	// ErrNilCryptor 加密器为空
	ErrNilCryptor = errors.New("cryptor is nil")
)