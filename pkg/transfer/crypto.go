/**
  @creator: cilang
  @since: 2024
  @desc: 流量加解密处理
**/

package transfer

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// Cryptor 加解密接口
type Cryptor interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}

// TrafficCryptor 流量加密器
type TrafficCryptor struct {
	key []byte
	gcm cipher.AEAD
}

// NewTrafficCryptor 创建新的流量加密器
// key 必须是 16, 24 或 32 字节
func NewTrafficCryptor(key []byte) (*TrafficCryptor, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, ErrInvalidKey
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &TrafficCryptor{key: key, gcm: gcm}, nil
}

// Encrypt 加密数据
func (tc *TrafficCryptor) Encrypt(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, tc.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, ErrEncryptFailed
	}
	return tc.gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt 解密数据
func (tc *TrafficCryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	nonceSize := tc.gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, ErrCiphertextTooShort
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return tc.gcm.Open(nil, nonce, ciphertext, nil)
}

// NonceSize 返回 nonce 大小
func (tc *TrafficCryptor) NonceSize() int {
	return tc.gcm.NonceSize()
}

// Key 返回密钥
func (tc *TrafficCryptor) Key() []byte {
	return tc.key
}

// CryptTraffic 流量加密
// 如果 cryptor 为 nil，返回原始数据
func CryptTraffic(data []byte, cryptor *TrafficCryptor) ([]byte, error) {
	if cryptor == nil {
		return data, nil
	}
	return cryptor.Encrypt(data)
}

// DecryptTraffic 流量解密
// 如果 cryptor 为 nil，返回原始数据
func DecryptTraffic(data []byte, cryptor *TrafficCryptor) ([]byte, error) {
	if cryptor == nil {
		return data, nil
	}
	return cryptor.Decrypt(data)
}

// EncryptWithCryptor 使用 Cryptor 接口加密
func EncryptWithCryptor(data []byte, cryptor Cryptor) ([]byte, error) {
	if cryptor == nil {
		return data, nil
	}
	return cryptor.Encrypt(data)
}

// DecryptWithCryptor 使用 Cryptor 接口解密
func DecryptWithCryptor(data []byte, cryptor Cryptor) ([]byte, error) {
	if cryptor == nil {
		return data, nil
	}
	return cryptor.Decrypt(data)
}

// NOPCryptor 空加密器（不进行加密/解密操作）
type NOPCryptor struct{}

// NewNOPCryptor 创建空加密器
func NewNOPCryptor() *NOPCryptor {
	return &NOPCryptor{}
}

// Encrypt 返回原始数据
func (c *NOPCryptor) Encrypt(plaintext []byte) ([]byte, error) {
	return plaintext, nil
}

// Decrypt 返回原始数据
func (c *NOPCryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	return ciphertext, nil
}