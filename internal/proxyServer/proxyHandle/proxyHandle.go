/**
  @creator: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @github: https://github.com/OpencvLZG
  @since: 2023/12/21
  @desc: //TODO
**/

package proxyHandle

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"trafficForward/internal/constant"
	"trafficForward/internal/proxyServer/middleHandle"
	"trafficForward/internal/proxyServer/trafficHandle"
)

func errLog(msg string, err error) {
	log.Println("proxyHandle:" + msg + "err:" + err.Error())
}

// 转发流量 内部使用
func proxyTransfer(c net.Conn, s net.Conn) {
	go middleHandle.MiddleHandle(c, s)
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
	s, err := net.DialTimeout("tcp", request.Host, constant.DefaultOutTime)
	if err != nil {
		errLog("remote host connect failed"+request.Host, err)
		return
	}
	mdManage := middleHandle.MdManage
	for _, handle := range mdManage.HandleChain {
		handle(c, s)
	}
	proxyTransfer(c, s)
}

// tlsConfig := util2.TLSUtil{Organization: "CiproxyOrganization"}
// cert, err := tlsConfig.GenCertificate()
// if err != nil {
// log.Fatal(err)
// }
// config := &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}

// HttpsProxyHandle Https处理
func HttpsProxyHandle(c net.Conn) {
	buf := bufio.NewReader(c)
	request, err := http.ReadRequest(buf)
	if err != nil {
		return
	}
	s, err := net.DialTimeout("tcp", request.Host, constant.DefaultOutTime)
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
