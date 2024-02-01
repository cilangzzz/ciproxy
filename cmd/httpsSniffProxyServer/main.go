/**
  @creator: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @github: https://github.com/OpencvLZG
  @since: 2023/12/24
  @desc: //TODO
**/

package main

import (
	"bufio"
	"bytes"
	"flag"
	"github.com/opencvlzg/ciproxy"
	"github.com/opencvlzg/ciproxy/constants/proxyMethod"
	"github.com/opencvlzg/ciproxy/proxyServer/serve"
	"io"
	"net/http"
	"os"
)

var buffer = bytes.NewBuffer([]byte{})
var writer = io.MultiWriter(buffer)
var buf = bufio.NewReader(buffer)

func logger() {

	for {
		request, err := http.ReadRequest(buf)
		if err != nil {
			//return
		} else {
			println(request.Method, request.Host, buf.Size())

		}

	}
}

func main() {
	ip := flag.String("ip", "127.0.0.1", "Server Ip Address")
	port := flag.String("port", "8888", "Server Port")
	method := flag.String("method", proxyMethod.HttpsSniffProxy, "Server METHOD NORMAL,TUNNEL, SNIFF")
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
	file, err := os.Create("data.txt")
	if err != nil {
		return
	}
	defer file.Close()
	ciproxy.DefaultWriter = bufio.NewWriter(file)
	ciproxy.DefaultWriter = writer
	//go logger()
	//middleware := middleHandle.MdManage
	//middleware.Add(func(client net.Conn, target net.Conn) {
	//	// Todo Some regular u want implement
	//})
	proxyServe.Start()
}
