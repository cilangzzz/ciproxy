/**
  @creator: cilang
  @since: 2024
  @desc: ALPN 协议选择器
**/

package mitm

import (
	"crypto/tls"
	"github.com/opencvlzg/ciproxy/internal/util"
	"log"
	"net"
	"strings"
)

// ALPNSelector ALPN 协议选择器
type ALPNSelector struct {
	interceptor *Interceptor
}

// NewALPNSelector 创建新的 ALPN 选择器
func NewALPNSelector(interceptor *Interceptor) *ALPNSelector {
	return &ALPNSelector{
		interceptor: interceptor,
	}
}

// SelectProtocol 根据 ALPN 选择协议
func (s *ALPNSelector) SelectProtocol(clientHello *tls.ClientHelloInfo) (string, error) {
	// 检查客户端支持的协议
	for _, proto := range clientHello.SupportedProtos {
		switch proto {
		case "h2":
			if s.interceptor.EnableHTTP2 {
				return "h2", nil
			}
		case "http/1.1":
			return "http/1.1", nil
		}
	}
	return "http/1.1", nil // 默认 HTTP/1.1
}

// HandleMITMConnection 完整的 MITM 连接处理
func (s *ALPNSelector) HandleMITMConnection(clientConn net.Conn, host string, ctx Context) error {
	// 1. 解析主机名
	if !strings.Contains(host, ":") {
		host += ":443"
	}

	// 2. 生成 TLS 配置（支持 ALPN）
	tlsConfig, err := util.GenerateTlsConfigWithALPN(host, s.interceptor.EnableHTTP2)
	if err != nil {
		log.Println("generate tls config error:", err)
		return err
	}

	// 3. 建立到目标服务器的 TLS 连接
	serverTLS, err := tls.Dial("tcp", host, &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"h2", "http/1.1"},
	})
	if err != nil {
		log.Println("dial tls error:", err)
		return err
	}

	// 4. 获取协商的协议
	proto := serverTLS.ConnectionState().NegotiatedProtocol
	if proto == "" {
		proto = "http/1.1"
	}

	// 5. 发送 200 Connection Established
	_, err = clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	if err != nil {
		log.Println("write connection established error:", err)
		serverTLS.Close()
		return err
	}

	// 6. 升级客户端连接为 TLS
	clientTLS, err := s.upgradeClientTLS(clientConn, tlsConfig, proto)
	if err != nil {
		log.Println("upgrade client tls error:", err)
		serverTLS.Close()
		return err
	}

	// 7. 根据协议选择处理器
	return s.interceptor.HandleConnection(clientTLS, serverTLS, proto, ctx)
}

// upgradeClientTLS 升级客户端连接为 TLS
func (s *ALPNSelector) upgradeClientTLS(conn net.Conn, config *tls.Config, proto string) (net.Conn, error) {
	// 设置 NextProtos
	if proto != "" {
		config.NextProtos = []string{proto}
	}

	tlsConn := tls.Server(conn, config)
	err := tlsConn.Handshake()
	if err != nil {
		return nil, err
	}

	return tlsConn, nil
}
