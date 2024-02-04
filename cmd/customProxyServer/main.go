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
	"github.com/opencvlzg/ciproxy"
	"net"
	"net/http"
	"strings"
)

func proxyTransfer(c net.Conn, s net.Conn) {
	//go MiddleHandle(c, s)
	go ciproxy.Transfer(c, s)
	go ciproxy.Transfer(s, c)
}

func main() {
	ip := flag.String("ip", "127.0.0.1", "Server Ip Address")
	port := flag.String("port", "8888", "Server Port")
	method := flag.String("method", ciproxy.CustomProxy, "Server METHOD NORMAL,TUNNEL, SNIFF")
	protocol := flag.String("protocol", "TCP", "Connect Protocol")
	logPath := flag.String("log", "log/proxy.log", "log file path")
	flag.Parse()
	proxyServe := ciproxy.ProxyServe{
		Ip:       *ip,
		Port:     *port,
		Method:   *method,
		Protocol: *protocol,
		LogPath:  *logPath,
	}
	//customProxyHandle := func(c net.Conn) {
	//	proxyHandle.HttpsProxyHandle(c)
	//}
	customProxyHandle := func(c *ciproxy.Context) {
		buf := bufio.NewReader(c.ClientConn)
		request, err := http.ReadRequest(buf)
		if err != nil {
			return
		}
		if !strings.HasSuffix(request.Host, ":443") {
			request.Host += ":443"
		}
		s, err := net.DialTimeout("tcp", request.Host, ciproxy.DefaultOutTime)
		if err != nil {
			return
		}
		switch request.Method {
		case "CONNECT":
			_, err := c.ClientConn.Write([]byte("HTTP/1.1 200 Connection Established \r\n\r\n"))
			if err != nil {
				return
			}
		default:

		}
		proxyTransfer(c.ClientConn, s)
	}
	proxyServe.SetProxyHandle(customProxyHandle)
	proxyServe.Start()
}
