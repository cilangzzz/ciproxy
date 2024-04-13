/**
  @creator: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @github: https://github.com/OpencvLZG
  @since: 2023/12/21
  @desc: //TODO
**/

// Package ciproxy proxyHandle 代理响应处理头
package ciproxy

import (
	"bufio"
	"crypto/tls"
	"github.com/cilangzzz/ciproxy/internal/util"
	"log"
	"net"
	"net/http"
	"strings"
)

// 转发流量 内部使用
func proxyTransfer(c net.Conn, s net.Conn) {
	//go middleHandle.MiddleHandle(c, s)
	go Transfer(c, s)
	go Transfer(s, c)
}

// 转发流量 同时输出 内部使用
func proxyLogTransfer(c net.Conn, s net.Conn) {
	//go middleHandle.MiddleHandle(c, s)
	go TeeTransfer(c, s)
	go TeeTransfer(s, c)
}

// HttpProxyHandle Http处理
func HttpProxyHandle(c *Context) {

	buf := bufio.NewReader(c.ClientConn)
	request, err := http.ReadRequest(buf)
	if err != nil {
		return
	}
	c.ServerConn, err = net.DialTimeout("tcp", request.Host, DefaultOutTime)
	if err != nil {
		log.Println("remote host connect failed"+request.Host, err)
		return
	}
	proxyTransfer(c.ClientConn, c.ServerConn)
}

// HttpsProxyHandle Https处理
func HttpsProxyHandle(c *Context) {
	buf := bufio.NewReader(c.ClientConn)
	request, err := http.ReadRequest(buf)
	if err != nil {
		return
	}
	if !strings.HasSuffix(request.Host, ":443") {
		request.Host += ":443"
	}
	s, err := net.DialTimeout("tcp", request.Host, DefaultOutTime)
	if err != nil {
		log.Println("remote host connect failed"+request.Host, err)
		return
	}
	switch request.Method {
	case "CONNECT":
		_, err := c.ClientConn.Write([]byte("HTTP/1.1 200 Connection Established \r\n\r\n"))
		if err != nil {
			log.Println("write hello failed"+request.Host+request.Method, err)
			return
		}
	default:

	}
	proxyTransfer(c.ClientConn, s)

}

// HttpsSniffProxyHandle https中间人处理
func HttpsSniffProxyHandle(c *Context) {
	cReader := bufio.NewReader(c.ClientConn)
	request, err := http.ReadRequest(cReader)
	if err != nil {
		return
	}
	tlsCnf, err := util.GenerateTlsConfig(request.Host)
	if err != nil {
		return
	}
	if !strings.Contains(request.Host, ":443") {
		request.Host += ":443"
	}
	tlsS, err := tls.Dial("tcp", request.Host, tlsCnf)
	if err != nil {
		log.Println("remote host connect failed", err)
		return
	}
	_, err = c.ClientConn.Write([]byte("HTTP/1.1 200 Connection Established \r\n\r\n"))
	if err != nil {
		log.Println("write hello failed"+request.Host+request.Method, err)
		return
	}
	tlsC, err := upgradeTls(c.ClientConn, tlsCnf)
	if err != nil {
		log.Println("upgrade tls failed", err)
		closeConn(tlsC)
		closeConn(tlsS)
		return
	}
	//_, err = c.Write(util.HttpContext("元神"))
	//if err != nil {
	//	fmt.Printf("%s", err)
	//}
	//_, err = tlsC.Write(util.HttpContext("元神"))
	//if err != nil {
	//	fmt.Printf("%s", err)
	//}
	//fmt.Printf("%s\n", tlsS.RemoteAddr())
	//proxyTransfer(tlsC, tlsS)
	proxyTransfer(tlsC, tlsS)

}

// HttpsSniffDetailProxyHandle https中间人处理
func HttpsSniffDetailProxyHandle(c *Context) {
	cReader := bufio.NewReader(c.ClientConn)
	request, err := http.ReadRequest(cReader)
	if err != nil {
		return
	}
	tlsCnf, err := util.GenerateTlsConfig(request.Host)
	if err != nil {
		return
	}
	if !strings.Contains(request.Host, ":443") {
		request.Host += ":443"
	}
	tlsS, err := tls.Dial("tcp", request.Host, tlsCnf)
	if err != nil {
		log.Println("remote host connect failed", err)
		return
	}
	c.TlsClientConn = tlsS
	_, err = c.ClientConn.Write([]byte("HTTP/1.1 200 Connection Established \r\n\r\n"))
	if err != nil {
		log.Println("write hello failed"+request.Host+request.Method, err)
		return
	}
	tlsC, err := upgradeTls(c.ClientConn, tlsCnf)
	if err != nil {
		log.Println("upgrade tls failed", err)
		return
	}
	c.TlsClientConn = tlsC
	//_, err = c.Write(util.HttpContext("元神"))
	//if err != nil {
	//	fmt.Printf("%s", err)
	//}
	//_, err = tlsC.Write(util.HttpContext("元神"))
	//if err != nil {
	//	fmt.Printf("%s", err)
	//}
	//fmt.Printf("%s\n", tlsS.RemoteAddr())
	//proxyTransfer(tlsC, tlsS)
	//bytes := util.HttpContext("你好")
	//go tlsC.Write(bytes)

	go TeeDoRequestTransfer(c)

}

