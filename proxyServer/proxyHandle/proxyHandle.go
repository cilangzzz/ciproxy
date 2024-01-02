/**
  @creator: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @github: https://github.com/OpencvLZG
  @since: 2023/12/21
  @desc: //TODO
**/

// Package proxyHandle 代理响应处理头
package proxyHandle

import (
	"bufio"
	"crypto/tls"
	"github.com/opencvlzg/ciproxy/constants/connectConfig"
	"github.com/opencvlzg/ciproxy/proxyServer/trafficHandle"
	"github.com/opencvlzg/ciproxy/util"
	"log"
	"net"
	"net/http"
	"strings"
)

func errLog(msg string, err error) {
	log.Println("proxyHandle:" + msg + " err:" + err.Error())
}

// 转发流量 内部使用
func proxyTransfer(c net.Conn, s net.Conn) {
	//go middleHandle.MiddleHandle(c, s)
	go trafficHandle.Transfer(c, s)
	go trafficHandle.Transfer(s, c)
}

// HttpProxyHandle Http处理
func HttpProxyHandle(c net.Conn) {
	buf := bufio.NewReader(c)
	request, err := http.ReadRequest(buf)
	if err != nil {
		return
	}
	s, err := net.DialTimeout("tcp", request.Host, connectConfig.DefaultOutTime)
	if err != nil {
		errLog("remote host connect failed"+request.Host, err)
		return
	}
	proxyTransfer(c, s)
}

// HttpsProxyHandle Https处理
func HttpsProxyHandle(c net.Conn) {
	buf := bufio.NewReader(c)
	request, err := http.ReadRequest(buf)
	if err != nil {
		return
	}
	if !strings.HasSuffix(request.Host, ":443") {
		request.Host += ":443"
	}
	s, err := net.DialTimeout("tcp", request.Host, connectConfig.DefaultOutTime)
	if err != nil {
		errLog("remote host connect failed"+request.Host, err)
		return
	}
	switch request.Method {
	case "CONNECT":
		_, err := c.Write([]byte("HTTP/1.1 200 Connection Established \r\n\r\n"))
		if err != nil {
			errLog("write hello failed"+request.Host+request.Method, err)
			return
		}
	default:

	}
	proxyTransfer(c, s)

}

// HttpsSniffProxyHandle https中间人处理
func HttpsSniffProxyHandle(c net.Conn) {
	cReader := bufio.NewReader(c)
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
		errLog("remote host connect failed", err)
		return
	}
	_, err = c.Write([]byte("HTTP/1.1 200 Connection Established \r\n\r\n"))
	if err != nil {
		errLog("write hello failed"+request.Host+request.Method, err)
		return
	}
	tlsC, err := upgradeTls(c, tlsCnf)
	if err != nil {
		errLog("upgrade tls failed", err)
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
	proxyTransfer(tlsC, tlsS)

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
//			errLog("load ca certificate failed", err)
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
		errLog("close conn failed", err)
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
		errLog("tls handshake failed", err)
		return nil, err
	}

	return tlsC, nil
}

//// httpsTunnelResponse CONNECT 方法响应
//func httpsTunnelResponse(c net.Conn) {
//	_, err := c.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
//	if err != nil {
//		errLog("write hello failed", err)
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
//		errLog("upgrade tls failed", err)
//		return
//	}
//	//buf := bufio.NewReader(tlsC)
//	//request, err := http.ReadRequest(buf)
//	//if err != nil {
//	//	errLog("https encode filed", err)
//	//	return
//	//}
//	//println(request.Body)
//	println("握手成功")
//	tlsS, err := tls.Dial("tcp", "www.figma.com", conf)
//	if err != nil {
//		errLog("remote host connect failed", err)
//		return
//	}
//	proxyTransfer(tlsC, tlsS)
//}
