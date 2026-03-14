/**
  @creator: cilang
  @since: 2024
  @desc: 隧道流量转发
**/

package transfer

import (
	"io"
	"log"
	"sync"
)

// 默认配置
const (
	// DefaultBufferSize 默认缓冲区大小
	DefaultBufferSize = 4096
)

// TransferOption 转发选项函数
type TransferOption func(*TransferOptions)

// TransferOptions 转发选项
type TransferOptions struct {
	BufferSize int
	Cryptor    Cryptor
	OnError    func(error)
	OnData     func([]byte, bool) // 数据回调，bool 表示是否为加密方向
}

// DefaultTransferOptions 默认转发选项
var DefaultTransferOptions = TransferOptions{
	BufferSize: DefaultBufferSize,
}

// WithBufferSize 设置缓冲区大小
func WithBufferSize(size int) TransferOption {
	return func(o *TransferOptions) {
		if size > 0 {
			o.BufferSize = size
		}
	}
}

// WithCryptor 设置加密器
func WithCryptor(c Cryptor) TransferOption {
	return func(o *TransferOptions) {
		o.Cryptor = c
	}
}

// WithErrorCallback 设置错误回调
func WithErrorCallback(fn func(error)) TransferOption {
	return func(o *TransferOptions) {
		o.OnError = fn
	}
}

// WithDataCallback 设置数据回调
func WithDataCallback(fn func([]byte, bool)) TransferOption {
	return func(o *TransferOptions) {
		o.OnData = fn
	}
}

// Transfer 流量转发（单向）
func Transfer(dst io.WriteCloser, src io.ReadCloser, opts ...TransferOption) error {
	opt := DefaultTransferOptions
	for _, o := range opts {
		o(&opt)
	}

	defer func() {
		dst.Close()
		src.Close()
	}()

	buf := make([]byte, opt.BufferSize)
	_, err := io.CopyBuffer(dst, src, buf)
	if err != nil && opt.OnError != nil {
		opt.OnError(err)
	}
	return err
}

// TeeTransfer 带 Tee 的流量转发
func TeeTransfer(dst io.WriteCloser, src io.ReadCloser, teeWriter io.Writer, opts ...TransferOption) error {
	opt := DefaultTransferOptions
	for _, o := range opts {
		o(&opt)
	}

	defer func() {
		dst.Close()
		src.Close()
	}()

	teeReader := io.TeeReader(src, teeWriter)
	buf := make([]byte, opt.BufferSize)
	_, err := io.CopyBuffer(dst, teeReader, buf)
	if err != nil && opt.OnError != nil {
		opt.OnError(err)
	}
	return err
}

// BidirectionalTransfer 双向转发
func BidirectionalTransfer(conn1, conn2 io.ReadWriteCloser, opts ...TransferOption) error {
	var wg sync.WaitGroup
	var transferErr error
	var mu sync.Mutex

	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := Transfer(conn1, conn2, opts...); err != nil {
			mu.Lock()
			transferErr = err
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		if err := Transfer(conn2, conn1, opts...); err != nil {
			mu.Lock()
			transferErr = err
			mu.Unlock()
		}
	}()

	wg.Wait()
	return transferErr
}

// TunnelTransfer 加密隧道转发
// 客户端 -> 服务器：解密后转发
// 服务器 -> 客户端：加密后转发
func TunnelTransfer(client, server io.ReadWriteCloser, cryptor *TrafficCryptor) {
	var wg sync.WaitGroup

	wg.Add(2)

	// 客户端 -> 服务器：解密后转发
	go func() {
		defer wg.Done()
		defer client.Close()
		defer server.Close()
		decryptAndForward(client, server, cryptor)
	}()

	// 服务器 -> 客户端：加密后转发
	go func() {
		defer wg.Done()
		defer client.Close()
		defer server.Close()
		encryptAndForward(server, client, cryptor)
	}()

	wg.Wait()
}

// TunnelTransferWithCryptor 使用 Cryptor 接口的隧道转发
func TunnelTransferWithCryptor(client, server io.ReadWriteCloser, cryptor Cryptor) {
	var wg sync.WaitGroup

	wg.Add(2)

	// 客户端 -> 服务器：解密后转发
	go func() {
		defer wg.Done()
		defer client.Close()
		defer server.Close()
		decryptAndForward(client, server, cryptor)
	}()

	// 服务器 -> 客户端：加密后转发
	go func() {
		defer wg.Done()
		defer client.Close()
		defer server.Close()
		encryptAndForward(server, client, cryptor)
	}()

	wg.Wait()
}

// decryptAndForward 解密后转发
func decryptAndForward(src io.Reader, dst io.Writer, cryptor interface{}) {
	buf := make([]byte, DefaultBufferSize)
	for {
		n, err := src.Read(buf)
		if err != nil {
			return
		}

		data := buf[:n]

		// 尝试解密
		switch c := cryptor.(type) {
		case *TrafficCryptor:
			decrypted, err := DecryptTraffic(data, c)
			if err != nil {
				log.Println("decrypt error:", err)
				return
			}
			data = decrypted
		case Cryptor:
			decrypted, err := c.Decrypt(data)
			if err != nil {
				log.Println("decrypt error:", err)
				return
			}
			data = decrypted
		case nil:
			// 无加密器，直接转发
		}

		if _, err := dst.Write(data); err != nil {
			return
		}
	}
}

// encryptAndForward 加密后转发
func encryptAndForward(src io.Reader, dst io.Writer, cryptor interface{}) {
	buf := make([]byte, DefaultBufferSize)
	for {
		n, err := src.Read(buf)
		if err != nil {
			return
		}

		data := buf[:n]

		// 尝试加密
		switch c := cryptor.(type) {
		case *TrafficCryptor:
			encrypted, err := CryptTraffic(data, c)
			if err != nil {
				log.Println("encrypt error:", err)
				return
			}
			data = encrypted
		case Cryptor:
			encrypted, err := c.Encrypt(data)
			if err != nil {
				log.Println("encrypt error:", err)
				return
			}
			data = encrypted
		case nil:
			// 无加密器，直接转发
		}

		if _, err := dst.Write(data); err != nil {
			return
		}
	}
}

// CopyWithCallback 带回调的数据复制
func CopyWithCallback(dst io.Writer, src io.Reader, callback func([]byte) error, bufferSize int) error {
	if bufferSize <= 0 {
		bufferSize = DefaultBufferSize
	}

	buf := make([]byte, bufferSize)
	for {
		n, err := src.Read(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		data := buf[:n]
		if callback != nil {
			if err := callback(data); err != nil {
				return err
			}
		}

		if _, err := dst.Write(data); err != nil {
			return err
		}
	}
}