// TunnelProxyHandle 加密代理
func TunnelProxyHandle(c *Context) {
	buf := bufio.NewReader(c.ClientConn)
	request, err := http.ReadRequest(buf)
	if err != nil {
		return
	}
	if !strings.HasSuffix(request.Host, ":443") {
		request.Host += ":443"
	}
	s, err := net.DialTimeout("tcp", request.Host, DefaultOutTime)
	if err != nil {
		log.Println("remote host connect failed"+request.Host, err)
		return
	}
	switch request.Method {
	case "CONNECT":
		_, err := c.ClientConn.Write([]byte("HTTP/1.1 200 Connection Established \r\n\r\n"))
		if err != nil {
			log.Println("write hello failed"+request.Host+request.Method, err)
			return
		}
	default:

	}
	Transfer(c.ClientConn, s)
}

// isWebSocketUpgrade 判断是否为websocket链接
func isWebSocketUpgrade(req *http.Request) bool {
	upgradeHeader := req.Header.Get("Upgrade")
	connectionHeader := req.Header.Get("Connection")

	return strings.ToLower(upgradeHeader) == "websocket" && strings.Contains(strings.ToLower(connectionHeader), "upgrade")
}

// WebsocketProxyHandle websocket 代理
func WebsocketProxyHandle(c *Context) {
	buf := bufio.NewReader(c.ClientConn)
	request, err := http.ReadRequest(buf)
	if err != nil {
		return
	}
	if !strings.HasSuffix(request.Host, ":443") {
		request.Host += ":443"
	}
	s, err := net.DialTimeout("tcp", request.Host, DefaultOutTime)
	if err != nil {
		log.Println("remote host connect failed"+request.Host, err)
		return
	}
	if isWebSocketUpgrade(request) {
		_, err = s.Write([]byte("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\n\r\n"))
		if err != nil {
			log.Println("upgrade failed"+request.Host, err)
			return
		}
	} else {
		return
	}
	Transfer(c.ClientConn, s)
	//switch request.Method {
	//case "CONNECT":
	//	_, err := c.ClientConn.Write([]byte("HTTP/1.1 200 Connection Established \r\n\r\n"))
	//	if err != nil {
	//		log.Println("write hello failed"+request.Host+request.Method, err)
	//		return
	//	}
	//default:
	//
	//}
	//proxyTransfer(c.ClientConn, s)
}

// TestProxyHandle 测试代理头
//
//	func TestProxyHandle(c net.Conn) {
//		buf := make([]byte, 512)
//		_, err := c.Read(buf)
//		fmt.Printf("%s", buf)
//		_, err = c.Write([]byte("HTTP/1.1 200 Connection Established \r\n\r\n"))
//		if err != nil {
//			//return
//		}
//
//		cert, err := util.LoadCertificate("./cert/www.cilang.buzz/inter.crt", "./cert/www.cilang.buzz/inter.key")
//		if err != nil {
//			log.Println("load ca certificate failed", err)
//			panic(err)
//			return
//		}
//		conf := &tls.Config{
//			Certificates:       []tls.Certificate{*cert},
//			InsecureSkipVerify: true,
//			MinVersion:         tls.VersionTLS12,
//			MaxVersion:         tls.VersionTLS13,
//		}
//		_, err = upgradeTls(c, conf)
//		if err != nil {
//			return
//		}
//	}
//
// closeConn 关闭连接
func closeConn(c net.Conn) {
	err := c.Close()
	if err != nil {
		log.Println("close conn failed", err)
		return
	}
}

// upgradeTls 从tcp升级到tls连接
func upgradeTls(c net.Conn, conf *tls.Config) (net.Conn, error) {

	tlsC := tls.Server(c, conf)
	//defer func() {
	//	_ = tlsC.Close()
	//}()
	err := tlsC.Handshake()
	if err != nil {
		log.Println("tls handshake failed", err)
		return nil, err
	}

	return tlsC, nil
}

//// httpsTunnelResponse CONNECT 方法响应
//func httpsTunnelResponse(c net.Conn) {
//	_, err := c.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
//	if err != nil {
//		log.Println("write hello failed", err)
//		return
//	}
//}

//// httpsTlsResponse 通信协议识别
//func httpsTlsResponse(c net.Conn) {
//	conf, err := util.GenerateTlsConfig("www.figma.com")
//	if err != nil {
//		return
//	}
//	tlsC, err := upgradeTls(c, conf)
//	if err != nil {
//		log.Println("upgrade tls failed", err)
//		return
//	}
//	//buf := bufio.NewReader(tlsC)
//	//request, err := http.ReadRequest(buf)
//	//if err != nil {
//	//	log.Println("https encode filed", err)
//	//	return
//	//}
//	//println(request.Body)
//	println("握手成功")
//	tlsS, err := tls.Dial("tcp", "www.figma.com", conf)
//	if err != nil {
//		log.Println("remote host connect failed", err)
//		return
//	}
//	proxyTransfer(tlsC, tlsS)
//}
