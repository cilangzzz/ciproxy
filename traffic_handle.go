/**
  @author: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @since: 2023/6/11
  @desc: //TODO
**/

// Package ciproxy 流量转发
package ciproxy

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/opencvlzg/ciproxy/pkg/util"
	"io"
	"log"
	"net/http"
)

// TrafficCryptor 流量加密器
type TrafficCryptor struct {
	key []byte
	gcm cipher.AEAD
}

// NewTrafficCryptor 创建新的流量加密器
func NewTrafficCryptor(key []byte) (*TrafficCryptor, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("key must be 16, 24 or 32 bytes")
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
		return nil, err
	}

	return tc.gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt 解密数据
func (tc *TrafficCryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	nonceSize := tc.gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return tc.gcm.Open(nil, nonce, ciphertext, nil)
}

// CryptTraffic 流量加密
func CryptTraffic(data []byte, cryptor *TrafficCryptor) ([]byte, error) {
	if cryptor == nil {
		return data, nil
	}
	return cryptor.Encrypt(data)
}

// DecryptTraffic 流量解密
func DecryptTraffic(data []byte, cryptor *TrafficCryptor) ([]byte, error) {
	if cryptor == nil {
		return data, nil
	}
	return cryptor.Decrypt(data)
}

// TunnelTransfer 加密流量转发
func TunnelTransfer(client, server io.ReadWriteCloser, cryptor *TrafficCryptor) {
	// 客户端到服务器：解密后转发
	go func() {
		defer client.Close()
		defer server.Close()
		buf := make([]byte, 4096)
		for {
			n, err := client.Read(buf)
			if err != nil {
				return
			}
			decrypted, err := DecryptTraffic(buf[:n], cryptor)
			if err != nil {
				log.Println("decrypt error:", err)
				return
			}
			_, err = server.Write(decrypted)
			if err != nil {
				return
			}
		}
	}()

	// 服务器到客户端：加密后转发
	go func() {
		defer client.Close()
		defer server.Close()
		buf := make([]byte, 4096)
		for {
			n, err := server.Read(buf)
			if err != nil {
				return
			}
			encrypted, err := CryptTraffic(buf[:n], cryptor)
			if err != nil {
				log.Println("encrypt error:", err)
				return
			}
			_, err = client.Write(encrypted)
			if err != nil {
				return
			}
		}
	}()
}

// Transfer traffic transfer 流量Io转发
func Transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer func(destination io.WriteCloser) {
		err := destination.Close()
		if err != nil {
			//log.Println("close io writer failed", err)
		}
	}(destination)
	defer func(source io.ReadCloser) {
		err := source.Close()
		if err != nil {
			//log.Println("close io writer failed", err)
		}
	}(source)
	_, err := io.Copy(destination, source)
	if err != nil {
		log.Println("copy data transfer failed", err)
	}
}

// TeeTransfer traffic transfer 流量Io转发
func TeeTransfer(destination io.WriteCloser, source io.ReadCloser) {
	defer func(destination io.WriteCloser) {
		err := destination.Close()
		if err != nil {
			//log.Println("close io writer failed", err)
		}
	}(destination)
	defer func(source io.ReadCloser) {
		err := source.Close()
		if err != nil {
			//log.Println("close io writer failed", err)
		}
	}(source)
	teeReader := io.TeeReader(source, DefaultWriter)
	//http.Read
	_, err := io.Copy(destination, teeReader)
	if err != nil {
		log.Println("copy data transfer failed", err)
	}
}

// TeeDoRequestTransfer traffic transfer 流量Io转发,手动处理请求
func TeeDoRequestTransfer(c *Context) {
	//defer func(destination io.WriteCloser) {
	//	err := destination.Close()
	//	if err != nil {
	//		//log.Println("close io writer failed", err)
	//	}
	//}(c.TlsClientConn)
	//defer func(source io.ReadCloser) {
	//	err := source.Close()
	//	if err != nil {
	//		//log.Println("close io writer failed", err)
	//	}
	//}(c.TlsServerConn)

	//teeReader := io.TeeReader(c, DefaultWriter)
	// http.Read
	cReader := bufio.NewReader(c.TlsClientConn)
	request, err := http.ReadRequest(cReader)
	if err != nil {
		return
	}
	//println(request.Host)
	newReq, _ := util.NewRequest(request)
	// 创建 HTTP 请求, TODO 需要加入池优化性能
	// TODO 从TLS连接创建一个HTTP客户端
	//client := &http.Client{
	//	Transport: &http.Transport{
	//		TLSClientConfig: tlsConfig,
	//		DialTLS: func(network, addr string) (net.Conn, error) {
	//			return tlsClient, nil
	//		},
	//	},
	//}
	//request.RequestURI = ""
	client := http.Client{}
	//request.RemoteAddr = c.TslServerConn.RemoteAddr().String()
	//request.URL.Scheme = "http"
	//request.URL.Host = "cn.bing.com"
	//u, err := url.Parse(request.URL.Scheme + "://" + request.Host + request.URL.String())
	//if err != nil {
	//	panic(err)
	//}
	//request.URL = u
	//c.Req = req
	newReq.URL.Scheme = "https"
	newReq.URL.Host = request.Host
	resp, err := client.Do(newReq)
	//resp, err := client.Get("https://cn.bing.com")
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
	//bytes := util.HttpContext("你好")
	//c.Write(bytes)

}
