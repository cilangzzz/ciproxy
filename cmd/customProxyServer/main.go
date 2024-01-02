/**
  @creator: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @github: https://github.com/OpencvLZG
  @since: 2024/1/2
  @desc: //TODO
**/

package main

import (
	"bufio"
	"flag"
	"github.com/opencvlzg/ciproxy/constants/connectConfig"
	"github.com/opencvlzg/ciproxy/constants/proxyMethod"
	"github.com/opencvlzg/ciproxy/proxyServer/middleHandle"
	"github.com/opencvlzg/ciproxy/proxyServer/serve"
	"github.com/opencvlzg/ciproxy/proxyServer/trafficHandle"
	"net"
	"net/http"
	"strings"
)

func proxyTransfer(c net.Conn, s net.Conn) {
	go middleHandle.MiddleHandle(c, s)
	go trafficHandle.Transfer(c, s)
	go trafficHandle.Transfer(s, c)
}

func main() {
	ip := flag.String("ip", "127.0.0.1", "Server Ip Address")
	port := flag.String("port", "8888", "Server Port")
	method := flag.String("method", proxyMethod.CustomProxy, "Server METHOD NORMAL,TUNNEL, SNIFF")
	protocol := flag.String("protocol", "TCP", "Connect Protocol")
	logPath := flag.String("log", "log/proxy.log", "log file path")
	flag.Parse()
	proxyServe := serve.ProxyServe{
		Ip:       *ip,
		Port:     *port,
		Method:   *method,
		Protocol: *protocol,
		LogPath:  *logPath,
	}
	//customProxyHandle := func(c net.Conn) {
	//	proxyHandle.HttpsProxyHandle(c)
	//}
	customProxyHandle := func(c net.Conn) {
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
			return
		}
		switch request.Method {
		case "CONNECT":
			_, err := c.Write([]byte("HTTP/1.1 200 Connection Established \r\n\r\n"))
			if err != nil {
				return
			}
		default:

		}
		proxyTransfer(c, s)
	}
	proxyServe.SetProxyHandle(customProxyHandle)
	middleHandle.Add(func(client net.Conn, target net.Conn) {
		// Todo Some regular u want implement
		println(client.RemoteAddr().String())
	})
	proxyServe.Start()
}
